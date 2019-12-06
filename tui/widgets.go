
package tui

import (
	// "log"
	"math"
	"fmt"
	"time"
	"os"
	"os/signal"
	"syscall"
	"strings"
	"strconv"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"

	"github.com/jonaylor89/John_Naylor_CMSC312_2019/kernel"
)

const (
	PROMPT = "[os_simulator]$ "
)

var (
	p0 *widgets.Plot
	p  *widgets.Paragraph
	l  *widgets.List
	l0  *widgets.List
	i *TextBox
	grid *ui.Grid
	
	updateInterval = time.Second / 10

	PRINTABLE_KEYS = append(
		strings.Split(
			"0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ,./<>?;:'\"[]\\{}|`~!@#$%^&*()-_=+",
			"",
		),
		"<Space>",
		"<Tab>",
		"<Enter>",
	)
)

func InitWidgets(k *kernel.Kernel) {
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

	i = NewTextBox()
	i.SetText(PROMPT)
	i.SetRect(25, 25, 50, 40)
	i.Border = false
	i.ShowCursor = true

}

func Map(vs []*kernel.Process, f func(*kernel.Process) string) []string {
    vsm := make([]string, len(vs))
    for i, v := range vs {
        vsm[i] = f(v)
    }
    return vsm
}

func Render() {
	// TUI
	grid = ui.NewGrid()

	grid.Set(
		ui.NewRow(1.0/10,
			ui.NewCol(1.0/1, p),
		),
		ui.NewRow(1.0/3,
			ui.NewCol(1.0/2, l),
			ui.NewCol(1.0/2, l0),
		),
		ui.NewRow(1.0/3,
			ui.NewCol(1.0/1, p0),
		),
		ui.NewRow(1.0/10,
			ui.NewCol(1.0/1, i),
		),
	)

	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)

	ui.Render(grid)
}

func Launch(args []string, ch chan *kernel.Process) bool {

	switch args[0] {

	case "load":
		if len(args) != 3 {
			// fmt.Println("`load` requires a filename and number of processes as an argument")
			break
		}

		filename := args[1]
		numOfProc, err := strconv.Atoi(args[2])
		if err != nil {
			// fmt.Println("Could not get number of processes")
			break
		}

		if numOfProc <= 0 {
			// fmt.Println("`load` number of processes must be postive")
			break
		}

		err = kernel.LoadTemplate(filename, numOfProc, ch)
		if err != nil {
			break
		}

	case "exit", "quit", "q", ":qw":
		return true

	default:
		break
	}

	return false
}

func EventLoop(ch chan *kernel.Process) {
	drawTicker := time.NewTicker(updateInterval).C

	// handles kill signal sent to go
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
			case "<C-c>":
				return
			case "<Resize>":
				payload := e.Payload.(ui.Resize)
				grid.SetRect(0, 0, payload.Width, payload.Height)
				ui.Clear()

			case "<Left>":
				i.MoveCursorLeft()
			case "<Right>":
				i.MoveCursorRight()
			case "<Up>":
				i.MoveCursorUp()
			case "<Down>":
				i.MoveCursorDown()
			case "<Backspace>":
				i.Backspace()
			case "<Enter>":
				// Execute command that's set

				args := strings.Split(i.GetText(), " ")[1:]

				if quit := Launch(args, ch); quit {
					return
				}

				i.ClearText()
				i.SetText(PROMPT)
			case "<Tab>":
				i.InsertText("\t")
			case "<Space>":
				i.InsertText(" ")
			default:
				if ContainsString(PRINTABLE_KEYS, e.ID) {
					i.InsertText(e.ID)
				}
			}
		}
	}
}