package filesystem

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	posixpath "path"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	transferProgressChunkBytes  int64 = 256 * 1024
	transferProgressBuffer            = 64
	transferCleanupTimeout            = 5 * time.Second
	transferChecksumBufferBytes       = 32 * 1024
)

var nextTransferID atomic.Uint64 //nolint:gochecknoglobals // Atomic process-wide transfer ID sequence.

type TransferEngine struct {
	resolver SessionResolver
}

func NewTransferEngine(resolver SessionResolver) *TransferEngine {
	return &TransferEngine{resolver: resolver}
}

func ValidateTransferTopology(source, destination Location) error {
	if source.Provider != ProviderLocal && destination.Provider != ProviderLocal &&
		source.SessionID != destination.SessionID {
		return NewUnsupportedError(
			source.Provider,
			OperationRemoteCrossSessionMove,
			source.Path,
			"cross-session remote to remote transfer is deferred for v1",
		)
	}
	if locationsShareFilesystem(source, destination) && pathsOverlap(source.Path, destination.Path, true) {
		return NewUnsupportedError(
			source.Provider,
			OperationCopy,
			destination.Path,
			"source and destination must not be the same path or nested within the source",
		)
	}
	return nil
}

func locationsShareFilesystem(source, destination Location) bool {
	if source.Provider == ProviderLocal && destination.Provider == ProviderLocal {
		return true
	}
	return source.Provider == destination.Provider && source.SessionID == destination.SessionID
}

func pathsOverlap(source, destination Path, rejectDescendant bool) bool {
	if source.Kind() != destination.Kind() {
		return false
	}
	if source.IsLocal() {
		sourcePath, sourceErr := canonicalLocalPath(source.String())
		destinationPath, destinationErr := canonicalLocalPath(destination.String())
		if sourceErr != nil || destinationErr != nil {
			return filepath.Clean(source.String()) == filepath.Clean(destination.String())
		}
		relative, err := filepath.Rel(sourcePath, destinationPath)
		return err == nil && (relative == "." || rejectDescendant && relative != ".." &&
			!strings.HasPrefix(relative, ".."+string(filepath.Separator)))
	}

	sourcePath := posixpath.Clean(source.String())
	destinationPath := posixpath.Clean(destination.String())
	if sourcePath == destinationPath {
		return true
	}
	return rejectDescendant && sourcePath != "/" && strings.HasPrefix(destinationPath, sourcePath+"/") ||
		rejectDescendant && sourcePath == "/" && strings.HasPrefix(destinationPath, "/")
}

func canonicalLocalPath(value string) (string, error) {
	absolute, err := filepath.Abs(value)
	if err != nil {
		return "", err
	}
	candidate := absolute
	suffix := make([]string, 0)
	for {
		resolved, resolveErr := filepath.EvalSymlinks(candidate)
		if resolveErr == nil {
			segments := make([]string, 1, len(suffix)+1)
			segments[0] = resolved
			for i := len(suffix) - 1; i >= 0; i-- {
				segments = append(segments, suffix[i])
			}
			return filepath.Join(segments...), nil
		}
		parent := filepath.Dir(candidate)
		if parent == candidate {
			return filepath.Clean(absolute), nil
		}
		suffix = append(suffix, filepath.Base(candidate))
		candidate = parent
	}
}

func InferTransferDirection(source, destination Location) TransferDirection {
	switch {
	case source.Provider == ProviderLocal && destination.Provider == ProviderLocal:
		return TransferLocal
	case source.Provider == ProviderLocal && destination.Provider != ProviderLocal:
		return TransferUpload
	case source.Provider != ProviderLocal && destination.Provider == ProviderLocal:
		return TransferDownload
	default:
		return TransferRemote
	}
}

func (e *TransferEngine) Start(ctx context.Context, request TransferRequest) (Transfer, error) {
	if e == nil || e.resolver == nil {
		return nil, errors.New("transfer engine requires a session resolver")
	}
	if err := ValidateTransferTopology(request.Source, request.Destination); err != nil {
		return nil, err
	}
	if request.Direction == "" {
		request.Direction = InferTransferDirection(request.Source, request.Destination)
	}
	if request.Operation == "" {
		request.Operation = inferTransferOperation(request)
	}

	transferCtx, cancel := context.WithCancel(ctx)
	t := &managedTransfer{
		id:       newTransferID(),
		request:  request,
		ctx:      transferCtx,
		cancel:   cancel,
		progress: make(chan Progress, transferProgressBuffer),
		done:     make(chan struct{}),
		resolver: e.resolver,
	}

	go t.run()
	return t, nil
}

type managedTransfer struct {
	id       TransferID
	request  TransferRequest
	ctx      context.Context
	cancel   context.CancelFunc
	resolver SessionResolver

	progress chan Progress
	done     chan struct{}

	mu  sync.RWMutex
	err error
}

func (t *managedTransfer) ID() TransferID {
	return t.id
}

func (t *managedTransfer) Operation() Operation {
	return t.request.Operation
}

func (t *managedTransfer) Direction() TransferDirection {
	return t.request.Direction
}

func (t *managedTransfer) Progress() <-chan Progress {
	return t.progress
}

func (t *managedTransfer) Cancel(_ context.Context) error {
	t.cancel()
	return nil
}

func (t *managedTransfer) Wait(ctx context.Context) error {
	select {
	case <-t.done:
		t.mu.RLock()
		defer t.mu.RUnlock()
		return t.err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (t *managedTransfer) run() {
	defer close(t.done)
	defer close(t.progress)

	sourceSession, err := t.resolver.ResolveSession(t.ctx, t.request.Source)
	if err != nil {
		t.finish(
			err,
			Progress{TransferID: t.id, Operation: t.request.Operation, Current: t.request.Source.Path, Err: err},
		)
		return
	}
	defer sourceSession.Close()

	destinationSession, err := t.resolver.ResolveSession(t.ctx, t.request.Destination)
	if err != nil {
		t.finish(
			err,
			Progress{TransferID: t.id, Operation: t.request.Operation, Current: t.request.Destination.Path, Err: err},
		)
		return
	}
	defer destinationSession.Close()

	plan, err := t.planSource(sourceSession)
	if err != nil {
		t.finish(
			err,
			Progress{TransferID: t.id, Operation: t.request.Operation, Current: t.request.Source.Path, Err: err},
		)
		return
	}

	t.emit(Progress{
		TransferID: t.id,
		Operation:  t.request.Operation,
		Current:    t.request.Source.Path,
		Total:      plan.totalFiles,
		BytesTotal: plan.totalBytes,
	})

	state := &transferState{transfer: t, totalFiles: plan.totalFiles, totalBytes: plan.totalBytes}
	err = t.transferPath(sourceSession, destinationSession, t.request.Source, t.request.Destination, plan.stat, state)
	if err == nil && t.request.Operation == OperationCutMove {
		err = sourceSession.Delete(t.ctx, t.request.Source.Path, DeleteOptions{Recursive: true})
	}

	if err != nil {
		t.finish(err, Progress{
			TransferID: t.id,
			Operation:  t.request.Operation,
			Current:    state.current,
			Done:       state.doneFiles,
			Total:      state.totalFiles,
			BytesDone:  state.doneBytes,
			BytesTotal: state.totalBytes,
			Err:        err,
		})
		return
	}

	t.finish(nil, Progress{
		TransferID: t.id,
		Operation:  t.request.Operation,
		Current:    state.current,
		Done:       state.doneFiles,
		Total:      state.totalFiles,
		BytesDone:  state.doneBytes,
		BytesTotal: state.totalBytes,
	})
}

func (t *managedTransfer) finish(err error, progress Progress) {
	t.mu.Lock()
	t.err = err
	t.mu.Unlock()
	t.emit(progress)
}

func (t *managedTransfer) emit(progress Progress) {
	select {
	case t.progress <- progress:
	default:
	}
}

type transferPlan struct {
	stat       Stat
	totalFiles int64
	totalBytes int64
}

type transferState struct {
	transfer       *managedTransfer
	totalFiles     int64
	totalBytes     int64
	doneFiles      int64
	doneBytes      int64
	current        Path
	lastProgressAt int64
}

func (s *transferState) setCurrent(path Path) {
	s.current = path
}

func (s *transferState) addBytes(path Path, bytes int64) {
	if bytes <= 0 {
		return
	}
	s.doneBytes += bytes
	s.current = path
	if s.doneBytes-s.lastProgressAt < transferProgressChunkBytes && s.doneBytes != s.totalBytes {
		return
	}
	s.lastProgressAt = s.doneBytes
	s.transfer.emit(Progress{
		TransferID: s.transfer.id,
		Operation:  s.transfer.request.Operation,
		Current:    path,
		Done:       s.doneFiles,
		Total:      s.totalFiles,
		BytesDone:  s.doneBytes,
		BytesTotal: s.totalBytes,
	})
}

func (s *transferState) completeFile(path Path) {
	s.doneFiles++
	s.current = path
	s.transfer.emit(Progress{
		TransferID: s.transfer.id,
		Operation:  s.transfer.request.Operation,
		Current:    path,
		Done:       s.doneFiles,
		Total:      s.totalFiles,
		BytesDone:  s.doneBytes,
		BytesTotal: s.totalBytes,
	})
}

func (t *managedTransfer) planSource(session Session) (transferPlan, error) {
	stat, err := session.Stat(t.ctx, t.request.Source.Path)
	if err != nil {
		return transferPlan{}, err
	}
	files, bytes, err := countTransferTotals(t.ctx, session, t.request.Source.Path, stat)
	if err != nil {
		return transferPlan{}, err
	}
	return transferPlan{stat: stat, totalFiles: files, totalBytes: bytes}, nil
}

func countTransferTotals(ctx context.Context, session Session, path Path, stat Stat) (int64, int64, error) {
	if err := ctx.Err(); err != nil {
		return 0, 0, err
	}
	if !stat.IsDir || stat.IsSymlink {
		return 1, stat.Size, nil
	}
	entries, err := session.List(ctx, path)
	if err != nil {
		return 0, 0, err
	}
	var totalFiles int64
	var totalBytes int64
	for _, entry := range entries {
		files, bytes, err := countTransferTotals(ctx, session, entry.Path, entry.Stat)
		if err != nil {
			return 0, 0, err
		}
		totalFiles += files
		totalBytes += bytes
	}
	return totalFiles, totalBytes, nil
}

func (t *managedTransfer) transferPath(
	sourceSession Session,
	destinationSession Session,
	source Location,
	destination Location,
	stat Stat,
	state *transferState,
) error {
	if err := t.ctx.Err(); err != nil {
		return normalizeTransferError(destinationSession.Provider(), t.request.Operation, destination.Path, err)
	}
	if stat.IsDir && !stat.IsSymlink {
		return t.transferDirectory(sourceSession, destinationSession, source, destination, stat, state)
	}
	return t.transferFile(sourceSession, destinationSession, source, destination, stat, state)
}

func (t *managedTransfer) transferDirectory(
	sourceSession Session,
	destinationSession Session,
	source Location,
	destination Location,
	stat Stat,
	state *transferState,
) error {
	exists, destinationStat, err := t.destinationState(destinationSession, destination)
	if err != nil {
		return err
	}
	if exists {
		if !destinationStat.IsDir || destinationStat.IsSymlink {
			return NewConflictError(
				destinationSession.Provider(),
				t.request.Operation,
				destination.Path,
				"destination already exists",
			)
		}
		if t.request.Overwrite {
			return NewUnsupportedError(destinationSession.Provider(), t.request.Operation, destination.Path,
				"atomic directory replacement is unavailable for overwrite transfers")
		}
		return NewConflictError(
			destinationSession.Provider(),
			t.request.Operation,
			destination.Path,
			"destination already exists",
		)
	}

	tempDestination := destination
	tempDestination.Path = tempSiblingPath(destination.Path)
	cleanup := tempDestination
	cleanupRecursive := true
	defer func() {
		if err != nil {
			_ = t.cleanupLocation(cleanup, cleanupRecursive)
		}
	}()

	if err = destinationSession.Mkdir(
		t.ctx,
		tempDestination.Path,
		MkdirOptions{Mode: stat.Mode, Parents: true},
	); err != nil {
		return err
	}

	entries, err := sourceSession.List(t.ctx, source.Path)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		entrySource := source
		entrySource.Path = entry.Path
		entryDestination := tempDestination
		entryDestination.Path = joinTransferPath(tempDestination.Path, entry.Name)
		if transferErr := t.transferPath(
			sourceSession,
			destinationSession,
			entrySource,
			entryDestination,
			entry.Stat,
			state,
		); transferErr != nil {
			err = transferErr
			return err
		}
	}

	err = destinationSession.Rename(t.ctx, tempDestination.Path, destination.Path, RenameOptions{})
	return err
}

//nolint:nonamedreturns // The named error drives deferred temporary-file cleanup.
func (t *managedTransfer) transferFile(
	sourceSession Session,
	destinationSession Session,
	source Location,
	destination Location,
	stat Stat,
	state *transferState,
) (err error) {
	state.setCurrent(source.Path)

	exists, _, err := t.destinationState(destinationSession, destination)
	if err != nil {
		return err
	}
	if exists && !t.request.Overwrite {
		return NewConflictError(
			destinationSession.Provider(),
			t.request.Operation,
			destination.Path,
			"destination already exists",
		)
	}

	tempDestination := destination
	tempDestination.Path = tempSiblingPath(destination.Path)
	defer func() {
		if err != nil {
			_ = t.cleanupLocation(tempDestination, false)
		}
	}()

	sourceReader, err := sourceSession.Read(t.ctx, source.Path)
	if err != nil {
		return err
	}
	defer sourceReader.Close()

	hasher := sha256.New()
	progressReader := &verifyingProgressReader{
		ctx:       t.ctx,
		path:      source.Path,
		provider:  sourceSession.Provider(),
		operation: t.request.Operation,
		reader:    sourceReader,
		hasher:    hasher,
		state:     state,
	}

	err = destinationSession.Create(t.ctx, tempDestination.Path, progressReader, CreateOptions{
		Mode:      stat.Mode,
		Overwrite: true,
	})
	if err != nil {
		return normalizeTransferError(destinationSession.Provider(), t.request.Operation, tempDestination.Path, err)
	}

	sourceChecksum := hex.EncodeToString(hasher.Sum(nil))
	destinationChecksum, err := checksumSessionPath(
		t.ctx,
		destinationSession,
		tempDestination.Path,
		t.request.Operation,
	)
	if err != nil {
		return err
	}
	if sourceChecksum != destinationChecksum {
		return NewConflictError(destinationSession.Provider(), t.request.Operation, destination.Path,
			"transferred content verification failed")
	}

	renameOptions := RenameOptions{}
	if exists {
		renameOptions.Overwrite = true
	}
	err = destinationSession.Rename(t.ctx, tempDestination.Path, destination.Path, renameOptions)
	if err != nil {
		return normalizeAtomicReplaceError(
			destinationSession.Provider(),
			t.request.Operation,
			destination.Path,
			err,
			exists,
		)
	}

	state.completeFile(source.Path)
	return nil
}

func (t *managedTransfer) destinationState(session Session, destination Location) (bool, Stat, error) {
	stat, err := session.Stat(t.ctx, destination.Path)
	if err == nil {
		return true, stat, nil
	}
	if errors.Is(err, ErrNotFound) || errors.Is(err, os.ErrNotExist) {
		return false, Stat{}, nil
	}
	return false, Stat{}, err
}

func (t *managedTransfer) cleanupLocation(location Location, recursive bool) error {
	cleanupCtx, cancel := context.WithTimeout(context.Background(), transferCleanupTimeout)
	defer cancel()

	deleteWith := func(session Session) error {
		deleteErr := session.Delete(cleanupCtx, location.Path, DeleteOptions{Recursive: recursive})
		if deleteErr == nil || errors.Is(deleteErr, ErrNotFound) {
			return nil
		}
		return deleteErr
	}

	resolver := t.resolver.ResolveSession
	if freshResolver, ok := t.resolver.(FreshSessionResolver); ok {
		resolver = freshResolver.ResolveFreshSession
	}

	session, err := resolver(cleanupCtx, location)
	if err != nil {
		return err
	}
	defer session.Close()
	return deleteWith(session)
}

func checksumSessionPath(ctx context.Context, session Session, path Path, operation Operation) (string, error) {
	reader, err := session.Read(ctx, path)
	if err != nil {
		return "", err
	}
	defer reader.Close()

	hasher := sha256.New()
	buffer := make([]byte, transferChecksumBufferBytes)
	for {
		if err = ctx.Err(); err != nil {
			return "", normalizeTransferError(session.Provider(), operation, path, err)
		}
		var read int
		read, err = reader.Read(buffer)
		if read > 0 {
			_, _ = hasher.Write(buffer[:read])
		}
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return "", normalizeTransferError(session.Provider(), operation, path, err)
		}
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

type verifyingProgressReader struct {
	ctx       context.Context
	path      Path
	provider  ProviderKind
	operation Operation
	reader    io.Reader
	hasher    hashWriter
	state     *transferState
}

type hashWriter interface {
	Write([]byte) (int, error)
	Sum([]byte) []byte
}

func (r *verifyingProgressReader) Read(p []byte) (int, error) {
	if err := r.ctx.Err(); err != nil {
		return 0, normalizeTransferError(r.provider, r.operation, r.path, err)
	}
	n, err := r.reader.Read(p)
	if n > 0 {
		_, _ = r.hasher.Write(p[:n])
		r.state.addBytes(r.path, int64(n))
	}
	if err != nil && !errors.Is(err, io.EOF) {
		return n, normalizeTransferError(r.provider, r.operation, r.path, err)
	}
	return n, err
}

func inferTransferOperation(request TransferRequest) Operation {
	switch request.Direction {
	case TransferLocal:
		// Local transfers keep the requested operation below.
	case TransferUpload:
		return OperationTransferLocalToRemote
	case TransferDownload:
		return OperationTransferRemoteToLocal
	case TransferRemote:
		if request.Operation == OperationCutMove {
			return OperationRemoteSameSessionMove
		}
	}
	if request.Operation != "" {
		return request.Operation
	}
	return OperationCopy
}

func normalizeTransferError(provider ProviderKind, operation Operation, path Path, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) || errors.Is(err, ErrCanceled) {
		return NewCanceledError(provider, operation, path, err.Error())
	}
	if errors.Is(err, ErrDisconnected) {
		return NewDisconnectedError(provider, operation, path, err.Error())
	}
	if errors.Is(err, ErrConflict) || errors.Is(err, ErrPermission) || errors.Is(err, ErrNotFound) ||
		errors.Is(err, ErrUnsupported) {
		return err
	}
	if errors.Is(err, os.ErrNotExist) {
		return NewNotFoundError(provider, operation, path, err.Error())
	}
	if errors.Is(err, os.ErrPermission) {
		return NewPermissionError(provider, operation, path, err.Error())
	}
	return err
}

func normalizeAtomicReplaceError(
	provider ProviderKind,
	operation Operation,
	path Path,
	err error,
	destinationExists bool,
) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, ErrUnsupported) || errors.Is(err, ErrConflict) {
		return err
	}
	if destinationExists && provider == ProviderLocal && runtime.GOOS == "windows" {
		return NewUnsupportedError(
			provider,
			operation,
			path,
			"atomic overwrite rename is unavailable for this provider",
		)
	}
	return normalizeTransferError(provider, operation, path, err)
}

func tempSiblingPath(path Path) Path {
	base := transferBaseName(path)
	tempName := fmt.Sprintf(".%s.superfile-transfer-%d", base, time.Now().UnixNano())
	if path.IsRemote() {
		return path.Dir().Join(tempName)
	}
	return NewLocalPath(filepath.Join(filepath.Dir(path.String()), tempName))
}

func transferBaseName(path Path) string {
	if path.IsRemote() {
		return path.Base()
	}
	base := filepath.Base(path.String())
	if base == "." || base == string(filepath.Separator) || base == "" {
		return "root"
	}
	return base
}

func joinTransferPath(base Path, name string) Path {
	if base.IsRemote() {
		return base.Join(name)
	}
	return NewLocalPath(filepath.Join(base.String(), name))
}

func newTransferID() TransferID {
	return TransferID(fmt.Sprintf("transfer-%d", nextTransferID.Add(1)))
}
