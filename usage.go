package jiraui

const usage_format = `
Usage:
  jira-ui ls <Query Options>
  jira-ui ISSUE
  jira-ui

General Options:
  -e --endpoint=URI   URI to use for jira
  -l --log=FILE       FILE to use for log (default /dev/null)
  -h --help           Show this usage
  -u --user=USER      Username to use for authenticaion
  -v --verbose        Increase output logging
  --version           Print version

Ticket View Options:
  -t --template=FILE  Template file to use for viewing tickets
  -m --max_wrap=VAL   Maximum word-wrap width when viewing ticket text (0 disables)

Query Options:
  -q --query=JQL            Jira Query Language expression for the search
  -f --queryfields=FIELDS   Fields that are used in "list" view

`
