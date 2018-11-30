package main

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"syscall"
	"time"
	"unsafe"
)

func main() {
	Executable()


	log.Println(GetFileModTime("/Users/wang/go/src/github.com/lu566/go-examples/test-demo/xiaochuan.txt"))
}

func GetFileModTime(path string) int64 {
	f, err := os.Open(path)
	if err != nil {
		log.Println("open file error")
		return time.Now().Unix()
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		log.Println("stat fileinfo error")
		return time.Now().Unix()
	}

	return fi.ModTime().Unix()
}

func Executable() (string, error) {
	var mib [4]int32
	switch runtime.GOOS {
	case "freebsd":
		mib = [4]int32{1 /* CTL_KERN */, 14 /* KERN_PROC */, 12 /* KERN_PROC_PATHNAME */, -1}
	case "darwin":
		mib = [4]int32{1 /* CTL_KERN */, 38 /* KERN_PROCARGS */, int32(os.Getpid()), -1}
	}

	n := uintptr(0)
	// Get length.
	_, _, errNum := syscall.Syscall6(syscall.SYS___SYSCTL, uintptr(unsafe.Pointer(&mib[0])), 4, 0, uintptr(unsafe.Pointer(&n)), 0, 0)
	if errNum != 0 {
		return "", errNum
	}
	if n == 0 { // This shouldn't happen.
		return "", nil
	}
	buf := make([]byte, n)
	_, _, errNum = syscall.Syscall6(syscall.SYS___SYSCTL, uintptr(unsafe.Pointer(&mib[0])), 4, uintptr(unsafe.Pointer(&buf[0])), uintptr(unsafe.Pointer(&n)), 0, 0)
	if errNum != 0 {
		return "", errNum
	}
	if n == 0 { // This shouldn't happen.
		return "", nil
	}
	for i, v := range buf {
		if v == 0 {
			buf = buf[:i]
			break
		}
	}
	var err error
	execPath := string(buf)
	// execPath will not be empty due to above checks.
	// Try to get the absolute path if the execPath is not rooted.
	if execPath[0] != '/' {
		execPath, err = getAbs(execPath)
		if err != nil {
			return execPath, err
		}
	}
	// For darwin KERN_PROCARGS may return the path to a symlink rather than the
	// actual executable.
	if runtime.GOOS == "darwin" {
		if execPath, err = filepath.EvalSymlinks(execPath); err != nil {
			return execPath, err
		}
	}
	return execPath, nil
}


var initCwd, initCwdErr = os.Getwd()
func getAbs(execPath string) (string, error) {
	if initCwdErr != nil {
		return execPath, initCwdErr
	}
	// The execPath may begin with a "../" or a "./" so clean it first.
	// Join the two paths, trailing and starting slashes undetermined, so use
	// the generic Join function.
	return filepath.Join(initCwd, filepath.Clean(execPath)), nil
}
