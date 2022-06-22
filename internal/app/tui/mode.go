package tui

import "github.com/rivo/tview"

type scrDumpMode struct {
	global messenger
}

const screenDumpMode = "dump_mode"

func newScrDumpMode(m messenger) *scrDumpMode {
	return &scrDumpMode{m}
}

func (dm *scrDumpMode) Screen() (string, tview.Primitive) {
	menu := tview.NewList()
	menu.
		AddItem(" Channels/Threads ", " Dump Channels and Threads", 'm', func() {
			dm.global.sendMessage(wm_page, "")
		}).
		AddItem(" Export ", " Run Full Workspace Export", 'e', func() {
			dm.global.sendMessage(wm_page, "")
		}).
		AddItem(" Channels ", " Save channel list to a file", 'c', func() {
			dm.global.sendMessage(wm_page, "")
		}).
		AddItem(" Users ", " Save user list to a file.", 'u', func() {
			dm.global.sendMessage(wm_page, "")
		}).
		AddItem(" Exit ", " Exit Slackdump and return to OS.", 'x', func() {
			dm.global.sendMessage(wm_quit, nil)
		})
	applyListTheme(menu)

	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow).
		AddItem(makeHeader("Choose Login Mode"), headerHeight, 0, false).
		AddItem(menu, 0, 1, true)
	return screenDumpMode, flex
}
