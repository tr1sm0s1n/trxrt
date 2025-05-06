package screens

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/tr1sm0s1n/trxrt/api"
	"github.com/tr1sm0s1n/trxrt/config"
)

func DeployScreen(w fyne.Window) fyne.CanvasObject {
	keyEntry := widget.NewPasswordEntry()
	keyEntry.SetPlaceHolder("Private Key")
	keyEntry.Validator = validation.NewRegexp(`^[0-9a-fA-F]{64}$`, "Not a valid key.")

	bytecodeEntry := widget.NewMultiLineEntry()
	bytecodeEntry.SetPlaceHolder("0x...")
	bytecodeEntry.Wrapping = fyne.TextWrapWord
	bytecodeEntry.Validator = validation.NewRegexp(`^0x[0-9a-fA-F]+$`, "Not a valid bytecode.")

	loading := widget.NewProgressBarInfinite()
	loading.Hide()

	var form *widget.Form
	form = &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Key", Widget: keyEntry, HintText: "Your private key."},
			{Text: "Bytecode", Widget: bytecodeEntry, HintText: "Bytecode of the contract."},
		},
		OnCancel: func() {
			keyEntry.SetText("")
			bytecodeEntry.SetText("")
		},
		OnSubmit: func() {
			loading.Show()
			form.Disable()

			client, err := config.DialClient()
			if err != nil {
				loading.Hide()
				form.Enable()
				dialog.ShowInformation("Error", err.Error(), w)
				return
			}

			addr, err := api.DeployTx(client, keyEntry.Text, bytecodeEntry.Text)
			if err != nil {
				loading.Hide()
				form.Enable()
				dialog.ShowInformation("Error", err.Error(), w)
				return
			}

			loading.Hide()
			form.Enable()
			dialog.ShowInformation("Success", "Contract Address: "+addr.Hex(), w)
		},
	}
	return container.NewVBox(
		form,
		loading,
	)
}
