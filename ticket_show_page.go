package main

import (
	"fmt"
	ui "github.com/gizak/termui"
)

var ticketShowLineSelected = 0
var activeTicketShowList *ui.List

var currentTicketShowCache []string
var displayTicketShow []string

func prevTicketLine(n int) {
	ticketShowLineSelected = ticketShowLineSelected - n
	if ticketShowLineSelected < 0 {
		ticketShowLineSelected = 0
	}
}

func nextTicketLine(n int) {
	if ticketShowLineSelected < len(currentTicketShowCache)-n {
		ticketShowLineSelected = ticketShowLineSelected + n
	} else {
		ticketShowLineSelected = len(currentTicketShowCache) - 1
	}
}

func markActiveTicketLine() {
	for i, v := range currentTicketShowCache {
		selected := ""
		if i == ticketShowLineSelected {
			selected = "fg-white,bg-blue"
		}
		displayTicketShow[i] = fmt.Sprintf("[%s](%s)", v, selected)
	}
}

func updateTicketShowPage(ls *ui.List) {
	markActiveTicketLine()
	ui.Render(ls)
}

func handleTicketShowPage() {
	ui.Clear()
	ticketId := getTicketIdFromListLine(currentTicketListCache[ticketSelected])
	ticketShowLineSelected = 0
	currentTicketShowCache = JiraTicketAsStrings(ticketId)
	displayTicketShow = make([]string, len(currentTicketShowCache))
	ls := ui.NewList()
	ls.Items = displayTicketShow
	ls.ItemFgColor = ui.ColorYellow
	ls.Border = false
	ls.Height = ui.TermHeight()
	ls.Width = ui.TermWidth()
	ls.Overflow = "wrap"
	ls.Y = 0
	activeTicketShowList = ls
	markActiveTicketLine()
	ui.Render(ls)
}
