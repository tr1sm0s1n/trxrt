// Package main provides various examples of Fyne API capabilities.
package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/tr1sm0s1n/project-wallet-x/options"
)

const preferenceCurrentOption = "currentOption"

var topWindow fyne.Window

func main() {
	a := app.NewWithID("io.wallet.x")
	// a.SetIcon(data.FyneLogo)
	makeTray(a)
	logLifecycle(a)
	w := a.NewWindow("Wallet-X")
	topWindow = w

	w.SetMaster()

	content := container.NewStack()
	title := widget.NewLabel("Transaction Type")
	intro := widget.NewLabel("Choose your transaction type\namong legacy, access list, dynamic fee, blob and set code.")
	intro.Wrapping = fyne.TextWrapWord
	setOption := func(t options.Option) {
		if fyne.CurrentDevice().IsMobile() {
			child := a.NewWindow(t.Title)
			topWindow = child
			child.SetContent(t.View(topWindow))
			child.Show()
			child.SetOnClosed(func() {
				topWindow = w
			})
			return
		}

		title.SetText(t.Title)
		intro.SetText(t.Intro)

		if t.Title == "Welcome" {
			title.Hide()
			intro.Hide()
		} else {
			title.Show()
			intro.Show()
		}

		content.Objects = []fyne.CanvasObject{t.View(w)}
		content.Refresh()
	}

	option := container.NewBorder(
		container.NewVBox(title, widget.NewSeparator(), intro), nil, nil, nil, content)
	if fyne.CurrentDevice().IsMobile() {
		w.SetContent(makeNav(setOption, false, w))
	} else {
		split := container.NewHSplit(makeNav(setOption, true, w), option)
		split.Offset = 0.2
		w.SetContent(split)
	}
	w.Resize(fyne.NewSize(640, 460))
	w.ShowAndRun()
}

func logLifecycle(a fyne.App) {
	a.Lifecycle().SetOnStarted(func() {
		log.Println("Lifecycle: Started")
	})
	a.Lifecycle().SetOnStopped(func() {
		log.Println("Lifecycle: Stopped")
	})
	a.Lifecycle().SetOnEnteredForeground(func() {
		log.Println("Lifecycle: Entered Foreground")
	})
	a.Lifecycle().SetOnExitedForeground(func() {
		log.Println("Lifecycle: Exited Foreground")
	})
}

func makeTray(a fyne.App) {
	if desk, ok := a.(desktop.App); ok {
		h := fyne.NewMenuItem("Wallet-X", func() {})
		h.Icon = theme.HomeIcon()
		menu := fyne.NewMenu("Wallet-X", h)
		h.Action = func() {
			log.Println("System tray menu tapped")
			h.Label = "Welcome"
			menu.Refresh()
		}
		desk.SetSystemTrayMenu(menu)
	}
}

func makeNav(setOption func(option options.Option), loadPrevious bool, w fyne.Window) fyne.CanvasObject {
	a := fyne.CurrentApp()

	tree := &widget.Tree{
		ChildUIDs: func(uid string) []string {
			return options.OptionIndex[uid]
		},
		IsBranch: func(uid string) bool {
			children, ok := options.OptionIndex[uid]

			return ok && len(children) > 0
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			return widget.NewLabel("Collection Widgets")
		},
		UpdateNode: func(uid string, branch bool, obj fyne.CanvasObject) {
			t, ok := options.Options[uid]
			if !ok {
				fyne.LogError("Missing option panel: "+uid, nil)
				return
			}
			obj.(*widget.Label).SetText(t.Title)
		},
		OnSelected: func(uid string) {
			if t, ok := options.Options[uid]; ok {
				a.Preferences().SetString(preferenceCurrentOption, uid)
				setOption(t)
			}
		},
	}

	if loadPrevious {
		currentPref := a.Preferences().StringWithFallback(preferenceCurrentOption, "welcome")
		tree.Select(currentPref)
	}

	themes := container.NewGridWithColumns(2,
		widget.NewButton("Dark", func() {
			a.Settings().SetTheme(&forcedVariant{Theme: theme.DefaultTheme(), variant: theme.VariantDark})
		}),
		widget.NewButton("Light", func() {
			a.Settings().SetTheme(&forcedVariant{Theme: theme.DefaultTheme(), variant: theme.VariantLight})
		}),
	)

	bottom := container.NewVBox(themes, widget.NewButton("Fullscreen", func() {
		w.SetFullScreen(!w.FullScreen())
	}))

	return container.NewBorder(nil, bottom, nil, nil, tree)
}
