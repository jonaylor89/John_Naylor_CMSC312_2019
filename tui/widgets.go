
package tui

import (
	// "log"
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

	p := widgets.NewParagraph()
	p.Text = "Hello World!"
	p.SetRect(0, 0, 25, 5)

	l := widgets.NewList()
	l.Title = "List"
	l.Rows = Map(k.ReadyQ, func(p *kernel.Process) string {
		return fmt.Sprintf("%#v", p)
	})
	l.TextStyle = ui.NewStyle(ui.ColorYellow)
	l.WrapText = false
	l.SetRect(0, 0, 25, 8)

	grid := ui.NewGrid()

	grid.Set(
		ui.NewRow(1.0/3,
			ui.NewCol(1.0/3, l),
		),
		ui.NewRow(1.0/3,
			ui.NewCol(1.0/2, p),
		),
	)

	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)

	ui.Render(grid)
}