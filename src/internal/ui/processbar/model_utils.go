package processbar

import (
	"sort"

	"github.com/lithammer/shortuuid"
)

func (m *Model) cntProcesses() int {
	return len(m.processes)
}

func (m *Model) isValid() bool {
	return m.renderIndex <= m.cursor &&
		m.cursor <= m.renderIndex+cntRenderableProcess(m.height-borderSize)-1
}

func (m *Model) viewHeight() int {
	return m.height - borderSize
}

func (m *Model) viewWidth() int {
	return m.width - borderSize
}

func (m *Model) getSortedProcesses() []Process {
	// save process in the array and sort the process by finished or not,
	// completion percetage, or finish time
	// TODO : This is very inefficient and can be improved.
	// The whole design needs to be changed so that we dont need to recreate the slice
	// and sort on each render. Idea : Maintain two slices - completed, ongoing
	// Processes should be added / removed to the slice on correct time, and we dont
	// need to redo slice formation and sorting on each render.
	// TODO : One idea is that we can use google/btree to store processes
	// have process implement a Less() method, and we can do O(logn) inserts/deletes
	// To make sure its always stored in an order we want. And then iterate in O(n)
	// in render()
	var processes []Process
	for _, p := range m.processes {
		processes = append(processes, p)
	}
	// sort by the process
	sort.Slice(processes, func(i, j int) bool {
		doneI := (processes[i].State == Successful || processes[i].State == Failed)
		doneJ := (processes[j].State == Successful || processes[j].State == Failed)

		// sort by done or not
		if doneI != doneJ {
			return !doneI
		}

		// if both not done
		if !doneI {
			completionI := float64(processes[i].Done) / float64(processes[i].Total)
			completionJ := float64(processes[j].Done) / float64(processes[j].Total)
			return completionI < completionJ // Those who finish first will be ranked later.
		}

		// if both done sort by the doneTime
		return processes[j].DoneTime.Before(processes[i].DoneTime)
	})

	return processes
}

func (m *Model) newReqCnt() int {
	m.reqCnt++
	return m.reqCnt
}

// TODO: Maybe make sure that there isn't any existing process with this UUID
func (m *Model) newUUIDForProcess() string {
	return shortuuid.New()
}
