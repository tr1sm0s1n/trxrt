package screens

import (
	"strconv"
	"unicode"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/tr1sm0s1n/project-wallet-x/api"
	"github.com/tr1sm0s1n/project-wallet-x/config"
)

func DynamicFeeScreen(w fyne.Window) fyne.CanvasObject {
	keyEntry := widget.NewEntry()
	keyEntry.SetPlaceHolder("Enter your private key")
	keyEntry.OnChanged = func(text string) {
		if len(text) > 64 {
			keyEntry.SetText(text[:64]) // Trim if exceeds limit
		}
	}

	receiverEntry := widget.NewEntry()
	receiverEntry.SetPlaceHolder("Enter receiver address")
	receiverEntry.OnChanged = func(text string) {
		if len(text) > 42 {
			receiverEntry.SetText(text[:42]) // Trim if exceeds limit
		}
	}

	amountEntry := widget.NewEntry()
	amountEntry.SetPlaceHolder("Amount in ETH")
	amountEntry.OnChanged = func(text string) {
		filtered := ""
		for _, r := range text {
			if unicode.IsDigit(r) { // Allow only digits
				filtered += string(r)
			}
		}
		if text != filtered {
			amountEntry.SetText(filtered)
		}
	}

	loading := widget.NewProgressBarInfinite()
	loading.Hide()

	var submitButton *widget.Button
	submitButton = widget.NewButton("Submit", func() {
		key := keyEntry.Text
		receiver := receiverEntry.Text
		amount := amountEntry.Text

		loading.Show()
		submitButton.Disable()

		if key == "" || receiver == "" || amount == "" {
			dialog.ShowInformation("Error", "Please fill all fields.", w)
			return
		}

		client, err := config.DialClient()
		if err != nil {
			dialog.ShowInformation("Error", err.Error(), w)
			return
		}

		amountInt, _ := strconv.Atoi(amount)

		if err = api.DynamicFeeTx(client, key, receiver, int64(amountInt)); err != nil {
			dialog.ShowInformation("Error", err.Error(), w)
			return
		}

		loading.Hide()
		submitButton.Enable()
		dialog.ShowInformation("Success", "Transaction succeeded!", w)
		keyEntry.SetText("")
		receiverEntry.SetText("")
		amountEntry.SetText("")
	})

	return container.NewVBox(
		widget.NewLabel("Key:"),
		keyEntry,
		widget.NewLabel("Address:"),
		receiverEntry,
		widget.NewLabel("Transfer:"),
		amountEntry,
		submitButton,
		loading,
	)
}
