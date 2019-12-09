package cpu

import (
	"time"
)

// CPU : virtual CPU
type CPU struct {

	// TotalCycles : total number of cpu cycles run
	TotalCycles int

	// Speed : the minimum time between CPU cycles
	Speed time.Duration
}

// InitCPU : create new CPU
func InitCPU(speed time.Duration) *CPU {
	return &CPU{
		TotalCycles: 0,
		Speed:       speed,
	}
}

// RunCycle : execute a cpu cycle
func (cpu *CPU) RunCycle(runtime int) {

	// Increment number of CPU cycles
	cpu.TotalCycles++

	// Decrease the runtime needed for that process
	runtime--

	// Sleep before the next process
	time.Sleep(cpu.Speed * time.Nanosecond)
}
