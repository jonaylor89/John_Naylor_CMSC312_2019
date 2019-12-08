package tui

import (
	"fmt"
	"strconv"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/jonaylor89/John_Naylor_CMSC312_2019/sched"
)

type ProcWidget struct {
	*widgets.Table
	updateInterval time.Duration
	procs          *[]*sched.Process
}

func NewProcWidget(processes *[]*sched.Process) *ProcWidget {
	self := &ProcWidget{
		Table:          widgets.NewTable(),
		updateInterval: time.Second,
		procs:          processes,
	}

	self.TextAlignment = ui.AlignCenter
	self.RowSeparator = false

	self.update()

	go func() {
		for range time.NewTicker(self.updateInterval).C {
			self.Lock()
			self.update()
			self.Unlock()
		}
	}()

	return self
}

// update :  converts a []*kernel.Process to a [][]string and sets it to the table Rows
func (self *ProcWidget) update() {
	strings := make([][]string, len(*self.procs)+1)
	strings[0] = []string{"PID", "Name", "CPU", "Mem"}
	for i := range *self.procs {
		strings[i+1] = make([]string, 4)
		strings[i+1][0] = strconv.Itoa((*self.procs)[i].PID)
		strings[i+1][1] = (*self.procs)[i].Name
		strings[i+1][2] = fmt.Sprintf("%4s", strconv.Itoa((*self.procs)[i].Runtime))
		strings[i+1][3] = fmt.Sprintf("%4s", strconv.Itoa((*self.procs)[i].Memory))
	}

	self.Rows = strings
}
