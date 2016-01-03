go-jira-ui
----------

go-jira-ui is an ncurses command line tool for accessing JIRA.

It is built around the excellent [go-jira](https://github.com/Netflix-Skunkworks/go-jira) and
[termui](https://github.com/gizak/termui) libraries.

Currently it's focussed around browsing issues, supporting the following:

* Custom queries, with both a ticket and per-label view.
* Issue viewing
* Templatable list and issue views

In order to use this, you should configure an 'endpoint' as per the go-jira
documentation, essentially:

    $ cat ~/.jira.d/config.yml
    ---
    endpoint: https://jira.example.com/

This should be all that's needed to get going.

### Basic keys

    up/k         - previous line
    down/j       - next line
    C-f/<space>  - next page
    C-b          - previous page
    <enter>      - select item
    q            - go back / quit
    L            - Label view (query results page only)

### Configuration

It is very much recommended to read the
[go-jira](https://github.com/Netflix-Skunkworks/go-jira) documentation,
particularly surrounding the .jira.d configuration directories. go-jira-ui uses
this same mechanism, so can be used to load per-project defaults. It also
leverages the templating engine, so you can customise the view of both the
query output (use 'jira_ui_list' template), and the issue 'view' template.

go-jira-ui reads its own  `jira-ui-config.yml` file in these jira.d
directories, as not to pollute the go-jira config. You can add additional
queries to the top-level Query page:

    $ cat ~/jira.d/jira-ui-config.yml:
    queries:
      - name: "alice assigned"
        jql: "assignee = alice AND resolution = Unresolved"
      - name: "bob assigned"
        jql: "assignee = bob AND resolution = Unresolved"
      - name: "unresolved must-do"
        jql: "labels = 'must-do' AND resolution = Unresolved AND ( project = 'OPS' OR project = 'INFRA')"

Learning JQL is highly recommended, the Atlassian [Advanced
Searching](https://confluence.atlassian.com/jira/advanced-searching-179442050.html)
page is a good place to start.
