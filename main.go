package main

import (
	"bytes"
	"fmt"
	ui "github.com/gizak/termui"
	"github.com/op/go-logging"
	"os"
	"strings"
	"time"
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

type Query struct {
	Name string
	JQL  string
}

var origQueries = []Query{
	Query{"My Tickets", "--project OPS AND owner = 'mikepea'"},
	Query{"My Watched Tickets", "--project OPS AND watcher = 'mikepea'"},
	Query{"unlabelled", "--project OPS AND labels IS EMPTY"},
	Query{"Ops Queue", "--project OPS"},
}

var queries = []string{
	"My Tickets",
	"My Watched Tickets",
	"unlabelled",
	"OPS queue",
}

func markActiveQuery() {
	for i, v := range origQueries {
		selected := ""
		if i == querySelected {
			selected = "fg-white,bg-blue"
		}
		queries[i] = fmt.Sprintf("[%s](%s)", v.Name, selected)
	}
}

func updateQueries(ls *ui.List) {
	markActiveQuery()
	ls.Items = queries
	ui.Render(ls)
}

/*
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
*/

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

	ui.Loop()

}

func getJiraOpts() map[string]interface{} {

	user := os.Getenv("USER")
	home := os.Getenv("HOME")
	defaultQueryFields := "summary,created,updated,priority,status,reporter,assignee"
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
		"quiet":       false,
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
	return strings.Split(buf.String(), "\n")
}

var (
	log    = logging.MustGetLogger("jira")
	format = "%{color}%{time:2006-01-02T15:04:05.000Z07:00} %{level:-5s} [%{shortfile}]%{color:reset} %{message}"
)

func main() {

	opts := getJiraOpts()

	logging.SetLevel(logging.NOTICE, "")

	c := jira.New(opts)

	// check to see if we can run a query, otherwise force a login
	// TODO: make this quicker somehow
	if _, err := runJiraQuery("assignee = 'mikepea' AND resolution = Unresolved"); err != nil {
		//fmt.Println(err)
		c.CmdLogin()
	}

	lines := JiraQueryAsStrings("assignee = 'mikepea' AND resolution = Unresolved")
	fmt.Println(lines)

	// debug pause
	//time.Sleep(2 * time.Millisecond)
	time.Sleep(5 * time.Second)

	for exitNow != true {

		handleTicketQueryPage()

	}

}
