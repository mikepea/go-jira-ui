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
	ui.Handle("/sys/kbd/<enter>", func(ui.Event) {
		handleSelectKey()
	})
	ui.Handle("/sys/wnd/resize", func(ui.Event) {
		handleResize()
	})
}

func handleBackKey() {
	switch currentPage {
	case ticketQuery:
		ui.StopLoop()
		exitNow = true
	case ticketList:
		previousPage = currentPage
		currentPage = ticketQuery
	case ticketShow:
		previousPage = currentPage
		currentPage = ticketList
	}
	changePage()
}

func handleResize() {
	changePage()
}

func handleSelectKey() {
	switch currentPage {
	case ticketQuery:
		currentPage = ticketList
		previousPage = ticketQuery
	case ticketList:
		currentPage = ticketShow
		previousPage = ticketList
	}
	changePage()
}

func handleUpKey() {
	switch currentPage {
	case ticketQuery:
		prevQuery(1)
		updateQueryPage(activeQueryList)
	case ticketList:
		prevTicket(1)
		updateTicketListPage(activeTicketListList)
	case ticketShow:
		prevTicketLine(1)
		updateTicketShowPage(activeTicketShowList)
	}
}

func handleDownKey() {
	switch currentPage {
	case ticketQuery:
		nextQuery(1)
		updateQueryPage(activeQueryList)
	case ticketList:
		nextTicket(1)
		updateTicketListPage(activeTicketListList)
	case ticketShow:
		nextTicketLine(1)
		updateTicketShowPage(activeTicketShowList)
	}
}
