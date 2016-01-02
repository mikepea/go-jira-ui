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

func handleBackKey() {
	switch currentPage.(type) {
	case *QueryPage:
		ui.StopLoop()
		exitNow = true
	case *TicketListPage:
		previousPage = currentPage
		currentPage = &ticketQueryPage
	case *LabelListPage:
		previousPage = currentPage
		currentPage = &ticketListPage
	case *TicketShowPage:
		previousPage = currentPage
		currentPage = &ticketListPage
	}
	changePage()
}

func handleResize() {
	changePage()
}

func handleSelectKey() {
	switch currentPage.(type) {
	case *QueryPage:
		previousPage = currentPage
		currentPage = &ticketListPage
	case *TicketListPage:
		previousPage = currentPage
		currentPage = &ticketShowPage
	}
	changePage()
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
