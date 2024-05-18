package main

import (
	"fmt"
	"os"
	"strings"

	"io"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

func getInput() []string {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return []string{}
	}

	data, err := io.ReadAll(io.Reader(os.Stdin))
	if err != nil {
		fmt.Fprintln(os.Stderr, "reading stdin:", err)
		os.Exit(1)
	}

	lines := strings.Split(string(data), "\n")
	return lines
}

func filterItems(items []string, subs string) []string {
	subs = strings.ToLower(subs)
	result := []string{}

	for _, item := range items {
		if strings.Contains(strings.ToLower(item), subs) {
			result = append(result, item)
		}
	}

	return result
}

func main() {
	var myWindow fyne.Window

	myApp := app.New()

	drv := myApp.Driver()
	if drv, ok := drv.(desktop.Driver); ok {
		myWindow = drv.CreateSplashWindow()
	}

	myWindow.Resize(fyne.NewSize(500, 300))
	myWindow.CenterOnScreen()
	myWindow.SetPadded(true)

	// FIXME: not working
	// myWindow.Canvas().SetOnTypedKey(func(ev *fyne.KeyEvent) {
	// 	if ev.Name == fyne.KeyEscape {
	// 		myWindow.Close()
	// 	}
	// })

	sl := []string{}
	hasInput := false
	initial := getInput()

	if len(initial) > 0 {
		hasInput = true

		for _, line := range initial {
			sl = append(sl, line)
		}
	}

	data := binding.BindStringList(&sl)

	input := widget.NewEntry()
	input.SetPlaceHolder("Type you little maniac...")
	input.OnChanged = func(s string) {
		// fmt.Println(s)

		completions := []string{}
		if hasInput {
			completions = filterItems(initial, s)
		} else {
			completions = getCompletions(s)
		}

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
	list.Resize(fyne.NewSize(500, 400))

	list.OnSelected = func(id int) {
		items, err := data.Get()
		if err != nil {
			fmt.Println(err)
		}

		if len(items) > 0 {
			fmt.Println(items[id])
		}

		myWindow.Close()
	}

	// make list update as soon as data changes via input
	data.AddListener(binding.NewDataListener(func() {
		// fmt.Println("refreshing...")
		list.Refresh()
		list.Resize(fyne.NewSize(500, 400))
	}))

	cont := container.NewVBox(input, list)
	myWindow.SetContent(cont)
	myWindow.Canvas().Focus(input)
	myWindow.ShowAndRun()
}
