# Developing go-jira-ui

This is a notes doc mostly for myself when I pick up and put down this project,
but hopefully will also be useful for anyone wishing to contribute.


## Debugging

Delve is real useful, but needs to be run headless otherwise Bad Things Happen
with your terminal.

     dlv --headless --listen 127.0.0.1:5533 debug ./jira-ui/main.go -- BLAH-1

     dlv connect 127.0.0.1:5533
