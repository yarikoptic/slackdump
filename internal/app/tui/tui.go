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
	"   [$ptc:$cbc]         [-:-]",
	"   [$ptc:$cbc]    ▲    [-:-]",
	"   [$ptc:$cbc]   ▲ ▲   [-:-]",
	"   [$ptc:$cbc]         [-:-]",
	"   [$itc]Slackdump",
}

var defTheme = themeLotus1

type screen func() (title string, content tview.Primitive)

type UI struct {
	app      *tview.Application
	pages    *tview.Pages
	msgQueue chan msg
	debug    bool

	theme         tview.Theme
	colorReplacer *strings.Replacer
}

type Option func(*UI)

func WithTheme(theme tview.Theme) Option {
	return func(ui *UI) {
		ui.theme = theme
	}
}

func NewUI(opt ...Option) *UI {
	pages := tview.NewPages()

	ui := &UI{
		app:      tview.NewApplication(),
		pages:    pages,
		msgQueue: make(chan msg, 1),
		theme:    defTheme,
		debug:    false,
	}
	for _, fn := range opt {
		fn(ui)
	}
	ui.colorReplacer = initReplacer(ui.theme)
	tview.Styles = ui.theme

	return ui
}

func (ui *UI) Run(cfg app.Config, creds app.SlackCreds) error {
	screens := []screen{
		ui.makeScrLogin(creds),
		ui.scrMode,
	}

	status := tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(false).
		SetTextAlign(tview.AlignCenter).
		SetText(ui.colorize("[$ptc]F1[$itc] displays a Help screen, [$ptc]F3[$itc] exits."))

	// Create the main layout.
	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(ui.pages, 0, 14, true).
		AddItem(status, 1, 1, false)
	layout.SetBorder(true)

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

	return ui.d(grid)
}

func (ui *UI) lines(w io.Writer, lines []string) {
	for _, line := range lines {
		fmt.Fprintln(w, ui.colorize(line))
	}
}

func (ui *UI) linesEnum(w io.Writer, items []string) {
	for i, line := range items {
		fmt.Fprintln(w, ui.colorize(fmt.Sprintf("[$stc]%2d.[-]   %s", i+1, line)))
	}
}

func (ui *UI) makeInstructions(lines []string) *tview.TextView {
	p := tview.NewTextView().
		SetDynamicColors(true).
		SetWordWrap(true).
		SetRegions(true)
	p.SetTextColor(tview.Styles.PrimitiveBackgroundColor).
		SetBackgroundColor(tview.Styles.ContrastBackgroundColor).
		SetBorder(true).
		SetBorderColor(ui.theme.PrimaryTextColor)

	ui.linesEnum(p, lines)

	return p
}

// makeLogo creates a logo
func (ui *UI) makeLogo() *tview.TextView {
	p := tview.NewTextView().SetDynamicColors(true)
	ui.lines(p, logo)
	return p
}
