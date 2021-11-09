package wpsapi

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

/**
 *  * set a privilege to the given token
 *
 */

func SetPrivilege(t windows.Token, privilege *uint16, enablePrivilege bool) error {

	var luid windows.LUID
	err := windows.LookupPrivilegeValue(nil, privilege, &luid)
	if err != nil {
		return err
	}

	var tp windows.Tokenprivileges

	tp.PrivilegeCount = 1
	tp.Privileges[0].Luid = luid
	if enablePrivilege {
		tp.Privileges[0].Attributes = windows.SE_PRIVILEGE_ENABLED
	} else {
		tp.Privileges[0].Attributes = 0
	}

	err = windows.AdjustTokenPrivileges(t, false, &tp, uint32(unsafe.Sizeof(tp)), nil, nil)
	if err != nil {
		return err
	}

	return nil
}

/**
 * * Set a priv for the current process
 */

func SetPrivilegeForCurrentProcess(p string) {

	var currentToken windows.Token
	currentProcess := windows.CurrentProcess()

	err := windows.OpenProcessToken(currentProcess, windows.TOKEN_ADJUST_PRIVILEGES, &currentToken)
	if err != nil {
		panic(err)
	}
	defer currentToken.Close()

	priv := windows.StringToUTF16Ptr(p)
	err = SetPrivilege(currentToken, priv, true)
	if err != nil {
		panic(err)
	}

}
