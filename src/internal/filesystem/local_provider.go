package filesystem

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	trashWin "github.com/hymkor/trash-go"
	"github.com/rkoesters/xdg/trash"

	variable "github.com/yorukot/superfile/src/config"
	"github.com/yorukot/superfile/src/pkg/utils"
)

const localSessionLabel = "local"

type LocalProvider struct {
	capabilities CapabilitySet
}

type LocalSession struct {
	id           SessionID
	root         Location
	capabilities CapabilitySet
}

func NewLocalProvider() *LocalProvider {
	return &LocalProvider{capabilities: V1CapabilityMatrix()}
}

func (p *LocalProvider) Kind() ProviderKind {
	return ProviderLocal
}

func (p *LocalProvider) Name() string {
	return "Local filesystem"
}

func (p *LocalProvider) Capabilities() CapabilitySet {
	return p.capabilities
}

func (p *LocalProvider) Open(ctx context.Context, location Location) (Session, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	root := location
	root.Provider = ProviderLocal
	if root.Path.String() == "" {
		root.Path = NewLocalPath(string(filepath.Separator))
	}
	if !root.Path.IsLocal() {
		return nil, NewUnsupportedError(ProviderLocal, OperationNavigate, root.Path,
			"local provider requires local filesystem paths")
	}
	if root.Label == "" {
		root.Label = localSessionLabel
	}
	if root.SessionID == "" {
		root.SessionID = SessionID(localSessionLabel)
	}

	return &LocalSession{id: root.SessionID, root: root, capabilities: p.capabilities}, nil
}

func (s *LocalSession) ID() SessionID {
	return s.id
}

func (s *LocalSession) Provider() ProviderKind {
	return ProviderLocal
}

func (s *LocalSession) Root() Location {
	return s.root
}

func (s *LocalSession) Capabilities() CapabilitySet {
	return s.capabilities
}

func (s *LocalSession) List(ctx context.Context, path Path) ([]Entry, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	localPath, err := s.requireLocalPath(OperationList, path)
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(localPath)
	if err != nil {
		return nil, err
	}

	result := make([]Entry, 0, len(entries))
	for _, entry := range entries {
		entryPath := filepath.Join(localPath, entry.Name())
		info, infoErr := os.Lstat(entryPath)
		if infoErr != nil {
			return nil, infoErr
		}
		result = append(result, Entry{
			Name: entry.Name(),
			Path: NewLocalPath(entryPath),
			Stat: newLocalStat(entryPath, info),
		})
	}

	return result, nil
}

func (s *LocalSession) Stat(ctx context.Context, path Path) (Stat, error) {
	if err := ctx.Err(); err != nil {
		return Stat{}, err
	}

	localPath, err := s.requireLocalPath(OperationStat, path)
	if err != nil {
		return Stat{}, err
	}

	info, err := os.Lstat(localPath)
	if err != nil {
		return Stat{}, err
	}

	return newLocalStat(localPath, info), nil
}

func (s *LocalSession) Read(ctx context.Context, path Path) (io.ReadCloser, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	localPath, err := s.requireLocalPath(OperationPreviewRead, path)
	if err != nil {
		return nil, err
	}

	return os.Open(localPath)
}

func (s *LocalSession) Create(ctx context.Context, path Path, reader io.Reader, options CreateOptions) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	localPath, err := s.requireLocalPath(OperationCreateFile, path)
	if err != nil {
		return err
	}

	mode := options.Mode
	if mode == 0 {
		mode = utils.UserFilePerm
	}

	flags := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	if !options.Overwrite {
		flags = os.O_WRONLY | os.O_CREATE | os.O_EXCL
	}

	file, err := os.OpenFile(localPath, flags, mode)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return NewConflictError(ProviderLocal, OperationCreateFile, path, err.Error())
		}
		return err
	}
	defer file.Close()

	if reader == nil {
		return nil
	}

	_, err = io.Copy(file, reader)
	return err
}

func (s *LocalSession) Mkdir(ctx context.Context, path Path, options MkdirOptions) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	localPath, err := s.requireLocalPath(OperationMkdir, path)
	if err != nil {
		return err
	}

	mode := options.Mode
	if mode == 0 {
		mode = utils.UserDirPerm
	}

	if options.Parents {
		return os.MkdirAll(localPath, mode)
	}

	return os.Mkdir(localPath, mode)
}

func (s *LocalSession) Rename(ctx context.Context, source Path, destination Path, options RenameOptions) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	sourcePath, err := s.requireLocalPath(OperationRename, source)
	if err != nil {
		return err
	}
	destinationPath, err := s.requireLocalPath(OperationRename, destination)
	if err != nil {
		return err
	}

	if !options.Overwrite {
		if _, statErr := os.Stat(destinationPath); statErr == nil {
			return NewConflictError(ProviderLocal, OperationRename, destination, "destination already exists")
		} else if !errors.Is(statErr, os.ErrNotExist) {
			return statErr
		}
	}

	return os.Rename(sourcePath, destinationPath)
}

func (s *LocalSession) Delete(ctx context.Context, path Path, options DeleteOptions) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	localPath, err := s.requireLocalPath(OperationDeleteFile, path)
	if err != nil {
		return err
	}

	if options.UseTrash {
		return s.moveToTrash(localPath, path)
	}

	if options.Recursive {
		return os.RemoveAll(localPath)
	}

	return os.Remove(localPath)
}

func (s *LocalSession) Copy(ctx context.Context, source Path, destination Path, options CopyOptions) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	sourcePath, err := s.requireLocalPath(OperationCopy, source)
	if err != nil {
		return err
	}
	destinationPath, err := s.requireLocalPath(OperationCopy, destination)
	if err != nil {
		return err
	}

	return s.copyPath(ctx, source, sourcePath, destination, destinationPath, options)
}

func (s *LocalSession) Move(ctx context.Context, source Path, destination Path, options MoveOptions) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	sourcePath, err := s.requireLocalPath(OperationCutMove, source)
	if err != nil {
		return err
	}
	destinationPath, err := s.requireLocalPath(OperationCutMove, destination)
	if err != nil {
		return err
	}

	sameDev, err := localPathsSharePartition(sourcePath, destinationPath)
	if err != nil {
		return fmt.Errorf("failed to check partitions: %w", err)
	}

	if sameDev {
		if !options.Overwrite {
			if _, statErr := os.Stat(destinationPath); statErr == nil {
				return NewConflictError(ProviderLocal, OperationCutMove, destination, "destination already exists")
			} else if !errors.Is(statErr, os.ErrNotExist) {
				return statErr
			}
		}
		if err = os.Rename(sourcePath, destinationPath); err == nil {
			return nil
		}
	}

	err = s.copyPath(ctx, source, sourcePath, destination, destinationPath, CopyOptions(options))
	if err != nil {
		return fmt.Errorf("failed to copy: %w", err)
	}

	err = os.RemoveAll(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to remove source after copy: %w", err)
	}

	return nil
}

func (s *LocalSession) Chmod(ctx context.Context, path Path, mode os.FileMode) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	localPath, err := s.requireLocalPath(OperationChmod, path)
	if err != nil {
		return err
	}

	return os.Chmod(localPath, mode)
}

func (s *LocalSession) Transfer(_ context.Context, request TransferRequest) (Transfer, error) {
	operation := request.Operation
	if operation == "" {
		operation = OperationCopy
	}

	return nil, NewUnsupportedError(ProviderLocal, operation, s.root.Path,
		"local provider transfer orchestration is handled outside provider-local sessions")
}

func (s *LocalSession) Close() error {
	return nil
}

func (s *LocalSession) requireLocalPath(operation Operation, path Path) (string, error) {
	if !path.IsLocal() {
		return "", NewUnsupportedError(ProviderLocal, operation, path,
			"local provider requires local filesystem paths")
	}
	return path.String(), nil
}

func (s *LocalSession) copyPath(
	ctx context.Context,
	source Path,
	sourcePath string,
	destination Path,
	destinationPath string,
	options CopyOptions,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	sourceInfo, err := os.Stat(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to stat source: %w", err)
	}

	if sourceInfo.IsDir() {
		if !options.Recursive {
			return NewUnsupportedError(ProviderLocal, OperationCopy, source,
				"copying directories requires recursive copy")
		}
		return s.copyDir(ctx, source, sourcePath, destination, destinationPath, sourceInfo, options)
	}

	return s.copyFile(ctx, destination, sourcePath, destinationPath, sourceInfo, options)
}

func (s *LocalSession) copyDir(
	ctx context.Context,
	_ Path,
	sourcePath string,
	destination Path,
	destinationPath string,
	sourceInfo os.FileInfo,
	options CopyOptions,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if !options.Overwrite {
		if _, err := os.Stat(destinationPath); err == nil {
			return NewConflictError(ProviderLocal, OperationCopy, destination, "destination already exists")
		} else if !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}

	err := os.MkdirAll(destinationPath, sourceInfo.Mode())
	if err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	entries, err := os.ReadDir(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to read source directory: %w", err)
	}

	for _, entry := range entries {
		entrySourcePath := filepath.Join(sourcePath, entry.Name())
		entryDestinationPath := filepath.Join(destinationPath, entry.Name())
		entrySource := NewLocalPath(entrySourcePath)
		entryDestination := NewLocalPath(entryDestinationPath)

		entryInfo, infoErr := entry.Info()
		if infoErr != nil {
			return fmt.Errorf("failed to get entry info: %w", infoErr)
		}

		if entryInfo.IsDir() {
			err = s.copyDir(
				ctx,
				entrySource,
				entrySourcePath,
				entryDestination,
				entryDestinationPath,
				entryInfo,
				options,
			)
		} else {
			err = s.copyFile(ctx, entryDestination, entrySourcePath, entryDestinationPath, entryInfo, options)
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *LocalSession) copyFile(
	ctx context.Context,
	destination Path,
	sourcePath string,
	_ string,
	sourceInfo os.FileInfo,
	options CopyOptions,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	reader, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer reader.Close()

	err = s.Create(ctx, destination, reader, CreateOptions{Mode: sourceInfo.Mode(), Overwrite: options.Overwrite})
	if err != nil {
		if errors.Is(err, ErrConflict) {
			return err
		}
		return fmt.Errorf("failed to create destination file: %w", err)
	}

	return nil
}

func (s *LocalSession) moveToTrash(localPath string, path Path) error {
	switch runtime.GOOS {
	case utils.OsDarwin:
		return s.Move(context.Background(), path,
			NewLocalPath(filepath.Join(variable.DarwinTrashDirectory, filepath.Base(localPath))),
			MoveOptions{Overwrite: true, Recursive: true})
	case utils.OsWindows:
		return trashWin.Throw(localPath)
	default:
		return trash.Trash(localPath)
	}
}

func newLocalStat(path string, info os.FileInfo) Stat {
	stat := Stat{
		Name:       info.Name(),
		Size:       info.Size(),
		Mode:       info.Mode(),
		ModTime:    info.ModTime(),
		IsDir:      info.IsDir(),
		IsSymlink:  info.Mode()&os.ModeSymlink != 0,
		ProviderID: string(ProviderLocal),
	}
	if stat.IsSymlink {
		if target, err := os.Readlink(path); err == nil {
			stat.Target = NewLocalPath(target)
		}
	}
	return stat
}

func localPathsSharePartition(path1, path2 string) (bool, error) {
	absPath1, err := filepath.Abs(path1)
	if err != nil {
		return false, fmt.Errorf("failed to get absolute path of the first path: %w", err)
	}

	absPath2, err := filepath.Abs(path2)
	if err != nil {
		return false, fmt.Errorf("failed to get absolute path of the second path: %w", err)
	}

	if runtime.GOOS == utils.OsWindows {
		return localDriveLetter(absPath1) == localDriveLetter(absPath2), nil
	}

	return filepath.VolumeName(absPath1) == filepath.VolumeName(absPath2), nil
}

func localDriveLetter(path string) string {
	return strings.ToUpper(string(path[0]))
}
