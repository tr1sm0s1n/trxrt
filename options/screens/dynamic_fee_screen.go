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

func DynamicFeeScreen(w fyne.Window) fyne.CanvasObject {
	keyEntry := widget.NewPasswordEntry()
	keyEntry.SetPlaceHolder("Private Key")
	keyEntry.Validator = validation.NewRegexp(`^[0-9a-fA-F]{64}$`, "Not a valid key.")

	receiverEntry := widget.NewEntry()
	receiverEntry.SetPlaceHolder("0x...")
	receiverEntry.Validator = validation.NewRegexp(`^0x[0-9a-fA-F]{40}$`, "Not a valid address.")

	amountEntry := widget.NewEntry()
	amountEntry.SetPlaceHolder("1")
	amountEntry.Validator = validation.NewRegexp(`^[1-9][0-9]{0,8}$`, "Not a valid amount.")

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
			{Text: "Key", Widget: keyEntry, HintText: "Your private key."},
			{Text: "Address", Widget: receiverEntry, HintText: "Address of the receiver."},
			{Text: "Amount", Widget: amountEntry, HintText: "Tranfer amount in ETH."},
			{Text: "Gas", Widget: gasEntry, HintText: "Gas limit in wei."},
			{Widget: gasSlide},
			{Text: "Max Fee", Widget: maxFeeEntry, HintText: "maxFeePerGas in gwei."},
			{Widget: maxFeeSlide},
			{Text: "Max Priority Fee", Widget: maxPriorityFeeEntry, HintText: "maxPriorityFeePerGas in wei."},
			{Widget: maxPriorityFeeSlide},
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
			if err = api.DynamicFeeTx(client, keyEntry.Text, receiverEntry.Text, int64(amountInt), g, mf, mpf); err != nil {
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
