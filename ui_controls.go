package main

import (
	ui "github.com/gizak/termui"
)

func registerKeyboardHandlers() {
	ui.Handle("/sys/kbd/q", func(ui.Event) {
		handleBackKey()
	})
	ui.Handle("/sys/kbd/j", func(ui.Event) {
		handleDownKey()
	})
	ui.Handle("/sys/kbd/k", func(ui.Event) {
		handleUpKey()
	})
	ui.Handle("/sys/kbd/<down>", func(ui.Event) {
		handleDownKey()
	})
	ui.Handle("/sys/kbd/<up>", func(ui.Event) {
		handleUpKey()
	})
	ui.Handle("/sys/kbd/L", func(ui.Event) {
		handleLabelViewKey()
	})
	ui.Handle("/sys/kbd/<enter>", func(ui.Event) {
		handleSelectKey()
	})
	ui.Handle("/sys/kbd/<space>", func(ui.Event) {
		handlePageDownKey()
	})
	ui.Handle("/sys/kbd/C-f", func(ui.Event) {
		handlePageDownKey()
	})
	ui.Handle("/sys/kbd/C-b", func(ui.Event) {
		handlePageUpKey()
	})
	ui.Handle("/sys/kbd/E", func(ui.Event) {
		handleEditKey()
	})
	ui.Handle("/sys/kbd/C", func(ui.Event) {
		handleCommentKey()
	})
	ui.Handle("/sys/wnd/resize", func(ui.Event) {
		handleResize()
	})
}

func handleLabelViewKey() {
	switch currentPage.(type) {
	case *TicketListPage:
		previousPage = currentPage
		currentPage = &labelListPage
	}
	changePage()
}

func handleEditKey() {
	if obj, ok := currentPage.(TicketEditer); ok {
		obj.EditTicket()
	}
}

func handleCommentKey() {
	if obj, ok := currentPage.(TicketCommenter); ok {
		obj.CommentTicket()
	}
}

func handleBackKey() {
	if obj, ok := currentPage.(GoBacker); ok {
		obj.GoBack()
	} else {
		ui.StopLoop()
		exitNow = true
	}
}

func handleResize() {
	changePage()
}

func handleSelectKey() {
	if obj, ok := currentPage.(ItemSelecter); ok {
		obj.SelectItem()
	}
}

func handleUpKey() {
	if obj, ok := currentPage.(PagePager); ok {
		obj.PreviousLine(1)
		obj.Update()
	}
}

func handleDownKey() {
	if obj, ok := currentPage.(PagePager); ok {
		obj.NextLine(1)
		obj.Update()
	}
}

func handlePageUpKey() {
	if obj, ok := currentPage.(PagePager); ok {
		obj.PreviousPage()
		obj.Update()
	}
}

func handlePageDownKey() {
	if obj, ok := currentPage.(PagePager); ok {
		obj.NextPage()
		obj.Update()
	}
}
