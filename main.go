package main

import (
	// "bufio"
	// "fmt"
	"log"
	// "os"
	// "strconv"
	// "strings"
	// "time"

	ui "github.com/gizak/termui/v3"

	"github.com/jonaylor89/John_Naylor_CMSC312_2019/config"
	"github.com/jonaylor89/John_Naylor_CMSC312_2019/cpu"
	"github.com/jonaylor89/John_Naylor_CMSC312_2019/sched"
	"github.com/jonaylor89/John_Naylor_CMSC312_2019/memory"
	"github.com/jonaylor89/John_Naylor_CMSC312_2019/tui"
)

func main() {

	conf := config.ReadConfig("config.yml")

	// Message channel to kernel
	ch := make(chan *sched.Process, conf.ProcChanSize)
	defer close(ch)

	cpu1 := &cpu.CPU{
		TotalCycles: 0,
		Speed:       conf.CPU.ClockSpeed1,
	}

	// cpu2 := &cpu.CPU{
	// 	TotalCycles: 0,
	// 	Speed: conf.CPU.ClockSpeed2,
	// }

	mem := &memory.Memory{
		PageSize:       conf.Memory.PageSize,
		TotalRam:       conf.Memory.TotalRam,
		PageTable:      make(map[int]int),
		VirtualMemory:  make([]*memory.Page, 0),
		PhysicalMemory: make([]*memory.Page, 0, conf.Memory.TotalRam/conf.Memory.PageSize),
	}

	s := &sched.Scheduler{
		CPU:               cpu1,
		Mem:               mem,
		InMsg:             ch,
		ReadyQ:            []*sched.Process{},
		WaitingQ:          []*sched.Process{},
		MinimumFreeFrames: conf.MinimumFreeFrames,
		TimeQuantum:       conf.Sched.TimeQuantum,
		Mailboxes: []chan byte{
			make(chan byte, 10),
			make(chan byte, 10),
			make(chan byte, 10),
			make(chan byte, 10),
			make(chan byte, 10),
			make(chan byte, 10),
			make(chan byte, 10),
			make(chan byte, 10),
			make(chan byte, 10),
			make(chan byte, 10),
		},
	}

	// Run the scheduler
	go s.RunRoundRobin()
	// go k.RunFirstComeFirst

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	tui.InitWidgets(s)
	tui.Render()
	tui.EventLoop(ch)

}
