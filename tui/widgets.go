package tui

import (
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"

	"github.com/jonaylor89/John_Naylor_CMSC312_2019/sched"
)

const (
	PROMPT = "[os_simulator]$ "
)

var (
	header   *widgets.Paragraph
	readys   *ProcWidget
	waitings *ProcWidget
	mems     *MemWidget
	shell    *TextBox
	grid     *ui.Grid

	updateInterval = time.Second / 10

	PrintableKeys = append(
		strings.Split(
			"0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ,./<>?;:'\"[]\\{}|`~!@#$%^&*()-_=+",
			"",
		),
		"<Space>",
		"<Tab>",
		"<Enter>",
	)
)

func InitWidgets(s *sched.Scheduler) {
	mems = NewMemWidget(s.Mem)
	mems.SetRect(0, 0, 25, 5)

	header = widgets.NewParagraph()
	header.Text = " CMSC 312 Operating System Simulator "
	header.SetRect(0, 0, 25, 5)

	readys = NewProcWidget(&s.ReadyQ)
	readys.Title = " Ready Processes "
	readys.TextStyle = ui.NewStyle(ui.ColorYellow)
	// readys.WrapText = false
	readys.SetRect(0, 0, 25, 8)

	waitings = NewProcWidget(&s.WaitingQ)
	waitings.Title = " Waiting Processes "
	waitings.TextStyle = ui.NewStyle(ui.ColorYellow)
	// waitings.WrapText = false
	waitings.SetRect(0, 0, 25, 8)

	shell = NewTextBox()
	shell.SetText(PROMPT)
	shell.SetRect(25, 25, 50, 40)
	shell.Border = false
	shell.ShowCursor = true

}

func Map(vs []*sched.Process, f func(*sched.Process) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

func RenderTUI() {
	// TUI

	grid = ui.NewGrid()

	// et grid dimensions
	grid.Set(
		ui.NewRow(1.0/3,
			ui.NewCol(1.0/1, header),
		),
		ui.NewRow(1.0/3,
			ui.NewCol(1.0/2, readys),
			ui.NewCol(1.0/2, waitings),
		),
		ui.NewRow(1.0/3,
			ui.NewCol(1.0/1, mems),
		),
	)

	// Adapt to terminal width and height
	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight-1)
	shell.SetRect(0, termHeight-1, termWidth, termHeight)

	ui.Render(grid)
}

func Launch(args []string, ch chan *sched.Process) bool {

	// Interpret commands
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

		err = sched.LoadTemplate(filename, numOfProc, ch)
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

// EventLoop : Main tui event loop
func EventLoop(ch chan *sched.Process) {

	// framerate
	drawTicker := time.NewTicker(updateInterval).C

	// handles kill signal sent to go
	sigTerm := make(chan os.Signal, 2)
	signal.Notify(sigTerm, os.Interrupt, syscall.SIGTERM)

	uiEvents := ui.PollEvents()

	for {
		select {
		case <-sigTerm:
			return
		case <-drawTicker:
			ui.Render(grid)
			ui.Render(shell)
		case e := <-uiEvents:
			switch e.ID {
			case "<C-c>":
				return
			case "<Resize>":
				payload := e.Payload.(ui.Resize)
				grid.SetRect(0, 0, payload.Width, payload.Height-1)
				shell.SetRect(0, payload.Height-1, payload.Width, payload.Height)
				ui.Clear()

			case "<Left>":
				shell.MoveCursorLeft()
			case "<Right>":
				shell.MoveCursorRight()
			case "<Up>":
				shell.MoveCursorUp()
			case "<Down>":
				shell.MoveCursorDown()
			case "<Backspace>":
				shell.Backspace()
			case "<Enter>":
				// Execute command that's set

				args := strings.Split(shell.GetText(), " ")[1:]

				if quit := Launch(args, ch); quit {
					return
				}

				shell.ClearText()
				shell.SetText(PROMPT)
			case "<Tab>":
				shell.InsertText("\t")
			case "<Space>":
				shell.InsertText(" ")
			default:
				if ContainsString(PrintableKeys, e.ID) {
					shell.InsertText(e.ID)
				}
			}
		}
	}
}
