
package memory

// Mutex : Mutex lock
type Mutex struct {
	locked bool
}

func (m *Mutex) acquire() bool {
	if m.locked {
		return false
	} 

	m.locked = true
	return true
}

func (m *Mutex) release() {
	m.locked = false
}