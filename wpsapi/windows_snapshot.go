package wpsapi

import (
	"syscall"
	"unsafe"

	"github.com/gonutz/w32/v2"
)

type WindowsProcess struct {
	Name string
	ID   uint32
}

type WindowsModules struct {
	m []w32.MODULEENTRY32
}

type ProcessEntry struct {
	Process WindowsProcess
	Modules WindowsModules
}

type SystemSnapshot []ProcessEntry

func (SystemSnapshot) New() (wm SystemSnapshot, err error) {

	// get all the processes!
	processes, ok := w32.EnumAllProcesses()

	if !ok {
		return SystemSnapshot{}, ErrModuleProcessing("could not retrieve any processes")
	}

	for _, processId := range processes {
		var module w32.MODULEENTRY32
		module.Size = uint32(unsafe.Sizeof(module))

		snap := w32.CreateToolhelp32Snapshot(w32.TH32CS_SNAPMODULE, processId)
		if snap == 0 {
			return SystemSnapshot{}, ErrModuleProcessing("could not create a snapshot of the system's modules")
		}

		var windowsProcess WindowsProcess
		var windowsModules WindowsModules

		if w32.Module32First(snap, &module) {
			name := syscall.UTF16ToString(module.SzModule[:])
			windowsProcess.Name = name
			windowsProcess.ID = processId
		} else {
			continue
		}

		for w32.Module32Next(snap, &module) {
			windowsModules.m = append(windowsModules.m, module)
		}

		wm = append(wm, ProcessEntry{
			Process: windowsProcess,
			Modules: windowsModules,
		})

	}

	return
}

func (wm SystemSnapshot) CountInstances() (instances map[string]int) {
	next := wm.Iterator()
	instances = make(map[string]int)
	for {
		snap, err := next()
		if err != nil {
			break
		}
		processName := snap.Process.Name
		if _, ok := instances[processName]; !ok {
			instances[processName] = 1
		} else {
			instances[processName]++
		}
	}
	return
}

func (wm SystemSnapshot) Iterator() func() (ProcessEntry, error) {
	pos := 0
	return func() (ProcessEntry, error) {
		if pos < len(wm) {
			m := wm[pos]
			pos++
			return m, nil
		}
		return ProcessEntry{}, ErrOutOfRange("no more modules to iterate")
	}
}
