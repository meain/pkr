package main

import (
	"fmt"
	"log"
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

const (
	width  = 700
	height = 300
)

func getInput() []Completion {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return []Completion{}
	}

	data, err := io.ReadAll(io.Reader(os.Stdin))
	if err != nil {
		fmt.Fprintln(os.Stderr, "reading stdin:", err)
		os.Exit(1)
	}

	lines := strings.Split(string(data), "\n")

	completions := []Completion{}
	for _, line := range lines {
		completions = append(completions, Completion{text: line, selectable: true})
	}

	return completions
}

func filterItems(items []Completion, subs string) []Completion {
	subs = strings.ToLower(subs)
	result := []Completion{}

	for _, item := range items {
		if strings.Contains(strings.ToLower(item.text), subs) {
			result = append(result, item)
		}
	}

	return result
}

func convertToInterface(completions []Completion) []interface{} {
	result := make([]interface{}, len(completions))
	for i, v := range completions {
		result[i] = v
	}

	return result
}

func performAction(comp Completion) {
	if comp.action.Command == "" {
		fmt.Println(comp.text)
		return
	}

	cmd := comp.action.Command
	var err error

	fmt.Println("Performing:", cmd, comp.text)

	switch comp.action.Mode {
	case "stdin":
		_, err = commandOutputWithStdin(cmd, comp.text)
	case "args":
		_, err = commandOutputWithArgs(cmd, comp.text)
	default:
		if strings.Contains(comp.action.Command, "{{}}") {
			cmd = strings.ReplaceAll(cmd, "{{}}", comp.text)
		}

		splits := strings.Split(comp.text, " ")
		for i, split := range splits {
			subs := fmt.Sprintf("{{%d}}", i+1)
			if strings.Contains(comp.action.Command, subs) {
				cmd = strings.ReplaceAll(cmd, subs, split)
			}
		}

		_, err = commandOutputWithArgs(cmd, "")
	}

	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	initial := getInput()
	rc := map[string]Completer{}
	var err error

	fp := "pkr.yaml"
	if len(os.Args) > 1 {
		fp = os.Args[1]
	}

	if len(initial) == 0 {
		rc, err = getRegisteredCompleters(fp)
		if err != nil {
			log.Fatal(err)
		}
	}

	var myWindow fyne.Window

	myApp := app.New()
	myApp.Settings().SetTheme(CustomTheme())

	// Make a borderless window
	drv := myApp.Driver()
	if drv, ok := drv.(desktop.Driver); ok {
		myWindow = drv.CreateSplashWindow()
	}

	myWindow.Resize(fyne.NewSize(width, height))
	myWindow.CenterOnScreen()
	myWindow.SetPadded(true)

	sl := []Completion{}
	hasInput := false

	if len(initial) > 0 {
		hasInput = true
		sl = initial
	} else {
		sl = getCompletions("", rc)
	}

	isl := convertToInterface(sl)
	data := binding.BindUntypedList(&isl)

	list := widget.NewListWithData(
		data,
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			obj, _ := i.(binding.Untyped).Get()
			o.(*widget.Label).SetText(obj.(Completion).text)
		})

	list.OnSelected = func(id int) {
		items, err := data.Get()
		if err != nil {
			fmt.Println(err)
		}

		if len(items) > 0 {
			performAction(items[id].(Completion))
		}

		myWindow.Close()
	}

	// make list update as soon as data changes via input
	data.AddListener(binding.NewDataListener(func() {
		list.Refresh()
	}))

	input := NewSearchField(myWindow, list)
	input.SetPlaceHolder("Type you little maniac...")
	input.OnChanged = func(s string) {
		completions := []Completion{}
		if hasInput {
			completions = filterItems(initial, s)
		} else {
			completions = getCompletions(s, rc)
		}

		// FIXME: somehow a simple set is not working
		data.Set(convertToInterface([]Completion{}))
		data.Set(convertToInterface(completions))
	}

	input.OnSubmitted = func(s string) {
		items, err := data.Get()
		if err != nil {
			log.Fatal(err)
		}

		if len(items) > 0 {
			performAction(items[0].(Completion))
			myWindow.Close()
		}
	}

	cont := container.NewBorder(input, nil, nil, nil, list)
	cont.Resize(fyne.NewSize(width, height))

	myWindow.SetContent(cont)
	myWindow.Canvas().Focus(input)
	myWindow.ShowAndRun()
}
