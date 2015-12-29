package main

import (
	"fmt"
	ui "github.com/gizak/termui"
	//jira "github.com/mikepea/go-jira"
)

const (
	ticketQuery = 1
	ticketList  = 2
	ticketShow  = 3
)

var exitNow = false

var currentPage = ticketQuery
var previousPage = ticketQuery

var ticketSelected = 0
var querySelected = 0

func prevTicket(n int) {
	ticketSelected = ticketSelected - n
}

func nextTicket(n int) {
	ticketSelected = ticketSelected + n
}

func prevQuery(n int) {
	querySelected = querySelected - n
}

func nextQuery(n int) {
	querySelected = querySelected + n
}

var origQueries = []string{
	"[0] [My Tickets]",
	"[1] [My Watched Tickets]",
	"[2] [unlabelled]",
	"[3] [Ops]",
}

var queries = []string{
	"[0] [My Tickets]",
	"[1] [My Watched Tickets]",
	"[2] [unlabelled]",
	"[3] [Ops]",
}

var origTickets = []string{
	"[0] [github.com/gizak/ui]",
	"[1] [你好，世界]",
	"[2] [こんにちは世界]",
	"[3] [color output]",
	"[4] [output.go]",
	"[5] [random_out.go]",
	"[6] [dashboard.go]",
	"[7] [nsf/termbox-go]",
}

var tickets = []string{
	"[0] [github.com/gizak/ui]",
	"[1] [你好，世界]",
	"[2] [こんにちは世界]",
	"[3] [color output]",
	"[4] [output.go]",
	"[5] [random_out.go]",
	"[6] [dashboard.go]",
	"[7] [nsf/termbox-go]",
}

func markActiveTicket() {
	for i, v := range origTickets {
		if i == ticketSelected {
			tickets[i] = v + "(fg-white,bg-blue)"
		} else {
			tickets[i] = v + "()"
		}
	}
}

func markActiveQuery() {
	for i, v := range origQueries {
		if i == querySelected {
			queries[i] = v + "(fg-white,bg-blue)"
		} else {
			queries[i] = v + "()"
		}
	}
}

func updateQueries(ls *ui.List) {
	markActiveQuery()
	ls.Items = queries
	ui.Render(ls)
}

func updateTickets(ls *ui.List) {
	markActiveTicket()
	ls.Items = tickets
	ui.Render(ls)
}

func nextPage() {
	if currentPage == ticketList {
		currentPage = ticketShow
	} else if currentPage == ticketShow {
		currentPage = ticketQuery
	} else if currentPage == ticketQuery {
		currentPage = ticketList
	}
	ui.StopLoop()
}

func handleTicketQueryPage() {

	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	ls := ui.NewList()
	ls.Items = queries
	ls.ItemFgColor = ui.ColorYellow
	ls.BorderLabel = "List"
	ls.Height = 10
	ls.Width = 80
	ls.Y = 0
	markActiveQuery()
	ui.Render(ls)

	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
		exitNow = true
	})
	ui.Handle("/sys/kbd/j", func(ui.Event) {
		nextQuery(1)
		updateQueries(ls)
	})
	ui.Handle("/sys/kbd/k", func(ui.Event) {
		prevQuery(1)
		updateQueries(ls)
	})
	ui.Handle("/sys/kbd/n", func(ui.Event) {
		nextPage()
	})

	ui.Loop()

}

func handleTicketListPage() {

	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	ls := ui.NewList()
	ls.Items = tickets
	ls.ItemFgColor = ui.ColorYellow
	ls.BorderLabel = "List"
	ls.Height = 10
	ls.Width = 80
	ls.Y = 0
	markActiveTicket()
	ui.Render(ls)

	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
		exitNow = true
	})
	ui.Handle("/sys/kbd/j", func(ui.Event) {
		nextTicket(1)
		updateTickets(ls)
	})
	ui.Handle("/sys/kbd/k", func(ui.Event) {
		prevTicket(1)
		updateTickets(ls)
	})
	ui.Handle("/sys/kbd/n", func(ui.Event) {
		nextPage()
	})
	ui.Loop()

}

func handleTicketShowPage() {

	fmt.Println("TicketShow!")
	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()
	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
		exitNow = true
	})
	ui.Handle("/sys/kbd/j", func(ui.Event) {
	})
	ui.Handle("/sys/kbd/k", func(ui.Event) {
	})
	ui.Handle("/sys/kbd/n", func(ui.Event) {
		nextPage()
	})
	ui.Loop()

}

func main() {

	for exitNow != true {

		switch {
		case currentPage == ticketQuery:
			handleTicketQueryPage()
		case currentPage == ticketList:
			handleTicketListPage()
		case currentPage == ticketShow:
			handleTicketShowPage()
		}

	}

}
