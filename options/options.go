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
		"legacy": {"Legacy Tx",
			"Send a type 0x0 transaction.",
			screens.LegacyScreen,
		},
		"dynamic_fee": {"Dynamic Fee Tx",
			"Send a type 0x2 transaction.",
			screens.DynamicFeeScreen,
		},
		"blob": {"Blob Tx",
			"Send a type 0x3 transaction.",
			screens.BlobScreen,
		},
	}

	OptionIndex = map[string][]string{
		"": {"home", "legacy", "dynamic_fee", "blob"},
	}
)
