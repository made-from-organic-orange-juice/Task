package wpsapi

import (
	"golang.org/x/sys/windows"
)

type SystemSnapShot struct {
	Process ProcessInformation
	Modules Modules
}

func (SystemSnapShot) SystemSnapShotIterator() func() (SystemSnapShot, error) {
	proc, err := Processes{}.New()

	if err != nil {
		return func() (SystemSnapShot, error) {
			return SystemSnapShot{}, err
		}
	}

	nextProcess := proc.Iterator()

	return func() (SystemSnapShot, error) {
		process, err := nextProcess()

		if err != nil {
			return SystemSnapShot{}, err
		}

		// find a valid process!
		var h windows.Handle
		for {
			h, err = windows.OpenProcess(windows.PROCESS_QUERY_INFORMATION|windows.PROCESS_VM_READ, false, process.ID)
			if err != nil {
				process, err = nextProcess()
				// if we are at the end of the list
				if err != nil {
					return SystemSnapShot{}, err
				}
				// otherwise get the next item!
				continue
			}
			break
		}

		defer windows.CloseHandle(h)

		modules, err := Modules{}.New(h)
		if err != nil {
			return SystemSnapShot{}, err
		}

		return SystemSnapShot{
			Process: process,
			Modules: modules,
		}, nil

	}

}
