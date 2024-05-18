package main

import (
	"fmt"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func getCompletions(input string) []string {
	outs := []string{}

	if len(input) == 0 {
		return outs
	}

	for i := 1; i <= 10; i++ {
		outs = append(outs, fmt.Sprintf("%s %d", input, i))
	}

	return outs
}

func main() {
	go func() {
		window := new(app.Window)
		err := run(window)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run(window *app.Window) error {
	theme := material.NewTheme()
	var text widget.Editor
	var list widget.List
	var ops op.Ops
	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			outputs := getCompletions(text.Text())
			fmt.Println(text.Text())

			input := material.Editor(theme, &text, "Type here...")
			input.Font.Typeface = "Dank Mono"
			input.TextSize = 20
			input.Editor.SingleLine = true

			labels := []material.LabelStyle{}
			for _, out := range outputs {
				lbl := material.Label(theme, 20, out)
				lbl.Font.Typeface = "Dank Mono"
				lbl.TextSize = 20

				labels = append(labels, lbl)
			}

			lst := material.List(theme, &list)
			list.Alignment = layout.Start

			layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(input.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return lst.Layout(
						gtx,
						len(labels),
						func(gtx layout.Context, i int) layout.Dimensions {
							return labels[i].Layout(gtx)
						})
				}))

			// render
			e.Frame(gtx.Ops)
		default:
			fmt.Println(e)
		}
	}
}
