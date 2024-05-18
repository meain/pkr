package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("pkr")
	myWindow.Resize(fyne.NewSize(500, 500))

	sl := []string{}
	data := binding.BindStringList(&sl)

	input := widget.NewEntry()
	input.SetPlaceHolder("Type you little maniac...")
	input.OnChanged = func(s string) {
		// fmt.Println(s)
		completions := getCompletions(s)

		// FIXME: somehow a simple set is not working
		data.Set([]string{})
		data.Set(completions)
	}

	input.OnSubmitted = func(s string) {
		items, err := data.Get()
		if err != nil {
			fmt.Println(err)
			return
		}

		if len(items) > 0 {
			fmt.Println(items[0])
		}

		myWindow.Close()
	}

	list := widget.NewListWithData(
		data,
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			str, _ := i.(binding.String).Get()
			o.(*widget.Label).SetText(str)
		})
	list.Resize(fyne.NewSize(500, 300))

	// make list update as soon as data changes via input
	data.AddListener(binding.NewDataListener(func() {
		// fmt.Println("refreshing...")
		list.Refresh()
	}))

	cont := container.NewVBox(input, list)
	myWindow.SetContent(cont)
	myWindow.ShowAndRun()
}
