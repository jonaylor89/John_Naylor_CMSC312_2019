package memory

// Mutex : Mutex lock
type Mutex struct {
	locked bool
}

// Acquire : lock the mutex lock
func (m *Mutex) Acquire() bool {
	if m.locked {
		return false
	}

	m.locked = true
	return true
}

// Release : unlock the mutex lock
func (m *Mutex) Release() {
	m.locked = false
}
