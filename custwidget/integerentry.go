package custwidget

import (
	"strconv"

	"fyne.io/fyne/v2"
	// "fyne.io/fyne/v2/app"
	// "fyne.io/fyne/v2/driver/mobile"
	"fyne.io/fyne/v2/widget"
)

type IntegerEntry struct {
	widget.Entry
}

func NewIntegerEntry() *IntegerEntry {
	entry := &IntegerEntry{}
	entry.ExtendBaseWidget(entry)
	return entry
}

func (e *IntegerEntry) TypedRune(r rune) {
	if r >= '0' && r <= '9' {
		e.Entry.TypedRune(r)
	}
}

func (e *IntegerEntry) TypedShortcut(shortcut fyne.Shortcut) {
	paste, ok := shortcut.(*fyne.ShortcutPaste)
	if !ok {
		e.Entry.TypedShortcut(shortcut)
		return
	}
	content := paste.Clipboard.Content()
	if _, err := strconv.ParseInt(content, 2, 64); err == nil {
		e.Entry.TypedShortcut(shortcut)
	}
}
