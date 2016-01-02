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
		ticketQueryPage.PreviousLine(1)
		ticketQueryPage.Update()
	case ticketList:
		ticketListPage.PreviousLine(1)
		ticketListPage.Update()
	case ticketShow:
		ticketShowPage.PreviousLine(1)
		ticketShowPage.Update()
	}
}

func handleDownKey() {
	switch currentPage {
	case ticketQuery:
		ticketQueryPage.NextLine(1)
		ticketQueryPage.Update()
	case ticketList:
		ticketListPage.NextLine(1)
		ticketListPage.Update()
	case ticketShow:
		ticketShowPage.NextLine(1)
		ticketShowPage.Update()
	}
}

func handlePageUpKey() {
	switch currentPage {
	case ticketQuery:
		ticketQueryPage.PreviousPage()
		ticketQueryPage.Update()
	case ticketList:
		ticketListPage.PreviousPage()
		ticketListPage.Update()
	case ticketShow:
		ticketShowPage.PreviousPage()
		ticketShowPage.Update()
	}
}

func handlePageDownKey() {
	switch currentPage {
	case ticketQuery:
		ticketQueryPage.NextPage()
		ticketQueryPage.Update()
	case ticketList:
		ticketListPage.NextPage()
		ticketListPage.Update()
	case ticketShow:
		ticketShowPage.NextPage()
		ticketShowPage.Update()
	}
}
