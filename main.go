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

func lastLineDisplayed(ls *ui.List, firstLine int, correction int) int {
	return firstLine + ls.Height - correction
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
