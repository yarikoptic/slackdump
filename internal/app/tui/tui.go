package tui

// Debug with: delve --headless debug module as per https://github.com/rivo/tview/issues/351

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rusq/slackdump/v2/internal/app"
)

var logo = []string{
	"",
	"   [$lfg:$lbg]         [:$tbg]",
	"   [$lfg:$lbg]    ▲    [:$tbg]",
	"   [$lfg:$lbg]   ▲ ▲   [:$tbg]",
	"   [$lfg:$lbg]         [:$tbg]",
	"   [$tfg]Slackdump",
}

type theme struct {
	Background     tcell.Color
	InfoBackground tcell.Color
	Border         tcell.Color
	Label          tcell.Color

	LogoFg   string
	LogoBg   string
	TextFg   string
	TextBg   string
	HiTextFg string
}

// https://www.ditig.com/256-colors-cheat-sheet
var (
	themeLotus = theme{
		Background:     tcell.ColorTeal,
		InfoBackground: tcell.ColorBlue,
		Border:         tcell.ColorBlack,

		LogoFg:   "white",
		LogoBg:   "black",
		TextFg:   "black",
		TextBg:   "teal",
		HiTextFg: "white",
	}
	theme5151 = theme{
		Background: tcell.ColorGreen,
		Border:     tcell.ColorBlack,
		Label:      tcell.ColorBlack,

		LogoFg:   "#ffff00",
		LogoBg:   "black",
		TextFg:   "black",
		TextBg:   "green",
		HiTextFg: "#ffff00",
	}
)

var defTheme = themeLotus

type Screen func(nextScreen func()) (title string, content tview.Primitive)

type UI struct {
	app   *tview.Application
	theme theme
	debug bool
}

func NewUI() *UI {
	return &UI{theme: defTheme, app: tview.NewApplication(), debug: false}
}

func (ui *UI) Run(cfg app.Config, creds app.SlackCreds) error {
	screens := []Screen{
		ui.login,
	}

	pages := tview.NewPages()
	pages.SetBackgroundColor(ui.theme.Background)

	status := tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(false).
		SetTextAlign(tview.AlignCenter).
		SetText(ui.colorize("[$hfg]F1[$tfg] displays a Help screen."))
	status.SetBackgroundColor(ui.theme.Background)

	// Create the main layout.
	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(pages, 0, 14, true).
		AddItem(status, 1, 1, false)
	layout.SetBorder(true).SetBackgroundColor(ui.theme.Background).SetBorderColor(ui.theme.Border)

	for index, screen := range screens {
		title, primitive := screen(nil)
		pages.AddPage(title, primitive, true, index == 0)
		layout.SetTitle(title)
	}

	// Start the application.
	if err := ui.app.SetRoot(layout, true).EnableMouse(true).Run(); err != nil {
		return err
	}
	return nil
}

func (ui *UI) d(p tview.Primitive) tview.Primitive {
	type borderer interface {
		SetBorder(show bool) *tview.Box
	}
	if ui.debug {
		if b, ok := p.(borderer); ok {
			b.SetBorder(true)
		}
	}
	return p
}

func (ui *UI) modal(p tview.Primitive, width int, height int) tview.Primitive {
	grid := tview.NewGrid().
		SetColumns(0, width, 0).
		SetRows(0, height, 0).
		AddItem(p, 1, 1, 1, 1, 0, 0, true)

	grid.SetBackgroundColor(ui.theme.Background)
	return ui.d(grid)
}

func (ui *UI) colorize(text string) string {
	r := strings.NewReplacer(
		"$lfg", ui.theme.LogoFg,
		"$lbg", ui.theme.LogoBg,
		"$tfg", ui.theme.TextFg,
		"$tbg", ui.theme.TextBg,
		"$hfg", ui.theme.HiTextFg,
	)
	return r.Replace(text)
}
