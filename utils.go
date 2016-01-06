package main

import (
	"bytes"
	"fmt"
	"github.com/Netflix-Skunkworks/go-jira"
	ui "github.com/gizak/termui"
	"gopkg.in/coryb/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

func countLabelsFromQuery(query string) map[string]int {
	counts := make(map[string]int)
	data, _ := runJiraQuery(query)
	issues := data.(map[string]interface{})["issues"].([]interface{})
	for _, issue := range issues {
		issueLabels := issue.(map[string]interface{})["fields"].(map[string]interface{})["labels"]
		for _, v := range issueLabels.([]interface{}) {
			label := v.(string)
			counts[label] = counts[label] + 1
		}
	}
	return counts
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

	loadConfigs(opts)
	for k, v := range defaults {
		if _, ok := opts[k]; !ok {
			log.Debug("Setting %q to %#v from defaults", k, v)
			opts[k] = v
		}
	}
	return opts
}

func lastLineDisplayed(ls *ui.List, firstLine int, correction int) int {
	return firstLine + ls.Height - correction
}
