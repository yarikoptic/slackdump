package tui

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	logo = []string{
		"",
		"   [$ptc:$cbc]         [-:-]",
		"   [$ptc:$cbc]    ▲    [-:-]",
		"   [$ptc:$cbc]   ▲ ▲   [-:-]",
		"   [$ptc:$cbc]         [-:-]",
		"   [$itc]Slackdump",
	}
	logoSz       = maxLineLen(logo)
	headerHeight = len(logo) + 1
)

func modal(p tview.Primitive, width int, height int) tview.Primitive {
	grid := tview.NewGrid().
		SetColumns(0, width, 0).
		SetRows(0, height, 0).
		AddItem(p, 1, 1, 1, 1, 0, 0, true)

	return grid
}

type shadowOpts struct {
	HFillChar       string
	VFillChar       string
	FillCount       int
	Color           tcell.Color
	BackgroundColor tcell.Color
}

// shadow creates a drop shadow for a primitive.  If opts are nil
// defaults will be used.  opts are not verified for errors.
func shadow(p tview.Primitive, opts *shadowOpts) tview.Primitive {
	//   0        1              2
	//  +--+-------------------+
	//  |          ^           |    0
	//  +  <spans 2 row/col>   +--+
	//  |          v           |  | 1
	//  +--+-------------------+--+
	//     |      < spans 2 col > | 2
	//     +-------------------+--+
	if opts == nil {
		opts = &shadowOpts{
			HFillChar:       "▒",
			VFillChar:       "▒",
			FillCount:       120,
			Color:           tcell.ColorGray,
			BackgroundColor: tview.Styles.PrimaryTextColor,
		}
	}

	vertShadow := tview.NewTextView().
		SetTextColor(opts.Color).
		SetText(strings.Repeat(opts.VFillChar, opts.FillCount))
	horShadow := tview.NewTextView().
		SetTextColor(opts.Color).
		SetText(strings.Repeat(opts.HFillChar, opts.FillCount))
	item := tview.NewGrid().
		SetColumns(2, 0, 2).
		SetRows(1, 0, 1).
		AddItem(p, 0, 0, 2, 2, 0, 0, true).           // body
		AddItem(vertShadow, 1, 2, 1, 1, 0, 0, false). // vertical shadow
		AddItem(horShadow, 2, 1, 1, 2, 0, 0, false)   // horizontal shadow
	return item
}

func makeInstructions(lines []string) *tview.TextView {
	p := tview.NewTextView().
		SetDynamicColors(true).
		SetWordWrap(true).
		SetRegions(true)
	p.SetTextColor(tview.Styles.PrimitiveBackgroundColor).
		SetBackgroundColor(tview.Styles.ContrastBackgroundColor).
		SetBorder(true).
		SetBorderColor(tview.Styles.GraphicsColor)

	linesEnum(p, lines)

	return p
}

// makeHeader creates a logo
func makeHeader(title string) tview.Primitive {
	//    0              1                  2
	//  +----+-----------------------+----+
	//  | LO |_______________________|    | 0
	//  | .. |________TEXT___________|    | 1
	//  | GO |                       |    | 2
	//  +----+-----------------------+----+
	//
	tLogo := tview.NewTextView().SetDynamicColors(true)
	lines(tLogo, logo)
	tTitle := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetText(title)
	grid := tview.NewGrid().
		SetColumns(logoSz, 0, logoSz).
		SetRows(0, 1, 0).
		// column 0
		AddItem(tLogo, 0, 0, 3, 1, 0, 0, false).
		// column 1
		AddItem(tview.NewBox(), 0, 1, 1, 1, 0, 0, false).
		AddItem(tTitle, 1, 1, 1, 1, 1, 1, false).
		AddItem(tview.NewBox(), 2, 1, 1, 1, 0, 0, false).
		// column 2
		AddItem(tview.NewBox(), 0, 2, 3, 1, 0, 0, false)
	return grid
}
