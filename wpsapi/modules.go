package wpsapi

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

/**
 * * Modules is a list of all the modules for a specific process
 * windows.ModuleInfo contains these attributes:
 *  BaseOfDll   uintptr
 *	SizeOfImage uint32
 *	EntryPoint  uintptr
 */
type ModuleInformation struct {
	Path     string
	BaseName string
	Info     windows.ModuleInfo
}

type Modules struct {
	m []ModuleInformation
}

func (Modules) New(process windows.Handle) (mod Modules, err error) {

	// get the module
	const maxModules = 250
	module := [maxModules]windows.Handle{}
	var cbNeeded uint32
	var moduleInformation ModuleInformation
	const sizeModule = unsafe.Sizeof(module)
	err = windows.EnumProcessModulesEx(process, &module[0], uint32(sizeModule), &cbNeeded, windows.LIST_MODULES_DEFAULT)
	if err != nil {
		return
	}

	// calculates the amount of modules found
	n := cbNeeded / uint32(unsafe.Sizeof(&module[0]))

	// if the number of modules found is bigger than we expect
	// resize the module and try again!
	if n > maxModules {
		//TODO: resize the array and make a call to enumProcessModuleEx again with the right size!
		n = maxModules
	}

	for i := 0; i < int(n); i++ {
		// get module path.
		modulePath := make([]uint16, 512)
		err = windows.GetModuleFileNameEx(process, module[i], &modulePath[0], uint32(len(modulePath)))
		if err != nil {
			return
		}
		moduleInformation.Path = windows.UTF16ToString(modulePath)

		// get module base name
		err = windows.GetModuleBaseName(process, module[i], &modulePath[0], uint32(len(modulePath)))
		if err != nil {
			return
		}
		moduleInformation.BaseName = windows.UTF16ToString(modulePath)

		// get module information.
		var moduleInfo windows.ModuleInfo
		err = windows.GetModuleInformation(process, module[i], &moduleInfo, uint32(unsafe.Sizeof(moduleInfo)))
		if err != nil {
			return
		}
		moduleInformation.Info = moduleInfo

		mod.m = append(mod.m, moduleInformation)
	}

	return
}

/**
 *  * A closure for iteration through a process modules.
 *
 */

func (mod Modules) Iterator() func() (ModuleInformation, error) {
	pos := 0
	return func() (ModuleInformation, error) {
		if pos < len(mod.m) {
			m := mod.m[pos]
			pos++
			return m, nil
		}

		return ModuleInformation{}, ErrOutOfRange("At the end of Modules!")
	}
}
