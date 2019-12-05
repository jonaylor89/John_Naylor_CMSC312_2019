
package tui

import (
	// "log"
	"math"
	"fmt"
	"time"
	"os"
	"os/signal"
	"syscall"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"

	"github.com/jonaylor89/John_Naylor_CMSC312_2019/kernel"
)

var (
	p0 *widgets.Plot
	p  *widgets.Paragraph
	l  *widgets.List
	l0  *widgets.List
	grid *ui.Grid
	
	updateInterval = time.Second
)

func initWidgets() {
	
}

func Map(vs []*kernel.Process, f func(*kernel.Process) string) []string {
    vsm := make([]string, len(vs))
    for i, v := range vs {
        vsm[i] = f(v)
    }
    return vsm
}

func Render(k *kernel.Kernel) {
	// TUI

	sinData := func() [][]float64 {
		n := 220
		data := make([][]float64, 2)
		data[0] = make([]float64, n)
		data[1] = make([]float64, n)
		for i := 0; i < n; i++ {
			data[0][i] = 1 + math.Sin(float64(i)/5)
			data[1][i] = 1 + math.Cos(float64(i)/5)
		}
		return data
	}()

	p0 = widgets.NewPlot()
	p0.Title = "Memory Usage"
	p0.Data = sinData
	p0.SetRect(0, 0, 50, 15)
	p0.AxesColor = ui.ColorWhite
	p0.LineColors[0] = ui.ColorGreen
	p0.LineColors[1] = ui.ColorBlue

	p = widgets.NewParagraph()
	p.Text = "CMSC312 Operating System Simulator (press `q` to quit)"
	p.SetRect(0, 0, 25, 5)

	l = widgets.NewList()
	l.Title = "Ready Processes"
	l.Rows = Map(k.ReadyQ, func(p *kernel.Process) string {
		return fmt.Sprintf("%#v", p)
	})
	l.TextStyle = ui.NewStyle(ui.ColorYellow)
	l.WrapText = false
	l.SetRect(0, 0, 25, 8)

	l0 = widgets.NewList()
	l0.Title = "Waiting Processes"
	l0.Rows = Map(k.WaitingQ, func(p *kernel.Process) string {
		return fmt.Sprintf("%#v", p)
	})
	l0.TextStyle = ui.NewStyle(ui.ColorYellow)
	l0.WrapText = false
	l0.SetRect(0, 0, 25, 8)

	grid = ui.NewGrid()

	grid.Set(
		ui.NewRow(1.0/3,
			ui.NewCol(1.0/1, p),
		),
		ui.NewRow(1.0/3,
			ui.NewCol(1.0/2, l),
			ui.NewCol(1.0/2, l0),
		),
		ui.NewRow(1.0/3,
			ui.NewCol(1.0/1, p0),
		),
	)

	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)

	ui.Render(grid)
}

func EventLoop() {
	drawTicker := time.NewTicker(updateInterval).C

	// handles kill signal sent to gotop
	sigTerm := make(chan os.Signal, 2)
	signal.Notify(sigTerm, os.Interrupt, syscall.SIGTERM)

	uiEvents := ui.PollEvents()

	// previousKey := ""

	for {
		select {
		case <-sigTerm:
			return
		case <-drawTicker:
			ui.Render(grid)
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			case "<Resize>":
				payload := e.Payload.(ui.Resize)
				grid.SetRect(0, 0, payload.Width, payload.Height)
				ui.Clear()
			}
		}
	}
}