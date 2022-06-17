package tui

import (
	"fmt"
	"io"

	"github.com/rivo/tview"
)

func (ui *UI) login(nextPage func()) (title string, content tview.Primitive) {
	header := tview.NewTextView().SetDynamicColors(true)
	header.SetBackgroundColor(ui.theme.Background)
	ui.writeLines(header, logo)

	input := tview.NewInputField().
		SetLabel("Slack Workspace ").
		SetLabelColor(ui.theme.Label)
	input.SetBackgroundColor(ui.theme.Background)

	instructions := tview.NewTextView().SetDynamicColors(true).SetWordWrap(true)
	instructions.SetBackgroundColor(ui.theme.InfoBackground)
	instructions.SetBorder(true)
	ui.loginInstructions(instructions)

	instrFlex := ui.modal(instructions, 60, 6)

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(header, 8, 1, false).
		AddItem(ui.modal(input, 60, 3), 3, 1, true).
		AddItem(instrFlex, 6, 0, false)

	return " EZ-Login 3000 ", flex
}

func (ui *UI) writeLines(w io.Writer, lines []string) {
	for _, line := range lines {
		fmt.Fprintln(w, ui.colorize(line))
	}
}

func (ui *UI) loginInstructions(w io.Writer) {
	lines := []string{
		"Enter the Slack workspace name OR\n      paste the URL of your Slack workspace.",
		"Press ENTER, the browser will open.",
		"Login as usual (browser will close automatically).",
	}
	for i, line := range lines {
		fmt.Fprintln(w, ui.colorize(fmt.Sprintf("%2d.   %s", i+1, line)))
	}
}
