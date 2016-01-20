package jiraui

import (
	ui "github.com/gizak/termui"
)

type CommandBar struct {
	uiList *ui.List
	EditBox
}

func (p *CommandBar) Submit() {
	if obj, ok := currentPage.(CommandBoxer); ok {
		obj.SetCommandMode(false)
		obj.ExecuteCommand()
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
