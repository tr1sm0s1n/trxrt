package screens

import (
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/tr1sm0s1n/trxrt/api"
	"github.com/tr1sm0s1n/trxrt/config"
)

func SetCodeScreen(w fyne.Window) fyne.CanvasObject {
	authKeyEntry := widget.NewPasswordEntry()
	authKeyEntry.SetPlaceHolder("Authorizer's private key")
	authKeyEntry.Validator = validation.NewRegexp(`^[0-9a-fA-F]{64}$`, "Not a valid key.")

	keyEntry := widget.NewPasswordEntry()
	keyEntry.SetPlaceHolder("Transaction signer's private key")
	keyEntry.Validator = validation.NewRegexp(`^[0-9a-fA-F]{64}$`, "Not a valid key.")

	contractEntry := widget.NewEntry()
	contractEntry.SetPlaceHolder("0x...")
	contractEntry.Validator = validation.NewRegexp(`^0x[0-9a-fA-F]{40}$`, "Not a valid contract address.")

	var g float64 = 20000
	gas := binding.BindFloat(&g)
	gasEntry := widget.NewEntryWithData(binding.FloatToString(gas))
	gasEntry.Validator = validation.NewRegexp(`^([1-9]\d{0,4}|100000)$`, "Not a valid gas limit.")
	gasEntry.OnChanged = func(t string) {
		gasEntry.SetText(strings.Split(t, ".")[0])
	}
	gasSlide := widget.NewSliderWithData(0, 100000, gas)
	gasSlide.Step = 1000

	var mf float64 = 1
	maxFee := binding.BindFloat(&mf)
	maxFeeEntry := widget.NewEntryWithData(binding.FloatToString(maxFee))
	maxFeeEntry.Validator = validation.NewRegexp(`^(10(?:\.0)?|[1-9](?:\.\d)?|0\.[1-9])$`, "Not a valid gas price.")
	maxFeeEntry.OnChanged = func(t string) {
		f, _ := strconv.ParseFloat(t, 64)
		maxFeeEntry.SetText(strconv.FormatFloat(f, 'f', -1, 64))
	}
	maxFeeSlide := widget.NewSliderWithData(0, 10, maxFee)
	maxFeeSlide.Step = 0.1

	var mpf float64 = 1000000
	maxPriorityFee := binding.BindFloat(&mpf)
	maxPriorityFeeEntry := widget.NewEntryWithData(binding.FloatToString(maxPriorityFee))
	maxPriorityFeeEntry.Validator = validation.NewRegexp(`^([1-9]\d{0,6}|10000000)$`, "Not a valid maxPriorityFee limit.")
	maxPriorityFeeEntry.OnChanged = func(t string) {
		maxPriorityFeeEntry.SetText(strings.Split(t, ".")[0])
	}
	maxPriorityFeeSlide := widget.NewSliderWithData(0, 10000000, maxPriorityFee)
	maxPriorityFeeSlide.Step = 10000

	loading := widget.NewProgressBarInfinite()
	loading.Hide()

	var form *widget.Form
	form = &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Authorization Key", Widget: authKeyEntry, HintText: "Private Key to sign authorization (balance isn't mandatory)."},
			{Text: "Transaction Key", Widget: keyEntry, HintText: "Private Key to sign transaction (balance is mandatory)."},
			{Text: "Contract Address", Widget: contractEntry, HintText: "Address of the contract."},
			{Text: "Gas", Widget: gasEntry, HintText: "Gas limit in wei."},
			{Widget: gasSlide},
			{Text: "Max Fee", Widget: maxFeeEntry, HintText: "maxFeePerGas in gwei."},
			{Widget: maxFeeSlide},
			{Text: "Max Priority Fee", Widget: maxPriorityFeeEntry, HintText: "maxPriorityFeePerGas in wei."},
			{Widget: maxPriorityFeeSlide},
		},
		OnCancel: func() {
			authKeyEntry.SetText("")
			keyEntry.SetText("")
			contractEntry.SetText("")
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

			if err = api.SetCodeTx(client, keyEntry.Text, contractEntry.Text, authKeyEntry.Text, g, mf, mpf); err != nil {
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
