package main

import (
	"fmt"
	ui "github.com/gizak/termui"
	"strings"
)

var activeTicketListList *ui.List

var currentTicketListCache []string
var displayTickets []string

var ticketListLineSelected = 0
var displayTicketListFirstLine = 0

func prevTicket(n int) {
	ticketListLineSelected = ticketListLineSelected - n
	if ticketListLineSelected < 0 {
		ticketListLineSelected = 0
	}
	if ticketListLineSelected < displayTicketListFirstLine {
		displayTicketListFirstLine = ticketListLineSelected
	}
}

func nextTicket(n int) {
	if ticketListLineSelected < len(currentTicketListCache)-1 {
		ticketListLineSelected = ticketListLineSelected + n
	}
	if ticketListLineSelected > lastLineDisplayed(activeTicketListList, displayTicketListFirstLine, 3) {
		displayTicketListFirstLine = displayTicketListFirstLine + n
	}
}

func markActiveTicket() {
	for i, v := range currentTicketListCache {
		selected := ""
		if i == ticketListLineSelected {
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
	ls.Items = displayTickets[displayTicketListFirstLine:]
	ui.Render(ls)
}

func handleTicketListPage() {
	ui.Clear()
	ticketListLineSelected = 0
	queryName := origQueries[querySelected].Name
	queryJQL := origQueries[querySelected].JQL
	currentTicketListCache = displayQueryResults(queryJQL)
	displayTickets = make([]string, len(currentTicketListCache))
	ls := ui.NewList()
	ls.Items = displayTickets[displayTicketListFirstLine:]
	ls.ItemFgColor = ui.ColorYellow
	ls.BorderLabel = fmt.Sprintf("%s: %s", queryName, queryJQL)
	ls.Height = ui.TermHeight()
	ls.Width = ui.TermWidth()
	ls.Y = 0
	activeTicketListList = ls
	markActiveTicket()
	ui.Render(ls)
}
