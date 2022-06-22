package tui

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var allColors = make(map[tcell.Color]string, len(tcell.ColorNames))

func init() {
	for k, v := range tcell.ColorNames {
		allColors[v] = k
	}
}

// https://www.ditig.com/256-colors-cheat-sheet
var (
	themeLotus1 = tview.Theme{
		PrimitiveBackgroundColor:    tcell.ColorTeal,
		ContrastBackgroundColor:     tcell.ColorBlack,
		MoreContrastBackgroundColor: tcell.ColorBlue,
		BorderColor:                 tcell.ColorBlack,
		TitleColor:                  tcell.ColorBlack,
		GraphicsColor:               tcell.ColorMaroon,
		PrimaryTextColor:            tcell.ColorWhite,
		SecondaryTextColor:          tcell.ColorBlack,
		TertiaryTextColor:           tcell.ColorGrey,
		InverseTextColor:            tcell.ColorBlack,
		ContrastSecondaryTextColor:  tcell.ColorDarkBlue,
	}
)

func colorName(c tcell.Color) string {
	return allColors[c]
}

func initReplacer(th tview.Theme) *strings.Replacer {
	return strings.NewReplacer(
		"$pbc", colorName(th.PrimitiveBackgroundColor),
		"$cbc", colorName(th.ContrastBackgroundColor),
		"$mcbc", colorName(th.MoreContrastBackgroundColor),
		"$bc", colorName(th.BorderColor),
		"$tc", colorName(th.TitleColor),
		"$gc", colorName(th.GraphicsColor),
		"$ptc", colorName(th.PrimaryTextColor),
		"$stc", colorName(th.SecondaryTextColor),
		"$ttc", colorName(th.TertiaryTextColor),
		"$itc", colorName(th.InverseTextColor),
		"$cstc", colorName(th.ContrastSecondaryTextColor),
	)
}

func colorize(text string) string {
	colorReplacer := initReplacer(tview.Styles)
	return colorReplacer.Replace(text)
}
