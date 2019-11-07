package memory

type Memory struct {
	// PageSize : length of page contents in Mb as a power of 2
	PageSize int

	// TotalRam : Total amount of physical memory in the simulator in Mb as a power of 2
	TotalRam int

	// PageTable : map page id to index of frame
	PageTable map[int]int

	// VirtualMemory : pages in secondary memory
	VirtualMemory []Page

	// PhysicalMemory : Memory frames in RAM
	Physicalmemory []Page
}

// Page : a page of memory
type Page struct {
	PageID	 int // ID of page
	ProcID   int // Process ID of the process using this page
	contents [PageSize]byte
}

// GetPage : get a page of memory
func (m *Memory) GetPage(pageNum int) *Page {
	// Check for page in PhyiscalMemory

	if val, ok := m.PageTable[pageNum]; ok {
		return m.PhysicalMemory[val]
	}

	for _, page := range m.VirtualMemory {
		if page.PageID == pageNum {

			m.addToPhysicalMemory(page)

			return page
		}
	}

	// Page doesn't exist
	return nil

	// 		put page from vm in main memory in either a free space or replacement algorithm
	// 		add to page table
	// 		return page
}

// AddPage : Add a page of memory to memory pool
func (m *Memory) AddPage(p Page) []int {
	// Append new page to virtual memory
}

func (m *Memory) addToPhysicalMemory(p Page) {
	// if there is an empty space, put page in empty space	
	// if there isn't an empty space, run a replace procedure

	// Always add new entry to page table and remove old entry if replaced
}
