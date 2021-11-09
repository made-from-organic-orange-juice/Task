package wpsapi

import (
	"log"
	"unsafe"

	"golang.org/x/sys/windows"
)

type ProcessInformation struct {
	Name            string
	ID              uint32
	UnderlyingEntry windows.ProcessEntry32
}

type Processes struct {
	p []ProcessInformation
}

func (Processes) New() (proc Processes, err error) {

	// take a snapshot of all processes in the system
	h, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)

	if err != nil {
		log.Default().Println("In (Processes) New() line 28: ")
		return
	}

	// make sure we clean the snapshot object at the end
	defer windows.CloseHandle(h)

	// Set the size of the structure before using it
	const size = unsafe.Sizeof(windows.ProcessEntry32{})
	processEntry := windows.ProcessEntry32{Size: uint32(size)}

	// retrive information about the first process
	err = windows.Process32First(h, &processEntry)
	if err != nil {
		log.Default().Println("In (Processes) New() line 42: ")
		return
	}

	// add the first process
	proc.p = append(proc.p, ProcessInformation{
		Name:            windows.UTF16ToString(processEntry.ExeFile[:]),
		ID:              processEntry.ProcessID,
		UnderlyingEntry: processEntry,
	})

	// Now walk the snapshot of processes, and add the rest to the list
	for {
		//
		err = windows.Process32Next(h, &processEntry)
		if err != nil {
			break
		}

		proc.p = append(proc.p, ProcessInformation{
			Name:            windows.UTF16ToString(processEntry.ExeFile[:]),
			ID:              processEntry.ProcessID,
			UnderlyingEntry: processEntry,
		})

	}

	return proc, nil
}

/**
 *  * A closure for iteration through the systems processes
 */

func (proc Processes) Iterator() func() (ProcessInformation, error) {
	pos := 0
	return func() (ProcessInformation, error) {
		if pos < len(proc.p) {
			p := proc.p[pos]
			pos++
			return p, nil
		}
		return ProcessInformation{}, ErrOutOfRange("At the end of Processes!")
	}
}
