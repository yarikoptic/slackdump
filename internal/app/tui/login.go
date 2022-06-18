package tui

import (
	"fmt"
	"strconv"

	"github.com/rivo/tview"
)

func (ui *UI) login() (string, tview.Primitive) {
	loginScreens := []screen{
		ui.scrEzLogin,
		ui.scrCookie,
	}

	pages := tview.NewPages()
	pages.SetBackgroundColor(ui.theme.Background)

	info := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(false).
		SetHighlightedFunc(func(added, removed, remaining []string) {
			pages.SwitchToPage(added[0])
		})
	info.SetBackgroundColor(ui.theme.Background)

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(ui.makeLogo(), 7, 0, false).
		AddItem(ui.modal(info, 60, 3), 3, 0, false).
		AddItem(pages, 0, 1, true)

	for i, screen := range loginScreens {
		title, primitive := screen()
		pages.AddPage(strconv.Itoa(i), primitive, true, i == 0)
		fmt.Fprintf(info, ui.colorize(`%d ["%d"][$tfg]%s[white][""]  `), i+1, i, title)
		if i == 0 {
			info.Highlight(strconv.Itoa(i)).ScrollToHighlight()
		}
	}

	return "login", flex
}

func (ui *UI) scrEzLogin() (title string, content tview.Primitive) {
	items := []string{
		"Enter the Slack workspace name OR\n      paste the URL of your Slack workspace.",
		"Press ENTER, the browser will open.",
		"Login as usual (browser will close automatically).",
	}

	input := ui.newLoginInputField("Slack Workspace ")
	instrFlex := ui.modal(ui.makeInstructions(items), 60, 6)
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(ui.modal(input, 60, 3), 3, 1, true).
		AddItem(instrFlex, 6, 0, false)

	return "EZ-Login", flex
}

func (ui *UI) scrCookie() (string, tview.Primitive) {
	var (
		token  = ui.newLoginInputField("Token (xoxc-)")
		cookie = ui.newLoginInputField("Cookie (xoxd-)")
	)

	form := tview.NewForm().
		AddFormItem(token).
		AddFormItem(cookie)

	form.SetBackgroundColor(ui.theme.Background)
	form.SetLabelColor(ui.theme.Label)
	form.SetFieldBackgroundColor(ui.theme.Field)

	instr := ui.makeInstructions([]string{"Follow the steps on\n\n  https://github.com/rusq/slackdump/blob/master/doc/login-manual.rst"})

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(ui.modal(form, 72, 5), 5, 0, true).
		AddItem(ui.modal(instr, 72, 5), 0, 50, false)

	return "Cookie", flex
}

func (ui *UI) newLoginInputField(text string) *tview.InputField {
	input := tview.NewInputField().
		SetLabel(text + " ").
		SetLabelColor(ui.theme.Label).SetFieldBackgroundColor(ui.theme.Field)
	input.SetBackgroundColor(ui.theme.Background)
	return input
}
