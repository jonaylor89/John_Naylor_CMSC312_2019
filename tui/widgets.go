
package tui

import (
	// "log"
	"math"
	"fmt"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"

	"github.com/jonaylor89/John_Naylor_CMSC312_2019/kernel"
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

	p0 := widgets.NewPlot()
	p0.Title = "braille-mode Line Chart"
	p0.Data = sinData
	p0.SetRect(0, 0, 50, 15)
	p0.AxesColor = ui.ColorWhite
	p0.LineColors[0] = ui.ColorGreen

	p1 := widgets.NewPlot()
	p1.Title = "dot-mode line Chart"
	p1.Marker = widgets.MarkerDot
	p1.Data = [][]float64{[]float64{1, 2, 3, 4, 5}}
	p1.SetRect(50, 0, 75, 10)
	p1.DotMarkerRune = '+'
	p1.AxesColor = ui.ColorWhite
	p1.LineColors[0] = ui.ColorYellow
	p1.DrawDirection = widgets.DrawLeft

	p2 := widgets.NewPlot()
	p2.Title = "dot-mode Scatter Plot"
	p2.Marker = widgets.MarkerDot
	p2.Data = make([][]float64, 2)
	p2.Data[0] = []float64{1, 2, 3, 4, 5}
	p2.Data[1] = sinData[1][4:]
	p2.SetRect(0, 15, 50, 30)
	p2.AxesColor = ui.ColorWhite
	p2.LineColors[0] = ui.ColorCyan
	p2.PlotType = widgets.ScatterPlot

	p3 := widgets.NewPlot()
	p3.Title = "braille-mode Scatter Plot"
	p3.Data = make([][]float64, 2)
	p3.Data[0] = []float64{1, 2, 3, 4, 5}
	p3.Data[1] = sinData[1][4:]
	p3.SetRect(45, 15, 80, 30)
	p3.AxesColor = ui.ColorWhite
	p3.LineColors[0] = ui.ColorCyan
	p3.Marker = widgets.MarkerBraille
	p3.PlotType = widgets.ScatterPlot

	p := widgets.NewParagraph()
	p.Text = "CMSC312 Operating System Simulator"
	p.SetRect(0, 0, 25, 5)

	l := widgets.NewList()
	l.Title = "Ready Processes"
	l.Rows = Map(k.ReadyQ, func(p *kernel.Process) string {
		return fmt.Sprintf("%#v", p)
	})
	l.TextStyle = ui.NewStyle(ui.ColorYellow)
	l.WrapText = false
	l.SetRect(0, 0, 25, 8)

	l0 := widgets.NewList()
	l0.Title = "Waiting Processes"
	l0.Rows = Map(k.WaitingQ, func(p *kernel.Process) string {
		return fmt.Sprintf("%#v", p)
	})
	l0.TextStyle = ui.NewStyle(ui.ColorYellow)
	l0.WrapText = false
	l0.SetRect(0, 0, 25, 8)

	grid := ui.NewGrid()

	grid.Set(
		ui.NewRow(1.0/3,
			ui.NewCol(1.0/1, p),
		),
		ui.NewRow(1.0/3,
			ui.NewCol(1.0/2, l),
			ui.NewCol(1.0/2, l0),
		),
		ui.NewRow(1.0/3,
			ui.NewCol(1.0/4, p0),
			ui.NewCol(1.0/4, p1),
			ui.NewCol(1.0/4, p2),
			ui.NewCol(1.0/4, p3),
		),
	)

	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)

	ui.Render(grid)
}