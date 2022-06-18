package tui

// Debug with: delve --headless debug module as per https://github.com/rivo/tview/issues/351

import (
	"fmt"
	"io"
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
	Field          tcell.Color

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
		Label:          tcell.ColorBlack,
		Field:          tcell.ColorBlack,

		LogoFg:   "white",
		LogoBg:   "black",
		TextFg:   "black",
		TextBg:   "teal",
		HiTextFg: "white",
	}
	theme5151 = theme{
		Background:     tcell.ColorGreen,
		InfoBackground: tcell.ColorBlack,
		Border:         tcell.ColorBlack,
		Label:          tcell.ColorBlack,
		Field:          tcell.ColorBlack,

		LogoFg:   "#ffff00",
		LogoBg:   "black",
		TextFg:   "black",
		TextBg:   "green",
		HiTextFg: "#ffff00",
	}
)

var defTheme = theme5151

type screen func() (title string, content tview.Primitive)

type UI struct {
	app      *tview.Application
	pages    *tview.Pages
	msgQueue chan msg
	debug    bool

	theme theme
}

func NewUI() *UI {
	pages := tview.NewPages()

	return &UI{
		app:      tview.NewApplication(),
		pages:    pages,
		theme:    defTheme,
		msgQueue: make(chan msg, 1),
		debug:    false,
	}
}

func (ui *UI) Run(cfg app.Config, creds app.SlackCreds) error {
	screens := []screen{
		ui.scrCookie,
		ui.scrEzLogin,
	}

	ui.pages.SetBackgroundColor(ui.theme.Background)

	status := tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(false).
		SetTextAlign(tview.AlignCenter).
		SetText(ui.colorize("[$hfg]F1[$tfg] displays a Help screen, [$hfg]F3[$tfg] exits."))
	status.SetBackgroundColor(ui.theme.Background)

	// Create the main layout.
	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(ui.pages, 0, 14, true).
		AddItem(status, 1, 1, false)
	layout.SetBorder(true).
		SetBackgroundColor(ui.theme.Background).
		SetBorderColor(ui.theme.Border)

	for index, screen := range screens {
		title, primitive := screen()
		ui.pages.AddPage(title, primitive, true, index == 0)
	}

	ui.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyF1:
			// show help
			return nil
		case tcell.KeyF3:
			// show warning
			ui.sendMessage(wm_quit, nil)
			return nil
		}
		return event
	})
	// Start the message loop
	go ui.messageLoop()
	// Start the application.
	if err := ui.app.SetRoot(layout, true).EnableMouse(true).Run(); err != nil {
		return err
	}
	return nil
}

func (ui *UI) destroy() {
	ui.app.Stop()
	close(ui.msgQueue)
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
		"$lfg", ui.theme.LogoFg, // logo foreground
		"$lbg", ui.theme.LogoBg, // logo background
		"$tfg", ui.theme.TextFg, // text foreground
		"$tbg", ui.theme.TextBg, // text background
		"$hfg", ui.theme.HiTextFg, // hi text background
	)
	return r.Replace(text)
}

func (ui *UI) lines(w io.Writer, lines []string) {
	for _, line := range lines {
		fmt.Fprintln(w, ui.colorize(line))
	}
}

func (ui *UI) linesEnum(w io.Writer, items []string) {
	for i, line := range items {
		fmt.Fprintln(w, ui.colorize(fmt.Sprintf("%2d.   %s", i+1, line)))
	}
}

func (ui *UI) makeInstructions(lines []string) tview.Primitive {
	instructions := tview.NewTextView().
		SetDynamicColors(true).
		SetWordWrap(true)
	instructions.SetBackgroundColor(ui.theme.InfoBackground)
	instructions.SetBorder(true)

	ui.linesEnum(instructions, lines)

	return instructions
}

// makeLogo creates a logo
func (ui *UI) makeLogo() tview.Primitive {
	p := tview.NewTextView().SetDynamicColors(true)
	p.SetBackgroundColor(ui.theme.Background)
	ui.lines(p, logo)
	return p
}
