package main

import (
	_ "embed"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

//go:embed VictorMono-Regular.ttf
var fontFile []byte

type myTheme struct {
	fyne.Theme
}

func CustomTheme() fyne.Theme {
	return &myTheme{Theme: fyne.CurrentApp().Settings().Theme()}
}

func (m myTheme) Font(style fyne.TextStyle) fyne.Resource {
	return fyne.NewStaticResource("Victor Mono", fontFile)
}

func (m myTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == theme.ColorNameBackground {
		if variant == theme.VariantLight {
			return color.White
		}
		return color.Black
	}

	return theme.DefaultTheme().Color(name, variant)
}
