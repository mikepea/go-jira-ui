package jiraui

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/Netflix-Skunkworks/go-jira"
	ui "github.com/gizak/termui"
	"github.com/mitchellh/go-wordwrap"
	"gopkg.in/coryb/yaml.v2"
)

func countLabelsFromQuery(query string) map[string]int {
	data, _ := runJiraQuery(query)
	return countLabelsFromQueryData(data)
}

func countLabelsFromQueryData(data interface{}) map[string]int {
	counts := make(map[string]int)
	issues := data.(map[string]interface{})["issues"].([]interface{})
	for _, issue := range issues {
		issueLabels := issue.(map[string]interface{})["fields"].(map[string]interface{})["labels"]
		labels := issueLabels.([]interface{})
		if len(labels) == 0 {
			// "NOT LABELLED" isn't a valid label, so no possible conflict here.
			counts["NOT LABELLED"] = counts["NOT LABELLED"] + 1
		} else {
			for _, v := range labels {
				label := v.(string)
				counts[label] = counts[label] + 1
			}
		}
	}
	return counts
}

func RunExternalCommand(fn func() error) error {
	log.Debugf("ShellOut() called with %q", fn)
	deregisterEventHandlers()
	ui.Clear()
	stty := exec.Command("stty", "-f", "/dev/tty", "echo", "opost")
	_ = stty.Run()
	err := fn() // magic happens
	stty = exec.Command("stty", "-f", "/dev/tty", "-echo", "-opost")
	_ = stty.Run()
	registerEventHandlers()
	if err != nil {
		return err
	}
	return nil
}

func runShell() {
	_ = RunExternalCommand(
		func() error {
			cmd := exec.Command("bash")
			cmd.Stdout, cmd.Stderr, cmd.Stdin = os.Stdout, os.Stderr, os.Stdin
			return cmd.Run()
		})
}

func runJiraCmdEdit(ticketId string) {
	_ = RunExternalCommand(
		func() error {
			opts := getJiraOpts()
			c := jira.New(opts)
			return c.CmdEdit(ticketId)
		})
	switch c := currentPage.(type) {
	case Refresher:
		c.Refresh()
	}
	changePage()
}

func runJiraCmdCommentNoEditor(ticketId string, comment string) {
	opts := getJiraOpts()
	opts["comment"] = comment
	c := jira.New(opts)
	c.CmdComment(ticketId)
}

func runJiraCmdAssign(ticketId string, user string) {
	opts := getJiraOpts()
	c := jira.New(opts)
	c.CmdAssign(ticketId, user)
}

func runJiraCmdWatch(ticketId string, watcher string, remove bool) {
	opts := getJiraOpts()
	c := jira.New(opts)
	if watcher == "" {
		watcher = opts["user"].(string)
	}
	c.CmdWatch(ticketId, watcher, remove)
}

func runJiraCmdVote(ticketId string, up bool) {
	opts := getJiraOpts()
	c := jira.New(opts)
	c.CmdVote(ticketId, up)
}

func runJiraCmdLabels(ticketId string, action string, labels []string) {
	opts := getJiraOpts()
	c := jira.New(opts)
	err := c.CmdLabels(action, ticketId, labels)
	if err != nil {
		log.Errorf("Error writing labels: %q", err)
	}
}

func findTicketIdInString(line string) string {
	re := regexp.MustCompile(`[A-Z]{2,12}-[0-9]{1,6}`)
	return strings.TrimSpace(re.FindString(line))
}

func runJiraQuery(query string) (interface{}, error) {
	opts := getJiraOpts()
	opts["query"] = query
	c := jira.New(opts)
	return c.FindIssues()
}

func JiraQueryAsStrings(query string, templateName string) []string {
	opts := getJiraOpts()
	opts["query"] = query
	c := jira.New(opts)
	data, _ := c.FindIssues()
	buf := new(bytes.Buffer)
	if templateName == "" {
		templateName = "jira_ui_list"
	}
	template := c.GetTemplate(templateName)
	if template == "" {
		template = default_list_template
	}
	jira.RunTemplate(template, data, buf)
	return strings.Split(strings.TrimSpace(buf.String()), "\n")
}

func FetchJiraTicket(id string) (interface{}, error) {
	opts := getJiraOpts()
	c := jira.New(opts)
	return c.ViewIssue(id)
}

func JiraTicketAsStrings(data interface{}, templateName string) []string {
	opts := getJiraOpts()
	c := jira.New(opts)
	buf := new(bytes.Buffer)
	template := c.GetTemplate(templateName)
	log.Debugf("JiraTicketsAsStrings: template = %q", template)
	if template == "" {
		template = strings.Replace(default_view_template, "ENDPOINT", opts["endpoint"].(string), 1)
	}
	jira.RunTemplate(template, data, buf)
	return strings.Split(strings.TrimSpace(buf.String()), "\n")
}

func HelpTextAsStrings(data interface{}, templateName string) []string {
	opts := getJiraOpts()
	c := jira.New(opts)
	buf := new(bytes.Buffer)
	template := c.GetTemplate(templateName)
	if template == "" {
		template = default_help_template
	}
	log.Debugf("HelpTextAsStrings: template = %q", template)
	jira.RunTemplate(template, data, buf)
	return strings.Split(strings.TrimSpace(buf.String()), "\n")
}

func WrapText(lines []string, maxWidth uint) []string {
	out := make([]string, 0)
	insideNoformatBlock := false
	insideCodeBlock := false
	for _, line := range lines {
		if matched, _ := regexp.MatchString(`^\s+\{code`, line); matched {
			insideCodeBlock = !insideCodeBlock
		} else if strings.TrimSpace(line) == "{noformat}" {
			insideNoformatBlock = !insideNoformatBlock
		}
		if maxWidth == 0 || uint(len(line)) < maxWidth || insideCodeBlock || insideNoformatBlock {
			out = append(out, line)
			continue
		}
		if matched, _ := regexp.MatchString(`^[a-z_]+:\s`, line); matched {
			// don't futz with single line field+value.
			// If they are too long, that's their fault.
			out = append(out, line)
			continue
		}
		// wrap text, but preserve indenting
		re := regexp.MustCompile(`^\s*`)
		indenting := re.FindString(line)
		wrappedLines := strings.Split(wordwrap.WrapString(line, maxWidth-uint(len(indenting))), "\n")
		indentedWrappedLines := make([]string, len(wrappedLines))
		for i, wl := range wrappedLines {
			if i == 0 {
				// first line already has the indent
				indentedWrappedLines[i] = wl
			} else {
				indentedWrappedLines[i] = indenting + wl
			}
		}
		out = append(out, indentedWrappedLines...)
	}
	return out
}

func parseYaml(file string, v map[string]interface{}) {
	if fh, err := ioutil.ReadFile(file); err == nil {
		log.Debugf("Parsing YAML file: %s", file)
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
					log.Debugf("Setting %q to %#v from %s", k, v, file)
					opts[k] = v
				}
			}
		}
	}
}

func doLogin(opts map[string]interface{}) error {
	c := jira.New(opts)
	fmt.Printf("Logging in as %s:\n", opts["user"])
	return c.CmdLogin()
}

func ensureLoggedIntoJira() error {
	homeDir := os.Getenv("HOME")
	opts := getJiraOpts()
	testSessionQuery := fmt.Sprintf("reporter = %s", opts["user"])
	if _, err := os.Stat(fmt.Sprintf("%s/.jira.d/cookies.js", homeDir)); err != nil {
		return doLogin(opts)
	} else if data, err := runJiraQuery(testSessionQuery); err != nil {
		return doLogin(opts)
	} else if val, ok := data.(map[string]interface{})["errorMessages"]; ok {
		if len(val.([]interface{})) > 0 {
			return doLogin(opts)
		}
	}
	return nil
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

	for k, v := range cliOpts {
		if _, ok := opts[k]; !ok {
			log.Debugf("Setting %q to %#v from cli options", k, v)
			opts[k] = v
		}
	}

	loadConfigs(opts)
	for k, v := range defaults {
		if _, ok := opts[k]; !ok {
			log.Debugf("Setting %q to %#v from defaults", k, v)
			opts[k] = v
		}
	}
	return opts
}

func lastLineDisplayed(ls *ScrollableList, firstLine int, correction int) int {
	return firstLine + ls.Height - correction
}
