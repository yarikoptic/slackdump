package tui

import (
	"github.com/rivo/tview"
)

const pgLoginMenu = pageName("login_menu")

type scrLoginMode struct {
	global manager
}

func newLoginMode(m manager) *scrLoginMode {
	return &scrLoginMode{global: m}
}

func (l *scrLoginMode) Screen() (pageName, tview.Primitive) {
	menu := tview.NewList()
	menu.
		AddItem(" Login With Browser ", " Automatic login (EZ-Login 3000)", '1', func() {
			l.global.sendMessage(pgMain, wm_switch, pgLogin)
		}).
		AddItem(" Login With Token and Cookie ", " Login with token and cookie or file", '2', func() {
			l.global.sendMessage(pgMain, wm_switch, pgDumpMode)
		}).
		AddItem(" Exit ", " Exit Slackdump and return to OS.", 'x', func() {
			l.global.sendMessage(pgMain, wm_quit, nil)
		})
	applyListTheme(menu)

	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow).
		AddItem(makeHeader("CHOOSE LOGIN MODE"), headerHeight, 0, false).
		AddItem(menu, 0, 1, true)
	return pgLoginMenu, flex
}

func (l *scrLoginMode) WndProc(m msg) any {
	return false
}
