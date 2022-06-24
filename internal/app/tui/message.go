package tui

import (
	"fmt"
)

// you may find that this implementation reminds you of win32 messaging system,
// and you won't be mistaken.

// message is the type of message
type message uint

const (
	wm_quit      message = iota // quit the application
	wm_switch                   // switch to page
	wm_show                     // show page
	wm_killfocus                // sent to the window when it loses focus
	wm_close                    // closes the page
	wm_settext                  // set window text
)

// gID is the global message counter, it will increase with each message.
var gID = int64(1000) // start value is 1000.

// msg is the message iself.
type msg struct {
	id      int64 // message ID to track it across functions
	page    pageName
	message message
	param   any
}

func (m msg) String() string {
	return fmt.Sprintf("id: %10d MESSAGE: %d, WND: %q, PARAM: %v", m.id, m.message, m.page, m.param)
}

type manager interface {
	sendMessage(pageName, message, any) any
}

func (ui *UI) sendMessage(wnd pageName, m message, param any) any {
	currID := gID
	gID++ // increase global message ID.

	msg := msg{id: currID, page: wnd, message: m, param: param}
	result := ui.dispatchMessage(msg)
	ui.log.Debugf("result: %v", result)

	return result
}

func (ui *UI) dispatchMessage(msg msg) any {
	wnd, ok := ui.wnd[msg.page]
	if !ok {
		return ui.WndProc(msg)
	}
	return wnd.WndProc(msg)
}
