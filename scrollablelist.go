package main

import (
	ui "github.com/gizak/termui"
)

// A scrollable list with a cursor. To "deactivate" the cursor, just make the
// cursor colors the same as the item colors.
type ScrollableList struct {
	ui.Block

	// The items in the list
	Items []string

	// The window's offset relative to the start of `Items`
	Offset int

	// The foreground color for non-cursor items
	ItemFgColor ui.Attribute

	// The background color for non-cursor items
	ItemBgColor ui.Attribute

	// The foreground color for the cursor
	CursorFgColor ui.Attribute

	// The background color for the cursor
	CursorBgColor ui.Attribute

	// The position of the cursor relative to the start of `Items`
	Cursor int
}

// Add an element to the list
func (sl *ScrollableList) Add(s string) {
	sl.Items = append(sl.Items, s)
	sl.render()
}

func (sl *ScrollableList) render() {
	ui.Render(sl)
}

func (sl *ScrollableList) colorsForItem(i int) (fg, bg ui.Attribute) {
	if i == sl.Cursor {
		return sl.CursorFgColor, sl.CursorBgColor
	}
	return sl.ItemFgColor, sl.ItemBgColor
}

// Implements the termui.Bufferer interface
func (sl *ScrollableList) Buffer() ui.Buffer {
	buf := sl.Block.Buffer()
	start := min(sl.Offset, len(sl.Items))
	end := min(sl.Offset+sl.InnerHeight(), len(sl.Items))
	for i, item := range sl.Items[start:end] {
		fg, bg := sl.colorsForItem(start + i)
		cells := ui.DefaultTxBuilder.Build(item, fg, bg)
		cells = ui.DTrimTxCls(cells, sl.InnerWidth())
		offsetX := 0
		for _, cell := range cells {
			width := cell.Width()
			buf.Set(
				sl.InnerBounds().Min.X+offsetX,
				sl.InnerBounds().Min.Y+i,
				cell,
			)
			offsetX += width
		}
	}
	return buf
}

// Move the window up one row
func (sl *ScrollableList) ScrollUp() {
	if sl.Offset > 0 {
		sl.Offset -= 1
		if sl.Cursor >= sl.Offset+sl.InnerHeight() {
			sl.Cursor = sl.Offset + sl.InnerHeight() - 1
		}
		sl.render()
	}
}

// Move the window down one row
func (sl *ScrollableList) ScrollDown() {
	if sl.Offset < len(sl.Items) {
		sl.Offset += 1
		if sl.Offset > sl.Cursor {
			sl.Cursor = sl.Offset
		}
		sl.render()
	}
}

// Move the cursor down one row; moving the cursor out of the window will cause
// scrolling.
func (sl *ScrollableList) CursorDown() {
	if sl.Cursor < len(sl.Items)-1 {
		sl.Cursor += 1
		if sl.Cursor > sl.Offset+sl.InnerHeight()-1 {
			sl.Offset += 1
		}
		sl.render()
	}
}

// Move the cursor up one row; moving the cursor out of the window will cause
// scrolling.
func (sl *ScrollableList) CursorUp() {
	if sl.Cursor > 0 {
		sl.Cursor -= 1
		if sl.Cursor < sl.Offset {
			sl.Offset = sl.Cursor
		}
		sl.render()
	}
}

// Move the window down one frame; this will move the cursor as well.
func (sl *ScrollableList) PageDown() {
	if sl.Offset < len(sl.Items)-sl.InnerHeight() {
		sl.Offset += sl.InnerHeight()
		if sl.Offset > sl.Cursor {
			sl.Cursor = sl.Offset
		}
		sl.render()
	}
}

// Move the window up one frame; this will move the cursor as well.
func (sl *ScrollableList) PageUp() {
	sl.Offset = max(0, sl.Offset-sl.InnerHeight())
	if sl.Cursor >= sl.Offset+sl.InnerHeight() {
		sl.Cursor = sl.Offset + sl.InnerHeight() - 1
	}
	sl.render()
}

// Scroll to the bottom of the list
func (sl *ScrollableList) ScrollToBottom() {
	if len(sl.Items) >= sl.InnerHeight() {
		sl.Offset = len(sl.Items) - sl.InnerHeight()
		sl.render()
	}
}

// Scroll to the top of the list
func (sl *ScrollableList) ScrollToTop() {
	sl.Offset = 0
	sl.render()
}
