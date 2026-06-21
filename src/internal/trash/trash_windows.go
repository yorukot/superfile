//go:build windows

package trash

import (
	"fmt"
	"path/filepath"
	"runtime"
	"syscall"
	"unsafe"
)

type guid struct {
	data1 uint32
	data2 uint16
	data3 uint16
	data4 [8]byte
}

type comObject struct {
	vtable *uintptr
}

type fileOperationProgressSink struct {
	vtable   *fileOperationProgressSinkVtbl
	refs     uintptr
	recycled bool
	failedHR uintptr
}

type fileOperationProgressSinkVtbl struct {
	queryInterface uintptr
	addRef         uintptr
	release        uintptr
	start          uintptr
	finish         uintptr
	preRename      uintptr
	postRename     uintptr
	preMove        uintptr
	postMove       uintptr
	preCopy        uintptr
	postCopy       uintptr
	preDelete      uintptr
	postDelete     uintptr
	preNew         uintptr
	postNew        uintptr
	updateProgress uintptr
	resetTimer     uintptr
	pauseTimer     uintptr
	resumeTimer    uintptr
}

const (
	coinitApartmentThreaded = 0x2
	clsctxInprocServer      = 0x1

	sOK          = 0x00000000
	eNoInterface = 0x80004002
	ePointer     = 0x80004003

	fofSilent         = 0x0004
	fofNoConfirmation = 0x0010
	fofAllowUndo      = 0x0040
	fofNoConfirmMkdir = 0x0200
	fofNoErrorUI      = 0x0400

	fofxRecycleOnDelete = 0x00080000
	fofxEarlyFailure    = 0x00100000
	fofxAddUndoRecord   = 0x20000000
)

var (
	ole32   = syscall.NewLazyDLL("ole32.dll")
	shell32 = syscall.NewLazyDLL("shell32.dll")

	procCoInitializeEx              = ole32.NewProc("CoInitializeEx")
	procCoUninitialize              = ole32.NewProc("CoUninitialize")
	procCoCreateInstance            = ole32.NewProc("CoCreateInstance")
	procSHCreateItemFromParsingName = shell32.NewProc("SHCreateItemFromParsingName")

	clsidFileOperation = guid{
		0x3ad05575, 0x8857, 0x4850,
		[8]byte{0x92, 0x77, 0x11, 0xb8, 0x5b, 0xdb, 0x8e, 0x09},
	}
	iidIUnknown = guid{
		0x00000000, 0x0000, 0x0000,
		[8]byte{0xc0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x46},
	}
	iidIFileOperation = guid{
		0x947aab5f, 0x0a5c, 0x4c13,
		[8]byte{0xb4, 0xd6, 0x4b, 0xf7, 0x83, 0x6f, 0xc9, 0xf8},
	}
	iidIFileOperationProgressSink = guid{
		0x04b0f1a7, 0x9490, 0x44bc,
		[8]byte{0x96, 0xe1, 0x42, 0x96, 0xa3, 0x12, 0x52, 0xe2},
	}
	iidIShellItem = guid{
		0x43826d1e, 0xe718, 0x42ee,
		[8]byte{0xbc, 0x55, 0xa1, 0xe2, 0x61, 0xc3, 0x7b, 0xfe},
	}

	progressSinkVtbl = fileOperationProgressSinkVtbl{
		queryInterface: syscall.NewCallback(progressSinkQueryInterface),
		addRef:         syscall.NewCallback(progressSinkAddRef),
		release:        syscall.NewCallback(progressSinkRelease),
		start:          syscall.NewCallback(progressSinkStartOperations),
		finish:         syscall.NewCallback(progressSinkFinishOperations),
		preRename:      syscall.NewCallback(progressSinkPreRenameItem),
		postRename:     syscall.NewCallback(progressSinkPostRenameItem),
		preMove:        syscall.NewCallback(progressSinkPreMoveItem),
		postMove:       syscall.NewCallback(progressSinkPostMoveItem),
		preCopy:        syscall.NewCallback(progressSinkPreCopyItem),
		postCopy:       syscall.NewCallback(progressSinkPostCopyItem),
		preDelete:      syscall.NewCallback(progressSinkPreDeleteItem),
		postDelete:     syscall.NewCallback(progressSinkPostDeleteItem),
		preNew:         syscall.NewCallback(progressSinkPreNewItem),
		postNew:        syscall.NewCallback(progressSinkPostNewItem),
		updateProgress: syscall.NewCallback(progressSinkUpdateProgress),
		resetTimer:     syscall.NewCallback(progressSinkResetTimer),
		pauseTimer:     syscall.NewCallback(progressSinkPauseTimer),
		resumeTimer:    syscall.NewCallback(progressSinkResumeTimer),
	}
)

func Init() error {
	return nil
}

func Available(path string) bool {
	return path != ""
}

func Move(path string) (Result, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return Result{OriginalPath: path, Backend: BackendWindows}, err
	}

	result := Result{
		OriginalPath: absPath,
		Backend:      BackendWindows,
	}

	errCh := make(chan error, 1)
	go func() {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
		errCh <- recycleWithIFileOperation(absPath)
	}()
	if err := <-errCh; err != nil {
		return result, err
	}
	result.StrictlyRecycled = true
	return result, nil
}

func recycleWithIFileOperation(path string) error {
	hr, _, _ := procCoInitializeEx.Call(0, coinitApartmentThreaded)
	if failed(hr) {
		return hresultError("CoInitializeEx", hr)
	}
	defer procCoUninitialize.Call()

	var fileOperation *comObject
	hr, _, _ = procCoCreateInstance.Call(
		uintptr(unsafe.Pointer(&clsidFileOperation)),
		0,
		clsctxInprocServer,
		uintptr(unsafe.Pointer(&iidIFileOperation)),
		uintptr(unsafe.Pointer(&fileOperation)),
	)
	if failed(hr) {
		return hresultError("CoCreateInstance(IFileOperation)", hr)
	}
	defer fileOperation.release()

	flags := uintptr(fofSilent | fofNoConfirmation | fofAllowUndo | fofNoConfirmMkdir |
		fofNoErrorUI | fofxRecycleOnDelete | fofxAddUndoRecord | fofxEarlyFailure)
	if hr := fileOperation.call(5, flags); failed(hr) {
		return hresultError("IFileOperation.SetOperationFlags", hr)
	}

	item, err := shellItemFromPath(path)
	if err != nil {
		return err
	}
	defer item.release()

	sink := newProgressSink()
	if hr := fileOperation.call(18, uintptr(unsafe.Pointer(item)), uintptr(unsafe.Pointer(sink))); failed(hr) {
		return hresultError("IFileOperation.DeleteItem", hr)
	}
	if hr := fileOperation.call(21); failed(hr) {
		return hresultError("IFileOperation.PerformOperations", hr)
	}
	var aborted int32
	if hr := fileOperation.call(22, uintptr(unsafe.Pointer(&aborted))); failed(hr) {
		return hresultError("IFileOperation.GetAnyOperationsAborted", hr)
	}
	if aborted != 0 {
		return fmt.Errorf("Recycle Bin operation was aborted")
	}
	if sink.failedHR != 0 {
		return hresultError("IFileOperation.PostDeleteItem", sink.failedHR)
	}
	if !sink.recycled {
		return fmt.Errorf("Recycle Bin operation did not return a recycled item")
	}
	runtime.KeepAlive(sink)
	return nil
}

func shellItemFromPath(path string) (*comObject, error) {
	utf16Path, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return nil, err
	}
	var item *comObject
	hr, _, _ := procSHCreateItemFromParsingName.Call(
		uintptr(unsafe.Pointer(utf16Path)),
		0,
		uintptr(unsafe.Pointer(&iidIShellItem)),
		uintptr(unsafe.Pointer(&item)),
	)
	if failed(hr) {
		return nil, hresultError("SHCreateItemFromParsingName", hr)
	}
	return item, nil
}

func (obj *comObject) call(method uintptr, args ...uintptr) uintptr {
	vtable := *(**[64]uintptr)(unsafe.Pointer(obj))
	callArgs := append([]uintptr{uintptr(unsafe.Pointer(obj))}, args...)
	hr, _, _ := syscall.SyscallN(vtable[method], callArgs...)
	return hr
}

func (obj *comObject) release() {
	if obj != nil {
		_ = obj.call(2)
	}
}

func failed(hr uintptr) bool {
	return int32(hr&0xffffffff) < 0
}

func hresultError(operation string, hr uintptr) error {
	return fmt.Errorf("%s failed: HRESULT 0x%08X", operation, uint32(hr))
}

func newProgressSink() *fileOperationProgressSink {
	return &fileOperationProgressSink{
		vtable: &progressSinkVtbl,
		refs:   1,
	}
}

func progressSinkFromThis(this uintptr) *fileOperationProgressSink {
	return (*fileOperationProgressSink)(unsafe.Pointer(this))
}

func progressSinkQueryInterface(this uintptr, iid uintptr, object uintptr) uintptr {
	if object == 0 {
		return ePointer
	}
	*(*uintptr)(unsafe.Pointer(object)) = 0
	if iid == 0 {
		return ePointer
	}

	requested := *(*guid)(unsafe.Pointer(iid))
	if requested != iidIUnknown && requested != iidIFileOperationProgressSink {
		return eNoInterface
	}

	*(*uintptr)(unsafe.Pointer(object)) = this
	_ = progressSinkAddRef(this)
	return sOK
}

func progressSinkAddRef(this uintptr) uintptr {
	sink := progressSinkFromThis(this)
	sink.refs++
	return sink.refs
}

func progressSinkRelease(this uintptr) uintptr {
	sink := progressSinkFromThis(this)
	if sink.refs > 0 {
		sink.refs--
	}
	return sink.refs
}

func progressSinkStartOperations(_ uintptr) uintptr {
	return sOK
}

func progressSinkFinishOperations(_ uintptr, _ uintptr) uintptr {
	return sOK
}

func progressSinkPreRenameItem(_ uintptr, _ uintptr, _ uintptr, _ uintptr) uintptr {
	return sOK
}

func progressSinkPostRenameItem(_ uintptr, _ uintptr, _ uintptr, _ uintptr, _ uintptr, _ uintptr) uintptr {
	return sOK
}

func progressSinkPreMoveItem(_ uintptr, _ uintptr, _ uintptr, _ uintptr, _ uintptr) uintptr {
	return sOK
}

func progressSinkPostMoveItem(_ uintptr, _ uintptr, _ uintptr, _ uintptr, _ uintptr, _ uintptr, _ uintptr) uintptr {
	return sOK
}

func progressSinkPreCopyItem(_ uintptr, _ uintptr, _ uintptr, _ uintptr, _ uintptr) uintptr {
	return sOK
}

func progressSinkPostCopyItem(_ uintptr, _ uintptr, _ uintptr, _ uintptr, _ uintptr, _ uintptr, _ uintptr) uintptr {
	return sOK
}

func progressSinkPreDeleteItem(_ uintptr, _ uintptr, _ uintptr) uintptr {
	return sOK
}

func progressSinkPostDeleteItem(this uintptr, _ uintptr, _ uintptr, hrDelete uintptr, newlyCreated uintptr) uintptr {
	sink := progressSinkFromThis(this)
	if failed(hrDelete) {
		sink.failedHR = hrDelete
		return sOK
	}
	sink.recycled = newlyCreated != 0
	return sOK
}

func progressSinkPreNewItem(_ uintptr, _ uintptr, _ uintptr, _ uintptr) uintptr {
	return sOK
}

func progressSinkPostNewItem(
	_ uintptr,
	_ uintptr,
	_ uintptr,
	_ uintptr,
	_ uintptr,
	_ uintptr,
	_ uintptr,
	_ uintptr,
) uintptr {
	return sOK
}

func progressSinkUpdateProgress(_ uintptr, _ uintptr, _ uintptr) uintptr {
	return sOK
}

func progressSinkResetTimer(_ uintptr) uintptr {
	return sOK
}

func progressSinkPauseTimer(_ uintptr) uintptr {
	return sOK
}

func progressSinkResumeTimer(_ uintptr) uintptr {
	return sOK
}
