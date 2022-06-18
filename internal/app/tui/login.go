package tui

import (
	"github.com/rivo/tview"
)

func (ui *UI) makeLoginScreen(p tview.Primitive) tview.Primitive {

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(ui.makeLogo(), 6, 1, false).
		AddItem(p, 0, 1, true)

	return flex
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

	return "ezlogin", ui.makeLoginScreen(flex)
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
		AddItem(ui.modal(form, 72, 5), 0, 50, true).
		AddItem(ui.modal(instr, 72, 5), 0, 50, false)

	return "", ui.makeLoginScreen(flex)
}

func (ui *UI) newLoginInputField(text string) *tview.InputField {
	input := tview.NewInputField().
		SetLabel(text + " ").
		SetLabelColor(ui.theme.Label).SetFieldBackgroundColor(ui.theme.Field)
	input.SetBackgroundColor(ui.theme.Background)
	return input
}
