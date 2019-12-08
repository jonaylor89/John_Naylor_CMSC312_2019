package tui

import (
	// "log"
	// "fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"

	"github.com/jonaylor89/John_Naylor_CMSC312_2019/kernel"
)

const (
	PROMPT = "[os_simulator]$ "
)

var (
	p        *widgets.Paragraph
	readys   *ProcWidget
	waitings *ProcWidget
	mems     *MemWidget
	i        *TextBox
	grid     *ui.Grid

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
	mems = NewMemWidget(k.Mem)
	mems.SetRect(0, 0, 25, 5)

	p = widgets.NewParagraph()
	p.Text = " CMSC312 Operating System Simulator "
	p.SetRect(0, 0, 25, 5)

	readys = NewProcWidget(&k.ReadyQ)
	readys.Title = " Ready Processes "
	readys.TextStyle = ui.NewStyle(ui.ColorYellow)
	// readys.WrapText = false
	readys.SetRect(0, 0, 25, 8)

	waitings = NewProcWidget(&k.WaitingQ)
	waitings.Title = " Waiting Processes "
	waitings.TextStyle = ui.NewStyle(ui.ColorYellow)
	// waitings.WrapText = false
	waitings.SetRect(0, 0, 25, 8)

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
		ui.NewRow(1.0/3,
			ui.NewCol(1.0/1, p),
		),
		ui.NewRow(1.0/3,
			ui.NewCol(1.0/2, readys),
			ui.NewCol(1.0/2, waitings),
		),
		ui.NewRow(1.0/3,
			ui.NewCol(1.0/1, mems),
		),
	)

	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight-1)
	i.SetRect(0, termHeight-1, termWidth, termHeight)

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

	case "exit", "quit", "q", ":wq":
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
			ui.Render(i)
		case e := <-uiEvents:
			switch e.ID {
			case "<C-c>":
				return
			case "<Resize>":
				payload := e.Payload.(ui.Resize)
				grid.SetRect(0, 0, payload.Width, payload.Height-1)
				i.SetRect(0, payload.Height-1, payload.Width, payload.Height)
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
