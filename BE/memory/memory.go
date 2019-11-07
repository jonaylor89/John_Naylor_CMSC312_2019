package memory

import (
	"math"
)

var (
	pageNum = 0
)

type Memory struct {
	// PageSize : length of page contents in Mb as a power of 2
	PageSize int

	// TotalRam : Total amount of physical memory in the simulator in Mb as a power of 2
	TotalRam int

	// PageTable : map page id to index of frame
	PageTable map[int]int

	// VirtualMemory : pages in secondary memory
	VirtualMemory []*Page

	// PhysicalMemory : Memory frames in RAM
	PhysicalMemory []*Page
}

// Page : a page of memory
type Page struct {
	PageID	 int // ID of page
	ProcID   int // Process ID of the process using this page
	contents []byte
}

// GetPage : get a page of memory
func (m *Memory) Get(pageNum int) *Page {
	// Check for page in PhyiscalMemory

	if val, ok := m.PageTable[pageNum]; ok {
		return m.PhysicalMemory[val]
	}

	for _, page := range m.VirtualMemory {
		if page.PageID == pageNum {

			m.moveToPhysicalMemory(page)

			return page
		}
	}

	// Page doesn't exist
	return nil

	// 		put page from vm in main memory in either a free space or replacement algorithm
	// 		add to page table
	// 		return page
}

// AddPage : Add pages of memory to memory pool, return PageIDs
func (m *Memory) Add(requirement int, pid int) []int {
	// Append new page to virtual memory

	pageIds := []int{}

	numOfPages := int(math.Ceil(float64(requirement) / float64(m.PageSize)))

	for i := 0; i < numOfPages; i++ {

		pageNum++ 

		p := &Page{
			PageID: pageNum,
			ProcID: pid,
			contents: make([]byte, 0, 256),
		}

		pageIds = append(pageIds, pageNum)

		m.VirtualMemory = append(m.VirtualMemory, p)	
	}

	return pageIds
}

func (m *Memory) moveToPhysicalMemory(p *Page) {
	// if there is an empty space, put page in empty space	
	// if there isn't an empty space, run a replace procedure

	// Always add new entry to page table and remove old entry if replaced
}

// RemovePages : remove all pages associated with a pid
func (m *Memory) RemovePages(pid int) {

	// Remove pages from physical memory
	for i, page := range m.PhysicalMemory {
		if page.ProcID == pid {
			m.PhysicalMemory = remove(m.PhysicalMemory, i)
		}
	}	

	// Remove pages from virtual memory
	for i, page := range m.VirtualMemory {
		if page.ProcID == pid {
			m.VirtualMemory = remove(m.VirtualMemory, i)
		}
	}
}

func remove(slice []*Page, s int) []*Page {
	slice[s] = slice[len(slice)-1] // Copy last element to index i.
	// slice[len(slice)-1] = nil   	// Erase last element (write zero value)
	slice = slice[:len(slice)-1] // Truncate slice.

	return slice
}