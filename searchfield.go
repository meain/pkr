package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func NewSearchField(win fyne.Window, list *widget.List) *SearchField {
	search := widget.NewEntry()
	search.SetPlaceHolder(" Search")
	return &SearchField{search, win, list}
}

type SearchField struct {
	*widget.Entry
	win  fyne.Window
	list *widget.List
}

func (s *SearchField) TypedKey(k *fyne.KeyEvent) {
	switch k.Name {
	case fyne.KeyEscape:
		s.win.Close()
	case fyne.KeyUp:
		// move to prev item in list
	case fyne.KeyDown:
		// move to next item in list
	default:
		s.Entry.TypedKey(k)
	}
}
