package screens

import (
	"unicode"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func LegacyScreen(w fyne.Window) fyne.CanvasObject {
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
	amountEntry.SetPlaceHolder("Amount in wei")
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

	submitButton := widget.NewButton("Submit", func() {
		key := keyEntry.Text
		receiver := receiverEntry.Text
		amount := amountEntry.Text

		if key == "" || receiver == "" || amount == "" {
			dialog.ShowInformation("Error", "Please fill all fields.", w)
			return
		}

		dialog.ShowInformation("Success", "Transaction succeeded!", w)
	})

	return container.NewVBox(
		widget.NewLabel("Key:"),
		keyEntry,
		widget.NewLabel("Address:"),
		receiverEntry,
		widget.NewLabel("Transfer:"),
		amountEntry,
		submitButton,
	)
}
