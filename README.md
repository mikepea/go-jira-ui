go-jira-ui
----------

go-jira-ui is an ncurses command line tool for accessing JIRA.

It is built around the excellent [go-jira](https://github.com/Netflix-Skunkworks/go-jira) and
[termui](https://github.com/gizak/termui) libraries.

In order to use this, you should configure an 'endpoint' as per the go-jira
documentation:

    $ cat ~/.jira.d/config.yml
    ---
    endpoint: https://jira.example.com/
    user: bob   # if not same as $USER

This should be all that's needed to get going.

### Features

* Supply your own JQL queries to view
* Label view of a given query, to see categorisations easily
* View tickets from the query
* Edit/Comment on tickets from both list and detail view
* Drill into sub/blocker/related/mentioned tickets in details view
* Basic compatibility with [go-jira](https://github.com/Netflix-Skunkworks/go-jira) commandline and options loading

At present, edit and comment will exit after the update. This is a workaround
to an implementation issue, being tracked in [#8](https://github.com/mikepea/go-jira-ui/issues/8)

### Usage

`jira-ui` is intended to mirror the options of go-jira's `jira` tool, where
useful:

    jira-ui             # opens up in Query List page. Default interface.
    jira-ui ISSUE       # opens up Ticket Show page, with ISSUE loaded
    jira-ui ls -q JQL   # opens up Ticket List page, with results of JQL loaded.

### Basic keys

Actions:

    <enter>      - select query/ticket
    L            - Label view (query results page only)
    E            - Edit ticket
    C            - Comment on ticket

Navigation:

    up/k         - previous line
    down/j       - next line
    C-f/<space>  - next page
    C-b          - previous page
    g            - go to top of page
    G            - go to bottom of page
    q            - go back / quit
    C-c/Q        - quit


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
