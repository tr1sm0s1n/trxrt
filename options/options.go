package options

import (
	"fyne.io/fyne/v2"
	"github.com/tr1sm0s1n/project-wallet-x/options/screens"
)

type Option struct {
	Title, Intro string
	View         func(w fyne.Window) fyne.CanvasObject
}

var (
	Options = map[string]Option{
		"home": {"Home", "", screens.HomeScreen},
		"legacy": {"Legacy TX",
			"Send a type 0x0 transaction.",
			screens.LegacyScreen,
		},
	}

	OptionIndex = map[string][]string{
		"": {"home", "legacy"},
	}
)
