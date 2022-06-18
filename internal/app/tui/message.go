package tui

// you may find that this implementation reminds you of win32 messaging system,
// and you won't be mistaken.

// message is the type of message
type message uint

const (
	wm_quit message = iota // quit the application
)

// msg is the message iself.
type msg struct {
	message message
	param   any
}

func (ui *UI) sendMessage(m message, param any) {
	ui.msgQueue <- msg{message: m, param: param}
}

func (ui *UI) messageLoop() {
	for msg := range ui.msgQueue {
		switch msg.message {
		case wm_quit:
			ui.destroy()
			return
		}
	}
}
