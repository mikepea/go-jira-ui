package main

import (
	"bytes"
	"fmt"
	"github.com/Netflix-Skunkworks/go-jira"
	ui "github.com/gizak/termui"
	"github.com/op/go-logging"
	"gopkg.in/coryb/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

const (
	ticketQuery = iota
	ticketList  = iota
	labelList   = iota
	ticketShow  = iota
)

const (
	default_list_template = `{{ range .issues }}{{ .key | printf "%-20s"}}  {{ dateFormat "2006-01-02" .fields.created }}/{{ dateFormat "2006-01-02T15:04" .fields.updated }}  {{ .fields.summary | printf "%-75s"}} -- labels({{ join "," .fields.labels }})
{{ end }}`
)

var exitNow = false

var currentPage = ticketQuery
var previousPage = ticketQuery

var ticketQueryPage QueryPage
var ticketListPage TicketListPage
var ticketShowPage TicketShowPage

func changePage() {
	switch currentPage {
	case ticketQuery:
		ticketQueryPage.Create()
	case ticketList:
		ticketListPage.Create()
	case ticketShow:
		ticketShowPage.Create()
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
	defaultMaxResults := 1000

	opts := make(map[string]interface{})
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

	loadConfigs(opts)
	for k, v := range defaults {
		if _, ok := opts[k]; !ok {
			log.Debug("Setting %q to %#v from defaults", k, v)
			opts[k] = v
		}
	}
	return opts
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
	// TODO: this is a nasty hack, make it less so
	// template must start {key} and (for labels view) end with '-- labels()'
	template := c.GetTemplate("jira_ui_list")
	if template == "" {
		template = default_list_template
	}
	jira.RunTemplate(template, data, buf)
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

func parseYaml(file string, v map[string]interface{}) {
	if fh, err := ioutil.ReadFile(file); err == nil {
		log.Debug("Parsing YAML file: %s", file)
		yaml.Unmarshal(fh, &v)
	}
}

func loadConfigs(opts map[string]interface{}) {
	paths := jira.FindParentPaths(".jira.d/jira-ui-config.yml")
	paths = append(jira.FindParentPaths(".jira.d/config.yml"), paths...)
	paths = append([]string{"/etc/go-jira-ui.yml", "/etc/go-jira.yml"}, paths...)

	// iterate paths in reverse
	for i := len(paths) - 1; i >= 0; i-- {
		file := paths[i]
		if _, err := os.Stat(file); err == nil {
			tmp := make(map[string]interface{})
			parseYaml(file, tmp)
			for k, v := range tmp {
				if _, ok := opts[k]; !ok {
					log.Debug("Setting %q to %#v from %s", k, v, file)
					opts[k] = v
				}
			}
		}
	}
}

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

		ticketQueryPage.Create()
		ui.Loop()

	}

}
