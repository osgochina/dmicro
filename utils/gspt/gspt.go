package gspt

/*

#include "setproctitle.h"

*/
import "C"

import (
	"os"
	"unsafe"
)

const (
	// These values must match the return values for spt_init1() used in C.

	HaveNone        = 0
	HaveNative      = 1
	HaveReplacement = 2
)

var (
	HaveSetProcTitle     int
	HaveSetProcTitleFast int
)

func initialization() {
	HaveSetProcTitle = int(C.spt_init1())
	HaveSetProcTitleFast = int(C.spt_fast_init1())

	if HaveSetProcTitle == HaveReplacement {
		newArgs := make([]string, len(os.Args))
		for i, s := range os.Args {
			// Use cgo to force go to make copies of the strings.
			cs := C.CString(s)
			newArgs[i] = C.GoString(cs)
			C.free(unsafe.Pointer(cs))
		}
		os.Args = newArgs

		//env := os.Environ()
		//for _, kv := range env {
		//	skv := strings.SplitN(kv, "=", 2)
		//	_ = os.Setenv(skv[0], skv[1])
		//}

		argc := C.int(len(os.Args))
		arg0 := C.CString(os.Args[0])
		defer C.free(unsafe.Pointer(arg0))

		C.spt_init2(argc, arg0)

		// Restore the original title.
		//setProcTitle(os.Args[0])
	}
}

func setProcTitle(title string) {
	cs := C.CString(title)
	defer C.free(unsafe.Pointer(cs))

	C.spt_setproctitle(cs)
}

func SetProcTitle(title string) {
	initialization()
	setProcTitle(title)
}

func setProcTitleFast(title string) {
	cs := C.CString(title)
	defer C.free(unsafe.Pointer(cs))
	C.spt_setproctitle_fast(cs)
}

func SetProcTitleFast(title string) {
	initialization()
	setProcTitleFast(title)
}
