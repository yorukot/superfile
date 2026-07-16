package filesystem

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/pkg/sftp"
	cryptossh "golang.org/x/crypto/ssh"

	internalssh "github.com/yorukot/superfile/src/internal/ssh"
	"github.com/yorukot/superfile/src/pkg/utils"
)

const (
	sftpProviderName       = "SFTP remote filesystem"
	defaultSFTPOpenTimeout = 10 * time.Second
)

type SFTPProvider struct {
	request      internalssh.ClientConfigRequest
	capabilities CapabilitySet
}

type SFTPSession struct {
	id           SessionID
	root         Location
	capabilities CapabilitySet
	bundle       *internalssh.ClientConfigBundle
	sshClient    *cryptossh.Client
	client       *sftp.Client
	closeOnce    sync.Once
	closeErr     error
}

func NewSFTPProvider(request internalssh.ClientConfigRequest) *SFTPProvider {
	return &SFTPProvider{request: request, capabilities: V1CapabilityMatrix()}
}

func (p *SFTPProvider) Kind() ProviderKind {
	return ProviderSFTP
}

func (p *SFTPProvider) Name() string {
	return sftpProviderName
}

func (p *SFTPProvider) Capabilities() CapabilitySet {
	return p.capabilities
}

func (p *SFTPProvider) Open(ctx context.Context, location Location) (Session, error) {
	if err := ctx.Err(); err != nil {
		return nil, mapSFTPError(OperationNavigate, location.Path, err)
	}
	openTimeout := p.request.Timeout
	if openTimeout <= 0 {
		openTimeout = defaultSFTPOpenTimeout
	}
	openCtx, cancel := context.WithTimeout(ctx, openTimeout)
	defer cancel()

	root := location
	root.Provider = ProviderSFTP
	if root.Path.String() == "" {
		if p.request.Profile.StartPath != "" {
			root.Path = NewRemotePath(p.request.Profile.StartPath)
		} else {
			root.Path = RootRemotePath()
		}
	}
	if !root.Path.IsRemote() {
		return nil, NewUnsupportedError(ProviderSFTP, OperationNavigate, root.Path,
			"sftp provider requires remote POSIX paths")
	}
	if root.Label == "" {
		root.Label = p.request.Profile.Name
		if root.Label == "" {
			root.Label = p.request.Profile.Host
		}
	}
	if root.SessionID == "" {
		root.SessionID = SessionID(root.Label)
	}

	bundle, err := internalssh.BuildClientConfig(p.request)
	if err != nil {
		return nil, err
	}
	sshClient, err := bundle.DialContext(openCtx)
	if err != nil {
		_ = bundle.Close()
		return nil, err
	}
	sftpClient, err := runSFTPOperation(openCtx, sshClient.Close, func() (*sftp.Client, error) {
		return sftp.NewClient(sshClient)
	})
	if err != nil {
		_ = sshClient.Close()
		_ = bundle.Close()
		return nil, mapSFTPError(OperationNavigate, root.Path, err)
	}

	return &SFTPSession{
		id:           root.SessionID,
		root:         root,
		capabilities: p.capabilities,
		bundle:       bundle,
		sshClient:    sshClient,
		client:       sftpClient,
	}, nil
}

func (s *SFTPSession) ID() SessionID {
	return s.id
}

func (s *SFTPSession) Provider() ProviderKind {
	return ProviderSFTP
}

func (s *SFTPSession) Root() Location {
	return s.root
}

func (s *SFTPSession) Capabilities() CapabilitySet {
	return s.capabilities
}

func (s *SFTPSession) List(ctx context.Context, path Path) ([]Entry, error) {
	remotePath, err := s.requireRemotePath(OperationList, path)
	if err != nil {
		return nil, err
	}
	if contextErr := ctx.Err(); contextErr != nil {
		return nil, mapSFTPError(OperationList, path, contextErr)
	}

	infos, err := s.client.ReadDirContext(ctx, remotePath)
	if err != nil {
		return nil, mapSFTPError(OperationList, path, err)
	}

	entries := make([]Entry, 0, len(infos))
	base := NewRemotePath(remotePath)
	for _, info := range infos {
		entryPath := base.Join(info.Name())
		stat := s.newStat(entryPath, info)
		entries = append(entries, Entry{Name: info.Name(), Path: entryPath, Stat: stat})
	}
	return entries, nil
}

func (s *SFTPSession) Stat(ctx context.Context, path Path) (Stat, error) {
	remotePath, err := s.requireRemotePath(OperationStat, path)
	if err != nil {
		return Stat{}, err
	}
	if contextErr := ctx.Err(); contextErr != nil {
		return Stat{}, mapSFTPError(OperationStat, path, contextErr)
	}

	info, err := runSFTPOperation(ctx, s.Close, func() (os.FileInfo, error) {
		return s.client.Lstat(remotePath)
	})
	if err != nil {
		return Stat{}, mapSFTPError(OperationStat, path, err)
	}
	return s.newStat(path, info), nil
}

func (s *SFTPSession) Read(ctx context.Context, path Path) (io.ReadCloser, error) {
	remotePath, err := s.requireRemotePath(OperationPreviewRead, path)
	if err != nil {
		return nil, err
	}
	if contextErr := ctx.Err(); contextErr != nil {
		return nil, mapSFTPError(OperationPreviewRead, path, contextErr)
	}

	file, err := runSFTPOperation(ctx, s.Close, func() (*sftp.File, error) {
		return s.client.Open(remotePath)
	})
	if err != nil {
		return nil, mapSFTPError(OperationPreviewRead, path, err)
	}
	return contextReadCloser{ctx: ctx, path: path, operation: OperationPreviewRead, reader: file}, nil
}

func (s *SFTPSession) Create(ctx context.Context, path Path, reader io.Reader, options CreateOptions) error {
	remotePath, err := s.requireRemotePath(OperationCreateFile, path)
	if err != nil {
		return err
	}
	if contextErr := ctx.Err(); contextErr != nil {
		return mapSFTPError(OperationCreateFile, path, contextErr)
	}

	flags := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	if !options.Overwrite {
		flags |= os.O_EXCL
	}
	file, err := s.client.OpenFile(remotePath, flags)
	if err != nil {
		return mapSFTPCreateError(OperationCreateFile, path, err)
	}
	defer file.Close()

	mode := options.Mode
	if mode == 0 {
		mode = utils.UserFilePerm
	}
	if chmodErr := s.client.Chmod(remotePath, mode); chmodErr != nil {
		return mapSFTPError(OperationCreateFile, path, chmodErr)
	}

	if reader == nil {
		return nil
	}
	_, err = file.ReadFrom(contextReader{ctx: ctx, path: path, operation: OperationCreateFile, reader: reader})
	return mapSFTPError(OperationCreateFile, path, err)
}

func (s *SFTPSession) Mkdir(ctx context.Context, path Path, options MkdirOptions) error {
	remotePath, err := s.requireRemotePath(OperationMkdir, path)
	if err != nil {
		return err
	}
	if contextErr := ctx.Err(); contextErr != nil {
		return mapSFTPError(OperationMkdir, path, contextErr)
	}

	if options.Parents {
		err = s.client.MkdirAll(remotePath)
	} else {
		err = s.client.Mkdir(remotePath)
	}
	if err != nil {
		return mapSFTPCreateError(OperationMkdir, path, err)
	}

	mode := options.Mode
	if mode == 0 {
		mode = utils.UserDirPerm
	}
	return mapSFTPError(OperationMkdir, path, s.client.Chmod(remotePath, mode))
}

func (s *SFTPSession) Rename(ctx context.Context, source Path, destination Path, options RenameOptions) error {
	sourcePath, err := s.requireRemotePath(OperationRename, source)
	if err != nil {
		return err
	}
	destinationPath, err := s.requireRemotePath(OperationRename, destination)
	if err != nil {
		return err
	}
	if contextErr := ctx.Err(); contextErr != nil {
		return mapSFTPError(OperationRename, source, contextErr)
	}

	if !options.Overwrite {
		if exists, existsErr := s.exists(destination); existsErr != nil {
			return existsErr
		} else if exists {
			return NewConflictError(ProviderSFTP, OperationRename, destination, "destination already exists")
		}
		return mapSFTPError(OperationRename, source, s.client.Rename(sourcePath, destinationPath))
	}

	if err = s.client.PosixRename(sourcePath, destinationPath); err == nil {
		return nil
	}
	if isSFTPOpUnsupported(err) {
		return NewUnsupportedError(ProviderSFTP, OperationRename, destination,
			"sftp server does not support atomic overwrite rename")
	}
	return mapSFTPError(OperationRename, source, err)
}

func (s *SFTPSession) Delete(ctx context.Context, path Path, options DeleteOptions) error {
	if options.UseTrash {
		return NewUnsupportedError(ProviderSFTP, OperationDeleteFile, path,
			"remote trash is not supported for SFTP sessions")
	}
	if err := ctx.Err(); err != nil {
		return mapSFTPError(OperationDeleteFile, path, err)
	}

	stat, err := s.Stat(ctx, path)
	if err != nil {
		return err
	}
	if stat.IsDir && !stat.IsSymlink {
		if !options.Recursive {
			return mapSFTPError(OperationDeleteDir, path, s.client.RemoveDirectory(path.String()))
		}
		return s.deleteDir(ctx, path)
	}
	return mapSFTPError(OperationDeleteFile, path, s.client.Remove(path.String()))
}

func (s *SFTPSession) Copy(ctx context.Context, source Path, destination Path, options CopyOptions) error {
	if err := ctx.Err(); err != nil {
		return mapSFTPError(OperationCopy, source, err)
	}
	sourcePath, err := s.requireRemotePath(OperationCopy, source)
	if err != nil {
		return err
	}
	if _, err = s.requireRemotePath(OperationCopy, destination); err != nil {
		return err
	}

	info, err := s.client.Lstat(sourcePath)
	if err != nil {
		return mapSFTPError(OperationCopy, source, err)
	}
	if info.IsDir() && info.Mode()&os.ModeSymlink == 0 {
		if !options.Recursive {
			return NewUnsupportedError(ProviderSFTP, OperationCopy, source,
				"copying directories requires recursive copy")
		}
		return s.copyDir(ctx, source, destination, info, options)
	}
	return s.copyFile(ctx, source, destination, info, options)
}

func (s *SFTPSession) Move(ctx context.Context, source Path, destination Path, options MoveOptions) error {
	if err := s.Rename(ctx, source, destination, RenameOptions{Overwrite: options.Overwrite}); err == nil {
		return nil
	} else if !isSFTPRecoverableRenameError(err) {
		return err
	} else if destinationExists, existsErr := s.exists(destination); existsErr != nil {
		return existsErr
	} else if destinationExists {
		return err
	}

	err := s.Copy(ctx, source, destination, CopyOptions(options))
	if err != nil {
		return fmt.Errorf("failed to copy before remote move delete: %w", err)
	}
	return s.Delete(ctx, source, DeleteOptions{Recursive: true})
}

func (s *SFTPSession) Chmod(ctx context.Context, path Path, mode os.FileMode) error {
	return NewUnsupportedError(ProviderSFTP, OperationChmod, path,
		"remote chmod is deferred until the UI exposes chmod explicitly")
}

func (s *SFTPSession) Transfer(_ context.Context, request TransferRequest) (Transfer, error) {
	operation := request.Operation
	if operation == "" {
		operation = OperationCopy
	}
	return nil, NewUnsupportedError(ProviderSFTP, operation, request.Destination.Path,
		"transfer orchestration is handled by the provider-aware transfer engine")
}

func (s *SFTPSession) Close() error {
	s.closeOnce.Do(func() {
		if s.client != nil {
			s.closeErr = errors.Join(s.closeErr, s.client.Close())
		}
		if s.sshClient != nil {
			s.closeErr = errors.Join(s.closeErr, s.sshClient.Close())
		}
		if s.bundle != nil {
			s.closeErr = errors.Join(s.closeErr, s.bundle.Close())
		}
	})
	return s.closeErr
}

type sftpOperationResult[T any] struct {
	value T
	err   error
}

func runSFTPOperation[T any](ctx context.Context, cancel func() error, operation func() (T, error)) (T, error) {
	result := make(chan sftpOperationResult[T], 1)
	go func() {
		value, err := operation()
		result <- sftpOperationResult[T]{value: value, err: err}
	}()
	select {
	case completed := <-result:
		return completed.value, completed.err
	case <-ctx.Done():
		_ = cancel()
		var zero T
		return zero, ctx.Err()
	}
}

func (s *SFTPSession) requireRemotePath(operation Operation, path Path) (string, error) {
	if err := s.capabilities.RequireRemote(ProviderSFTP, operation, path); err != nil {
		return "", err
	}
	if !path.IsRemote() {
		return "", NewUnsupportedError(ProviderSFTP, operation, path,
			"sftp provider requires remote POSIX paths")
	}
	return path.String(), nil
}

func (s *SFTPSession) deleteDir(ctx context.Context, directory Path) error {
	entries, err := s.List(ctx, directory)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if err := ctx.Err(); err != nil {
			return mapSFTPError(OperationDeleteDir, directory, err)
		}
		if entry.Stat.IsDir && !entry.Stat.IsSymlink {
			if err := s.deleteDir(ctx, entry.Path); err != nil {
				return err
			}
		} else if err := s.client.Remove(entry.Path.String()); err != nil {
			return mapSFTPError(OperationDeleteFile, entry.Path, err)
		}
	}
	return mapSFTPError(OperationDeleteDir, directory, s.client.RemoveDirectory(directory.String()))
}

func (s *SFTPSession) copyDir(
	ctx context.Context,
	source Path,
	destination Path,
	sourceInfo os.FileInfo,
	options CopyOptions,
) error {
	if err := s.prepareDestination(destination, OperationCopy, options.Overwrite); err != nil {
		return err
	}
	if err := s.Mkdir(
		ctx,
		destination,
		MkdirOptions{Mode: sourceInfo.Mode(), Parents: true},
	); err != nil &&
		!errors.Is(err, ErrConflict) {
		return err
	}
	entries, err := s.List(ctx, source)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if contextErr := ctx.Err(); contextErr != nil {
			return mapSFTPError(OperationCopy, source, contextErr)
		}
		entrySource := source.Join(entry.Name)
		entryDestination := destination.Join(entry.Name)
		if entry.Stat.IsDir && !entry.Stat.IsSymlink {
			err = s.copyDir(ctx, entrySource, entryDestination, entry.Stat.fileInfo(), options)
		} else {
			err = s.copyFile(ctx, entrySource, entryDestination, entry.Stat.fileInfo(), options)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SFTPSession) copyFile(
	ctx context.Context,
	source Path,
	destination Path,
	sourceInfo os.FileInfo,
	options CopyOptions,
) error {
	if err := s.prepareDestination(destination, OperationCopy, options.Overwrite); err != nil {
		return err
	}
	sourceFile, err := s.client.Open(source.String())
	if err != nil {
		return mapSFTPError(OperationCopy, source, err)
	}
	defer sourceFile.Close()

	destinationFile, err := s.client.OpenFile(destination.String(), os.O_WRONLY|os.O_CREATE|os.O_TRUNC)
	if err != nil {
		return mapSFTPCreateError(OperationCopy, destination, err)
	}
	defer destinationFile.Close()

	_, err = destinationFile.ReadFrom(
		contextReader{ctx: ctx, path: source, operation: OperationCopy, reader: sourceFile},
	)
	if err != nil {
		return mapSFTPError(OperationCopy, source, err)
	}
	return mapSFTPError(OperationCopy, destination, s.client.Chmod(destination.String(), sourceInfo.Mode()))
}

func (s *SFTPSession) prepareDestination(destination Path, operation Operation, overwrite bool) error {
	exists, err := s.exists(destination)
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}
	if !overwrite {
		return NewConflictError(ProviderSFTP, operation, destination, "destination already exists")
	}
	return nil
}

func (s *SFTPSession) exists(path Path) (bool, error) {
	_, err := s.client.Lstat(path.String())
	if err == nil {
		return true, nil
	}
	if isSFTPNotFound(err) {
		return false, nil
	}
	mappedErr := mapSFTPError(OperationStat, path, err)
	if errors.Is(mappedErr, ErrNotFound) {
		return false, nil
	}
	return false, mappedErr
}

func (s *SFTPSession) newStat(path Path, info os.FileInfo) Stat {
	stat := Stat{
		Name:       info.Name(),
		Size:       info.Size(),
		Mode:       info.Mode(),
		ModTime:    info.ModTime(),
		IsDir:      info.IsDir(),
		IsSymlink:  info.Mode()&os.ModeSymlink != 0,
		ProviderID: string(ProviderSFTP),
	}
	if stat.IsSymlink {
		if target, err := s.client.ReadLink(path.String()); err == nil {
			stat.Target = NewRemotePath(target)
		}
	}
	return stat
}

func (s Stat) AsFileInfo() os.FileInfo {
	return statFileInfo{s: s}
}

func (s Stat) fileInfo() os.FileInfo {
	return s.AsFileInfo()
}

type statFileInfo struct {
	s Stat
}

func (i statFileInfo) Name() string       { return i.s.Name }
func (i statFileInfo) Size() int64        { return i.s.Size }
func (i statFileInfo) Mode() os.FileMode  { return i.s.Mode }
func (i statFileInfo) ModTime() time.Time { return i.s.ModTime }
func (i statFileInfo) IsDir() bool        { return i.s.IsDir }
func (i statFileInfo) Sys() any           { return nil }

type contextReader struct {
	ctx       context.Context
	path      Path
	operation Operation
	reader    io.Reader
}

func (r contextReader) Read(p []byte) (int, error) {
	if err := r.ctx.Err(); err != nil {
		return 0, mapSFTPError(r.operation, r.path, err)
	}
	n, err := r.reader.Read(p)
	if err != nil {
		return n, mapSFTPError(r.operation, r.path, err)
	}
	return n, nil
}

type contextReadCloser struct {
	ctx       context.Context
	path      Path
	operation Operation
	reader    io.ReadCloser
}

func (r contextReadCloser) Read(p []byte) (int, error) {
	if err := r.ctx.Err(); err != nil {
		return 0, mapSFTPError(r.operation, r.path, err)
	}
	n, err := r.reader.Read(p)
	if err != nil {
		return n, mapSFTPError(r.operation, r.path, err)
	}
	return n, nil
}

func (r contextReadCloser) Close() error {
	return r.reader.Close()
}

func mapSFTPCreateError(operation Operation, path Path, err error) error {
	if err == nil {
		return nil
	}
	if isSFTPConflict(err) {
		return NewConflictError(ProviderSFTP, operation, path, "destination already exists")
	}
	return mapSFTPError(operation, path, err)
}

func mapSFTPError(operation Operation, path Path, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return NewCanceledError(ProviderSFTP, operation, path, err.Error())
	}
	if errors.Is(err, os.ErrPermission) || isSFTPPermission(err) {
		return NewPermissionError(ProviderSFTP, operation, path, err.Error())
	}
	if errors.Is(err, os.ErrNotExist) || isSFTPNotFound(err) {
		return NewNotFoundError(ProviderSFTP, operation, path, err.Error())
	}
	if isSFTPDisconnected(err) {
		return NewDisconnectedError(ProviderSFTP, operation, path, err.Error())
	}
	if isSFTPOpUnsupported(err) {
		return NewUnsupportedError(ProviderSFTP, operation, path, err.Error())
	}
	return err
}

func isSFTPConflict(err error) bool {
	if err == nil {
		return false
	}
	message := strings.ToLower(err.Error())
	return errors.Is(err, os.ErrExist) || strings.Contains(message, "file exists") ||
		strings.Contains(message, "already exists")
}

func isSFTPPermission(err error) bool {
	var status *sftp.StatusError
	if !errors.As(err, &status) {
		return false
	}
	message := strings.ToLower(status.Error())
	return status.FxCode() == sftp.ErrSSHFxPermissionDenied || strings.Contains(message, "permission denied")
}

func isSFTPNotFound(err error) bool {
	var status *sftp.StatusError
	if !errors.As(err, &status) {
		return false
	}
	message := strings.ToLower(status.Error())
	return status.FxCode() == sftp.ErrSSHFxNoSuchFile || strings.Contains(message, "file does not exist") ||
		strings.Contains(message, "no such file")
}

func isSFTPDisconnected(err error) bool {
	var status *sftp.StatusError
	if errors.As(err, &status) {
		return status.FxCode() == sftp.ErrSSHFxNoConnection || status.FxCode() == sftp.ErrSSHFxConnectionLost
	}
	message := strings.ToLower(err.Error())
	return strings.Contains(message, "connection lost") || strings.Contains(message, "no connection") ||
		strings.Contains(message, "broken pipe") || strings.Contains(message, "failed to send packet") ||
		strings.Contains(message, "failed to receive packet") || strings.Contains(message, "connection reset") ||
		strings.Contains(message, "use of closed network connection")
}

func isSFTPOpUnsupported(err error) bool {
	var status *sftp.StatusError
	if !errors.As(err, &status) {
		return false
	}
	message := strings.ToLower(status.Error())
	return status.FxCode() == sftp.ErrSSHFxOpUnsupported || strings.Contains(message, "op unsupported") ||
		strings.Contains(message, "unsupported")
}

func isSFTPRecoverableRenameError(err error) bool {
	return errors.Is(err, ErrUnsupported) || errors.Is(err, ErrConflict) || isSFTPNotFound(err)
}
