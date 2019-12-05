
package tui

import (
	ui "github.com/gizak/termui/v3"
)

func EventLoop() {
	for e := range ui.PollEvents() {
		if e.Type == ui.KeyboardEvent {
			break
		}
	}
}