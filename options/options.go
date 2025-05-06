package options

import (
	"fyne.io/fyne/v2"
	"github.com/tr1sm0s1n/trxrt/options/screens"
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
		"set_code": {"Set Code Tx",
			"Send a type 0x4 transaction.",
			screens.SetCodeScreen,
		},
		"deploy": {"Deploy Tx",
			"Deploy a contract.",
			screens.DeployScreen,
		},
	}

	OptionIndex = map[string][]string{
		"": {"home", "legacy", "dynamic_fee", "blob", "set_code", "deploy"},
	}
)
