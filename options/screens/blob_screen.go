package screens

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/tr1sm0s1n/project-wallet-x/api"
	"github.com/tr1sm0s1n/project-wallet-x/config"
)

func BlobScreen(w fyne.Window) fyne.CanvasObject {
	keyEntry := widget.NewPasswordEntry()
	keyEntry.SetPlaceHolder("Private Key")
	keyEntry.Validator = validation.NewRegexp(`^[0-9a-fA-F]{64}$`, "Not a valid key.")

	receiverEntry := widget.NewEntry()
	receiverEntry.SetPlaceHolder("Ox...")
	receiverEntry.Validator = validation.NewRegexp(`^0x[0-9a-fA-F]{40}$`, "Not a valid address.")

	blobEntry := widget.NewEntry()
	blobEntry.SetPlaceHolder("Hello, World!")

	loading := widget.NewProgressBarInfinite()
	loading.Hide()

	var form *widget.Form
	form = &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Key", Widget: keyEntry, HintText: "Your private key."},
			{Text: "Address", Widget: receiverEntry, HintText: "Address of the receiver."},
			{Text: "Blob", Widget: blobEntry, HintText: "Blob data."},
		},
		OnCancel: func() {
			keyEntry.SetText("")
			receiverEntry.SetText("")
			blobEntry.SetText("")
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

			if err = api.BlobTx(client, keyEntry.Text, receiverEntry.Text, blobEntry.Text); err != nil {
				loading.Hide()
				form.Enable()
				dialog.ShowInformation("Error", err.Error(), w)
				return
			}

			loading.Hide()
			form.Enable()
			dialog.ShowInformation("Success", "Transaction succeeded!", w)
		},
	}
	return container.NewVBox(
		form,
		loading,
	)
}
