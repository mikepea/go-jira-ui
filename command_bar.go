package jiraui

import (
	ui "github.com/gizak/termui"
)

type CommandBar struct {
	uiList *ui.List
	EditBox
	previousCommands []string
}

func (p *CommandBar) Submit() {
	if obj, ok := currentPage.(CommandBoxer); ok {
		obj.SetCommandMode(false)
		obj.ExecuteCommand()
		p.previousCommands = append(p.previousCommands, string(p.text))
	}
	// currentPage may have changed
	if obj, ok := currentPage.(CommandBoxer); ok {
		obj.Update()
	}
}

func (p *CommandBar) PreviousCommand() {
	if len(p.previousCommands) == 0 {
		return
	}
	if obj, ok := currentPage.(CommandBoxer); ok {
		p.text = []byte(p.previousCommands[len(p.previousCommands)-1])
		obj.Update()
	}
}

func (p *CommandBar) Reset() {
	p.text = []byte(``)
	p.line_voffset = 0
	p.cursor_boffset = 0
	p.cursor_voffset = 0
	p.cursor_coffset = 0
}

func (p *CommandBar) Update() {
	ls := p.uiList
	strs := []string{string(p.text)}
	ls.Items = strs
	ui.Render(ls)
}

func (p *CommandBar) Create() {
	ls := ui.NewList()
	p.uiList = ls
	ls.ItemFgColor = ui.ColorGreen
	ls.Border = false
	ls.Height = 1
	ls.Width = ui.TermWidth()
	ls.X = 0
	ls.Y = ui.TermHeight() - 1
	p.Update()
}
