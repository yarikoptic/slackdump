package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const pgHelp pageName = "help"

type scrHelp struct {
	global manager
	wnd    *tview.TextView
}

func newScrHelp(m manager) *scrHelp {
	help := styleInstructions(tview.NewTextView())
	help.SetScrollable(true)

	help.SetTitle(helpTitle).
		SetBorder(true).
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyESC {
				m.sendMessage(pgMain, wm_close, pgHelp)
				return nil
			}
			return event
		})

	return &scrHelp{
		global: m,
		wnd:    help,
	}
}

const helpTitle = ` Context Help `

func (h *scrHelp) Screen() (pageName, tview.Primitive) {
	return pgHelp, modal(shadow(h.wnd, nil), 60, 20)
}

func (wnd *scrHelp) WndProc(msg msg) any {
	if msg.page != pgHelp {
		// not our message
		return nil
	}
	switch msg.message {
	case wm_settext:
		wnd.wnd.SetText(msg.param.(string))
		return true
	}
	return nil
}
