
package memory

var (

	// frames : basically a cache of pages for the simulator because of the lack of hardware
	frames []Page

	// pages : pages in secondary memory
	virtualMemory []Page
)

// Page : a page of memory
type Page struct {

}

type Mutex struct {
	locked bool
}

func (m *Mutex) acquire() bool {
	if (m.locked) {
		return false
	} else {
		m.locked = true
		return true
	}
}

func (m *Mutex) release() {
	m.locked = false	
}
