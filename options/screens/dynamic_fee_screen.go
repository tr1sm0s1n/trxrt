package screens

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/tr1sm0s1n/project-wallet-x/api"
	"github.com/tr1sm0s1n/project-wallet-x/config"
)

func DynamicFeeScreen(w fyne.Window) fyne.CanvasObject {
	keyEntry := widget.NewPasswordEntry()
	keyEntry.SetPlaceHolder("Private Key")
	keyEntry.Validator = validation.NewRegexp(`^[0-9a-fA-F]{64}$`, "Not a valid key.")

	receiverEntry := widget.NewEntry()
	receiverEntry.SetPlaceHolder("Ox...")
	receiverEntry.Validator = validation.NewRegexp(`^0x[0-9a-fA-F]{40}$`, "Not a valid address.")

	amountEntry := widget.NewEntry()
	amountEntry.SetPlaceHolder("1")
	amountEntry.Validator = validation.NewRegexp(`^[1-9][0-9]{0,8}$`, "Not a valid amount.")

	loading := widget.NewProgressBarInfinite()
	loading.Hide()

	var form *widget.Form
	form = &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Key", Widget: keyEntry, HintText: "Your private key."},
			{Text: "Address", Widget: receiverEntry, HintText: "Address of the receiver."},
			{Text: "Amount", Widget: amountEntry, HintText: "Tranfer amount in ETH."},
		},
		OnCancel: func() {
			keyEntry.SetText("")
			receiverEntry.SetText("")
			amountEntry.SetText("")
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

			amountInt, _ := strconv.Atoi(amountEntry.Text)
			if err = api.DynamicFeeTx(client, keyEntry.Text, receiverEntry.Text, int64(amountInt)); err != nil {
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
