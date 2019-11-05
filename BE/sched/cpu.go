
package sched

import (
	"time"
)

// CPU : virtual CPU
type CPU struct {

	// TotalCycles : total number of cpu cycles runh
	TotalCycles int

	// Speed : the minimum time between CPU cycles
	Speed time.Duration
}

// RunCycle : execute a cpu cycle
func (cpu *CPU) RunCycle(p *Process) {

	cpu.TotalCycles++

	p.runtime--
	time.Sleep(cpu.Speed * time.Millisecond)
}