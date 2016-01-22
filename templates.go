package jiraui

const (
	default_list_template = `{{ range .issues }}[{{ .key | printf "%-12s"}}](fg-red)  [{{ if .fields.assignee }}{{ .fields.assignee.name | printf "%-10s" }}{{else}}{{"Unassigned"| printf "%-10s" }}{{end}} ](fg-blue) [{{ .fields.status.name | printf "%-12s"}}](fg-blue) [{{ dateFormat "2006-01-02" .fields.created }}](fg-blue)/[{{ dateFormat "2006-01-02T15:04" .fields.updated }}](fg-green)  {{ .fields.summary | printf "%-75s"}}
{{ end }}`
	default_view_template = `
issue: [{{ .key }}](fg-red)
summary: [{{ .fields.summary }}](fg-blue)

self: {{ .self }}
browse: ENDPOINT/browse/{{ .key }}
priority: {{ .fields.priority.name }}
status: {{ .fields.status.name }}
votes: {{ .fields.votes.votes }}
created: {{ .fields.created }}
updated: {{ .fields.updated }}
assignee: {{ if .fields.assignee }}{{ .fields.assignee.name }}{{end}}
reporter: {{ if .fields.reporter }}{{ .fields.reporter.name }}{{end}}
issuetype: {{ .fields.issuetype.name }}
{{if eq .fields.issuetype.name "Epic" }}epic_links: [<click here to show>](fg-red){{end}}
{{if .fields.customfield_10001 }}epic: [{{ .fields.customfield_10001 }}](fg-red){{end}}
{{if .fields.parent }}parent: [{{ .fields.parent.key }}](fg-red) -- {{ .fields.parent.fields.summary }}{{end}}
subtasks:
{{ range .fields.subtasks }}  - [{{ .key }}](fg-red)[{{.fields.status.name}}] -- {{.fields.summary}}
{{end}}

[labels:](fg-green){{ range .fields.labels }} {{ . }}{{end}}
[components:](fg-green){{ range .fields.components }} {{ .name }}{{end}}
[watchers:](fg-green){{ range .fields.customfield_10304 }} {{ .name }}{{end}}
[blockers:](fg-green)
{{ range .fields.issuelinks }}{{if .outwardIssue}}  - [{{ .outwardIssue.key }}](fg-red)[{{.outwardIssue.fields.status.name}}] -- {{.outwardIssue.fields.summary}}
{{end}}{{end}}
[depends:](fg-green)
{{ range .fields.issuelinks }}{{if .inwardIssue}}  - [{{ .inwardIssue.key }}](fg-red)[{{.inwardIssue.fields.status.name}}] -- {{.inwardIssue.fields.summary}}
{{end}}{{end}}

[description:](fg-green)

  {{ or .fields.description "" | indent 2 }}

[comments:](fg-green)

{{ range .fields.comment.comments }}
  - [{{.author.name}} at {{.created}}](fg-blue)
    {{ or .body "" | indent 4}}

{{end}}
`
	default_help_template = `
[Quick reference for jira-ui](fg-white)

[Actions:](fg-blue)

    <enter>      - select query/ticket
    L            - Label view (query results page only)
    E            - Edit ticket
    S            - Select sort order (query results page only)

[Commands (a'la vim/tig):](fg-blue)

    :comment {single-line-comment} - add a short comment to ticket
    :label {labels}                - add labels to selected ticket
    :label add/remove {labels}     - add/remove labels to selected ticket
    :take                          - assign ticket to self
    :assign {user}                 - assign ticket to {user}
    :unassign                      - unassign ticket
    :watch                         - watch ticket
    :<up>                          - select previous command
    :quit or :q                    - quit

[Navigation:](fg-blue)

    up/k         - previous line
    down/j       - next line
    C-f/<space>  - next page
    C-b          - previous page
    }            - next paragraph/section/fast-move
    {            - previous paragraph/section/fast-move
    n            - next search match
    g            - go to top of page
    G            - go to bottom of page
    q            - go back / quit
    C-c/Q        - quit

[Configuration:](fg-blue)

  It is very much recommended to read the go-jira documentation,
  particularly surrounding the .jira.d configuration directories.

  go-jira-ui uses this same mechanism, so can be used to load per-project
  defaults. It also leverages the templating engine, so you can customise
  the view of both the query output (use 'jira_ui_list' template), and the
  issue 'view' template.

  go-jira-ui reads its own [jira-ui-config.yml](fg-green) file in these
  jira.d directories, as not to pollute the go-jira config. You can add
  additional queries & sort orderings to the top-level Query page:

    $ cat ~/jira.d/jira-ui-config.yml:
    sorts:
      - name: "sort by vote count"
            jql:  "ORDER BY votes DESC"
    queries:
      - name: "alice assigned"
        jql:  "assignee = alice AND resolution = Unresolved"
      - name: "bob assigned"
        jql:  "assignee = bob AND resolution = Unresolved"
      - name: "unresolved must-do"
        jql:  "labels = 'must-do' AND resolution = Unresolved AND ( project = 'OPS' OR project = 'INFRA')"

  Learning JQL is highly recommended, the Atlassian Advanced Searching
  page is a good place to start.

`
)
