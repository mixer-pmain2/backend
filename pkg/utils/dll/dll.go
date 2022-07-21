package dll

import "syscall"

type lib struct {
	dll *syscall.DLL
}

func Open(name string) *lib {
	dll := syscall.MustLoadDLL(name)

	return &lib{
		dll: dll,
	}
}

func (l *lib) Free() {
	syscall.FreeLibrary(l.dll.Handle)
}

func (l *lib) Proc(name string) *syscall.Proc {
	return l.dll.MustFindProc(name)
}
