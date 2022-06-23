package tui

import (
	"github.com/rivo/tview"
)

type scrLoginMode struct {
	global messenger
}

func newLoginMode(m messenger) *scrLoginMode {
	return &scrLoginMode{global: m}
}

func (l *scrLoginMode) Screen() (string, tview.Primitive) {
	menu := tview.NewList()
	menu.
		AddItem(" Login With Browser ", " Automatic login (EZ-Login 3000)", '1', func() {
			l.global.sendMessage(wm_page, "login")
		}).
		AddItem(" Login With Token and Cookie ", " Login with token and cookie or file", '2', func() {
			l.global.sendMessage(wm_page, screenDumpMode)
		}).
		AddItem(" Exit ", " Exit Slackdump and return to OS.", 'x', func() {
			l.global.sendMessage(wm_quit, nil)
		})
	applyListTheme(menu)

	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow).
		AddItem(makeHeader("Choose Login Mode"), headerHeight, 0, false).
		AddItem(menu, 0, 1, true)
	return "login_menu", flex
}
