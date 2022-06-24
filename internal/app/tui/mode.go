package tui

import "github.com/rivo/tview"

type scrDumpMode struct {
	global manager
}

const pgDumpMode pageName = "dump_mode"

func newScrDumpMode(m manager) *scrDumpMode {
	return &scrDumpMode{m}
}

func (wnd *scrDumpMode) Screen() (pageName, tview.Primitive) {
	menu := tview.NewList()
	menu.
		AddItem(" Channels/Threads ", " Dump Channels and Threads", 'm', func() {
			wnd.global.sendMessage(pgMain, wm_switch, "")
		}).
		AddItem(" Export ", " Run Full Workspace Export", 'e', func() {
			wnd.global.sendMessage(pgMain, wm_switch, "")
		}).
		AddItem(" Channels ", " Save channel list to a file", 'c', func() {
			wnd.global.sendMessage(pgMain, wm_switch, "")
		}).
		AddItem(" Users ", " Save user list to a file.", 'u', func() {
			wnd.global.sendMessage(pgMain, wm_switch, "")
		}).
		AddItem(" Exit ", " Exit Slackdump and return to OS.", 'x', func() {
			wnd.global.sendMessage(pgMain, wm_quit, nil)
		})
	applyListTheme(menu)

	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow).
		AddItem(makeHeader("Choose Login Mode"), headerHeight, 0, false).
		AddItem(menu, 0, 1, true)

	return pgDumpMode, flex
}

func (dm *scrDumpMode) WndProc(msg) any { return nil }
