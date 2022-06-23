package tui

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

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
