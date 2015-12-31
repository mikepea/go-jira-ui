package main

import (
	"fmt"
	ui "github.com/gizak/termui"
	"strings"
)

var activeTicketListList *ui.List

var currentTicketListCache []string
var displayTickets []string

var ticketSelected = 0

func prevTicket(n int) {
	ticketSelected = ticketSelected - n
	if ticketSelected < 0 {
		ticketSelected = 0
	}
}

func nextTicket(n int) {
	if ticketSelected < len(currentTicketListCache)-1 {
		ticketSelected = ticketSelected + n
	}
}

func markActiveTicket() {
	for i, v := range currentTicketListCache {
		selected := ""
		if i == ticketSelected {
			selected = "fg-white,bg-blue"
		}
		displayTickets[i] = fmt.Sprintf("[%s](%s)", v, selected)
	}
}

func getTicketIdFromListLine(line string) string {
	return strings.Split(line, " ")[0]
}

func updateTicketListPage(ls *ui.List) {
	markActiveTicket()
	ui.Render(ls)
}

func handleTicketListPage() {
	ui.Clear()
	ticketSelected = 0
	queryName := origQueries[querySelected].Name
	queryJQL := origQueries[querySelected].JQL
	currentTicketListCache = displayQueryResults(queryJQL)
	displayTickets = make([]string, len(currentTicketListCache))
	ls := ui.NewList()
	ls.Items = displayTickets
	ls.ItemFgColor = ui.ColorYellow
	ls.BorderLabel = fmt.Sprintf("%s: %s", queryName, queryJQL)
	ls.Height = ui.TermHeight()
	ls.Width = ui.TermWidth()
	ls.Y = 0
	activeTicketListList = ls
	markActiveTicket()
	ui.Render(ls)
}
