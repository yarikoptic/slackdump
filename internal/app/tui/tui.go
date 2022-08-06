package tui

// Debug with: delve --headless debug module as per https://github.com/rivo/tview/issues/351

import (
	"container/list"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rusq/dlog"
	"github.com/rusq/slackdump/v2/internal/app"
)

var defTheme = themeLotus1 // see colors.go

const pgMain pageName = ""

type UI struct {
	app   *tview.Application
	pages *tview.Pages

	theme         tview.Theme
	colorReplacer *strings.Replacer

	pageHistory *list.List

	log *dlog.Logger
	lf  io.WriteCloser // log file

	wnd map[pageName]window
}

type window interface {
	Screen() (pageName, tview.Primitive)
	WndProc(msg) any
}

type pageName string

type Option func(*UI)

func WithTheme(theme tview.Theme) Option {
	return func(ui *UI) {
		ui.theme = theme
	}
}

func NewUI(opt ...Option) *UI {
	pages := tview.NewPages()

	lf, err := os.OpenFile("tui.log", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		//TODO
		panic("unable to initialise logging")
	}

	ui := &UI{
		app:   tview.NewApplication(),
		pages: pages,
		wnd:   make(map[pageName]window),
		theme: defTheme,

		pageHistory: list.New(),

		log: dlog.New(lf, "", log.LstdFlags|log.Lshortfile, true),
		lf:  lf,
	}
	for _, fn := range opt {
		fn(ui)
	}
	ui.colorReplacer = initReplacer(ui.theme)
	tview.Styles = ui.theme

	return ui
}

func (ui *UI) Run(cfg app.Config, creds app.SlackCreds) error {
	screens := []window{
		newLoginMode(ui),
		newScrDumpMode(ui),
		newScrHelp(ui),
		// ui.makeScrLogin(creds),
	}

	for index, wnd := range screens {
		name, primitive := wnd.Screen()
		ui.pages.AddPage(string(name), primitive, true, index == 0)
		ui.wnd[name] = wnd
	}

	status := tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(false).
		SetTextAlign(tview.AlignCenter).
		SetText(colorize("[$ptc]F1[$itc] displays a Help screen, [$ptc]F3[$itc] exits, [$ptc]F9[$itc] Options."))

	// Create the main layout.
	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(ui.pages, 0, 14, true).
		AddItem(status, 1, 1, false)
	layout.SetBorder(true)

	ui.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyF1:
			curr, _ := ui.pages.GetFrontPage()
			topic := topics[pageName(curr)]
			ui.sendMessage(pgHelp, wm_settext, colorize(topic))
			ui.sendMessage(pgMain, wm_show, pgHelp)
			return nil
		case tcell.KeyF3:
			// show warning
			ui.sendMessage(pgMain, wm_quit, nil)
			return nil
		case tcell.KeyF9:
			// show parameters
			return nil
		}
		return event
	})

	// Start the application.
	if err := ui.app.SetRoot(layout, true).EnableMouse(true).Run(); err != nil {
		return err
	}
	return nil
}

// WndProc for the main window
func (ui *UI) WndProc(msg msg) any {
	// TODO type checking.
	switch msg.message {
	case wm_quit:
		ui.destroy()
		return true
	case wm_switch, wm_show:
		return ui.setFocus(msg)
	case wm_close:
		page, err := asPage(msg.param)
		if err != nil {
			ui.log.Printf("message: %s, error: %s", msg, err)
			return false
		}
		ui.pages.HidePage(page)
		ui.pageHistory.Remove(ui.pageHistory.Back())
		return true
	}
	return false
}

func (ui *UI) setFocus(msg msg) any {
	page, err := asPage(msg.param)
	if err != nil {
		ui.log.Printf("message: %s, error: %s", msg, err)
		return false
	}
	curr, _ := ui.pages.GetFrontPage()
	ui.sendMessage(pageName(curr), wm_killfocus, page)
	ui.pageHistory.PushBack(pageName(curr))

	switch msg.message {
	case wm_show:
		ui.pages.ShowPage(page)
	case wm_switch:
		ui.pages.SwitchToPage(page)
	default:
		panic("invalid call")
	}
	ui.pages.SendToFront(page)
	return true
}

func (ui *UI) App() *tview.Application {
	return ui.app
}

func (ui *UI) lastPage() pageName {
	el := ui.pageHistory.Back()
	if el == nil {
		return ""
	}
	return el.Value.(pageName)
}

func asPage(v any) (string, error) {
	switch val := v.(type) {
	case pageName:
		return string(val), nil
	case string:
		return val, nil
	default:
		return "", fmt.Errorf("invalid page type: %T", val)
	}
}

func (ui *UI) destroy() {
	ui.app.Stop()
	ui.log.Println("terminating")

	// close log
	ui.log.SetOutput(os.Stderr)
	ui.lf.Close()
}

func lines(w io.Writer, lines []string) {
	for _, line := range lines {
		fmt.Fprintln(w, colorize(line))
	}
}

func linesEnum(w io.Writer, items []string) {
	for i, line := range items {
		fmt.Fprintln(w, colorize(fmt.Sprintf("[$ptc]%2d.[-]   %s", i+1, line)))
	}
}

func maxLineLen(lines []string) int {
	var maxLen = 0
	for _, l := range lines {
		if lineLen := len(colorize(l)); lineLen > maxLen {
			maxLen = lineLen
		}
	}
	return maxLen
}
