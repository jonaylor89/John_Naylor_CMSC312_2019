package memory

const (
	// FrameLength : length of frame of memory
	FrameLength = 256
)

var (
	// pages : pages in secondary memory
	virtualMemory [1024]Page
)

// // RAM : virtual physical memory
// type RAM struct {
// 	// frames : basically a cache of pages for the simulator because of the lack of hardware
// 	frames []*Page
// }

// Page : a page of memory
type Page struct {
	PID    int // Process ID of the process using this page
	length int
	contents [256]byte
}

