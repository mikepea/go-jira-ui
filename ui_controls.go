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
	switch currentPage.(type) {
	case *QueryPage:
		currentPage.PreviousLine(1)
		currentPage.Update()
	case *TicketListPage:
		currentPage.PreviousLine(1)
		currentPage.Update()
	case *LabelListPage:
		currentPage.PreviousLine(1)
		currentPage.Update()
	case *TicketShowPage:
		currentPage.PreviousLine(1)
		currentPage.Update()
	}
}

func handleDownKey() {
	switch currentPage.(type) {
	case *QueryPage:
		currentPage.NextLine(1)
		currentPage.Update()
	case *TicketListPage:
		currentPage.NextLine(1)
		currentPage.Update()
	case *LabelListPage:
		currentPage.NextLine(1)
		currentPage.Update()
	case *TicketShowPage:
		currentPage.NextLine(1)
		currentPage.Update()
	}
}

func handlePageUpKey() {
	switch currentPage.(type) {
	case *QueryPage:
		currentPage.PreviousPage()
		currentPage.Update()
	case *TicketListPage:
		currentPage.PreviousPage()
		currentPage.Update()
	case *LabelListPage:
		currentPage.PreviousPage()
		currentPage.Update()
	case *TicketShowPage:
		currentPage.PreviousPage()
		currentPage.Update()
	}
}

func handlePageDownKey() {
	switch currentPage.(type) {
	case *QueryPage:
		currentPage.NextPage()
		currentPage.Update()
	case *TicketListPage:
		currentPage.NextPage()
		currentPage.Update()
	case *LabelListPage:
		currentPage.NextPage()
		currentPage.Update()
	case *TicketShowPage:
		currentPage.NextPage()
		currentPage.Update()
	}
}
