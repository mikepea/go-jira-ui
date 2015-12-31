package main

import (
	"bytes"
	"fmt"
	ui "github.com/gizak/termui"
	"github.com/op/go-logging"
	"os"
	"strings"
	//jira "github.com/mikepea/go-jira"
	"jira"
)

const (
	ticketQuery = 1
	ticketList  = 2
	ticketShow  = 3
)

var exitNow = false

var currentPage = ticketQuery
var previousPage = ticketQuery

var currentQuery = ""
var previousQuery = ""

var ticketSelected = 0
var querySelected = 0
var ticketShowLineSelected = 0

var activeQueryList *ui.List
var activeTicketListList *ui.List
var activeTicketShowList *ui.List

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

func prevQuery(n int) {
	querySelected = querySelected - n
	if querySelected < 0 {
		querySelected = 0
	}
}

func nextQuery(n int) {
	if querySelected < len(origQueries)-1 {
		querySelected = querySelected + n
	}
}

type Query struct {
	Name string
	JQL  string
}

var origQueries = []Query{
	Query{"My Assigned Tickets", "assignee = currentUser() AND resolution = Unresolved"},
	Query{"My Reported Tickets", "reporter = currentUser() AND resolution = Unresolved"},
	Query{"My Watched Tickets", "watcher = currentUser() AND resolution = Unresolved"},
	Query{"OPS unlabelled", "project = OPS AND labels IS EMPTY AND resolution = Unresolved"},
	Query{"Ops Queue", "project = OPS AND resolution = Unresolved"},
}

var displayQueries []string

var currentTicketListCache []string
var displayTickets []string

var currentTicketShowCache []string
var displayTicketShow []string

func displayQueryResults(query string) []string {
	results := JiraQueryAsStrings(query)
	return results
}

func markActiveQuery() {
	for i, v := range origQueries {
		selected := ""
		if i == querySelected {
			selected = "fg-white,bg-blue"
		}
		displayQueries[i] = fmt.Sprintf("[%s](%s)", v.Name, selected)
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

func markActiveTicketLine() {
	for i, v := range currentTicketShowCache {
		selected := ""
		if i == ticketShowLineSelected {
			selected = "fg-white,bg-blue"
		}
		displayTicketShow[i] = fmt.Sprintf("[%s](%s)", v, selected)
	}
}

func handleTicketQueryPage() {
	ui.Clear()
	ls := ui.NewList()
	displayQueries = make([]string, len(origQueries))
	ls.Items = displayQueries
	ls.ItemFgColor = ui.ColorYellow
	ls.BorderLabel = "List"
	ls.Height = ui.TermHeight()
	ls.Width = ui.TermWidth()
	ls.Y = 0
	activeQueryList = ls
	markActiveQuery()
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

func getTicketIdFromListLine(line string) string {
	return strings.Split(line, " ")[0]
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

func changePage() {
	switch currentPage {
	case ticketQuery:
		handleTicketQueryPage()
	case ticketList:
		handleTicketListPage()
	case ticketShow:
		handleTicketShowPage()
	}
}

func getJiraOpts() map[string]interface{} {

	user := os.Getenv("USER")
	home := os.Getenv("HOME")
	defaultQueryFields := "summary,created,updated,priority,status,reporter,assignee,labels"
	defaultSort := "priority asc, created"
	defaultMaxResults := 500

	defaults := map[string]interface{}{
		"user":        user,
		"endpoint":    os.Getenv("JIRA_ENDPOINT"),
		"queryfields": defaultQueryFields,
		"directory":   fmt.Sprintf("%s/.jira.d/templates", home),
		"sort":        defaultSort,
		"max_results": defaultMaxResults,
		"method":      "GET",
		"quiet":       true,
	}
	//opts := make(map[string]interface{})
	return defaults
}

func runJiraQuery(query string) (interface{}, error) {
	opts := getJiraOpts()
	opts["query"] = query
	c := jira.New(opts)
	return c.FindIssues()
}

func JiraQueryAsStrings(query string) []string {
	opts := getJiraOpts()
	opts["query"] = query
	c := jira.New(opts)
	data, _ := c.FindIssues()
	buf := new(bytes.Buffer)
	jira.RunTemplate(c.GetTemplate("list"), data, buf)
	return strings.Split(strings.TrimSpace(buf.String()), "\n")
}

func JiraTicketAsStrings(id string) []string {
	opts := getJiraOpts()
	c := jira.New(opts)
	data, _ := c.ViewIssue(id)
	buf := new(bytes.Buffer)
	jira.RunTemplate(c.GetTemplate("view"), data, buf)
	return strings.Split(strings.TrimSpace(buf.String()), "\n")
}

func updateQueryPage(ls *ui.List) {
	markActiveQuery()
	ls.Items = displayQueries
	ui.Render(ls)
}

func updateTicketListPage(ls *ui.List) {
	markActiveTicket()
	ui.Render(ls)
}

func updateTicketShowPage(ls *ui.List) {
	markActiveTicketLine()
	ui.Render(ls)
}

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

var (
	log    = logging.MustGetLogger("jira")
	format = "%{color}%{time:2006-01-02T15:04:05.000Z07:00} %{level:-5s} [%{shortfile}]%{color:reset} %{message}"
)

func main() {

	opts := getJiraOpts()

	logging.SetLevel(logging.NOTICE, "")

	c := jira.New(opts)

	// TODO: make this as quick as can be
	if _, err := runJiraQuery("assignee = CurrentUser() AND resolution = Unresolved"); err != nil {
		c.CmdLogin()
	}

	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	registerKeyboardHandlers()

	for exitNow != true {

		handleTicketQueryPage()
		ui.Loop()

	}

}
