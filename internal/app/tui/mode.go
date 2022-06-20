package tui

import "github.com/rivo/tview"

func (ui *UI) scrMode() (string, tview.Primitive) {
	title := tview.NewTextView()
	main := tview.NewPages()
	status := tview.NewTextView()
	// info := tview.NewTextView()

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(title, 1, 0, false).
		AddItem(main, 0, 1, false).
		AddItem(status, 1, 0, false)
	return "modes", flex
}
