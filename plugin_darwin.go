package vst2

//#cgo darwin LDFLAGS: -framework CoreFoundation
//#include <CoreFoundation/CoreFoundation.h>
//#include "vst2.h"
/*
char * MYCFStringCopyUTF8String(CFStringRef aString) {
	if (aString == NULL) {
    	return NULL;
	}

  	CFIndex length = CFStringGetLength(aString);
  	CFIndex maxSize = CFStringGetMaximumSizeForEncoding(length, kCFStringEncodingUTF8) + 1;
	char *buffer = (char *)malloc(maxSize);
	if (CFStringGetCString(aString, buffer, maxSize, kCFStringEncodingUTF8)) {
    	return buffer;
	}
	free(buffer); // If we failed
	return NULL;
}
*/
import "C"
import (
	"fmt"
	"unsafe"
)

const (
	// Extension of Vst2 files
	Extension = ".vst"
)

var (
	// ScanPaths of Vst2 files
	ScanPaths []string
)

func init() {
	ScanPaths = []string{
		"~/Library/Audio/Plug-Ins/VST",
		"/Library/Audio/Plug-Ins/VST",
	}
}

// Library used to instantiate new instances of plugin
type Library struct {
	entryPoint unsafe.Pointer
	library    uintptr
	Name       string
	Path       string
}

func (l *Library) load() error {
	//create C string
	cpath := C.CString(l.Path)
	defer C.free(unsafe.Pointer(cpath))
	//convert to CF string
	cfpath := C.CFStringCreateWithCString(0, cpath, C.kCFStringEncodingUTF8)
	defer C.CFRelease(C.CFTypeRef(cfpath))

	//get bundle url
	bundleURL := C.CFURLCreateWithFileSystemPath(C.kCFAllocatorDefault, cfpath, C.kCFURLPOSIXPathStyle, C.true)
	if bundleURL == 0 {
		return fmt.Errorf("Failed to create bundle url at %v", l.Path)
	}
	defer C.CFRelease(C.CFTypeRef(bundleURL))
	//open bundle and release it only if it failed
	bundle := C.CFBundleCreate(C.kCFAllocatorDefault, bundleURL)
	l.library = uintptr(bundle)
	//bundle ref should be released in the end of program with plugin.unload call

	//create C string
	cvstMain := C.CString(vstMain)
	defer C.free(unsafe.Pointer(cvstMain))
	//create CF string
	cfvstMain := C.CFStringCreateWithCString(0, cvstMain, C.kCFStringEncodingUTF8)
	defer C.CFRelease(C.CFTypeRef(cfvstMain))

	l.entryPoint = unsafe.Pointer(C.CFBundleGetFunctionPointerForName(bundle, cfvstMain))
	if l.entryPoint == nil {
		l.Close()
		return fmt.Errorf("Failed to find entry point in bundle %v", l.Path)
	}
	// l.Name = getBundleString(bundle, "CFBundleName")

	return nil
}

//Close cleans up library refs
//TODO: exceptions handling
func (l *Library) Close() error {
	C.CFRelease(C.CFTypeRef(C.CFBundleRef(l.library)))
	l.library = 0
	return nil
}

//get string from CFBundle
func getBundleString(bundle C.CFBundleRef, str string) string {
	//create C string
	cstring := C.CString(str)
	defer C.free(unsafe.Pointer(cstring))
	//convert to CF string
	cfstring := C.CFStringCreateWithCString(0, cstring, C.kCFStringEncodingUTF8)
	defer C.CFRelease(C.CFTypeRef(cfstring))

	bundleString := C.CFStringRef(C.CFBundleGetValueForInfoDictionaryKey(bundle, cfstring))
	defer C.CFRelease(C.CFTypeRef(bundleString))

	convertedString := C.MYCFStringCopyUTF8String(bundleString)
	defer C.free(unsafe.Pointer(convertedString))
	return C.GoString(convertedString)
}
