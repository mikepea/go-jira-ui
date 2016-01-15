package jiraui

import (
	ui "github.com/gizak/termui"
	"os"
)

func registerKeyboardHandlers() {
	ui.Handle("/sys/kbd/q", func(ui.Event) {
		handleBackKey()
	})
	ui.Handle("/sys/kbd/C-r", func(ui.Event) {
		handleRefreshKey()
	})
	ui.Handle("/sys/kbd/C-c", func(ui.Event) {
		ui.Close()
		os.Exit(0)
	})
	ui.Handle("/sys/kbd/Q", func(ui.Event) {
		ui.Close()
		os.Exit(0)
	})
	ui.Handle("/sys/kbd/}", func(ui.Event) {
		handleParaDownKey()
	})
	ui.Handle("/sys/kbd/{", func(ui.Event) {
		handleParaUpKey()
	})
	ui.Handle("/sys/kbd/j", func(ui.Event) {
		handleDownKey()
	})
	ui.Handle("/sys/kbd/k", func(ui.Event) {
		handleUpKey()
	})
	ui.Handle("/sys/kbd/<down>", func(ui.Event) {
		handleDownKey()
	})
	ui.Handle("/sys/kbd/<up>", func(ui.Event) {
		handleUpKey()
	})
	ui.Handle("/sys/kbd/g", func(ui.Event) {
		handleTopOfPageKey()
	})
	ui.Handle("/sys/kbd/G", func(ui.Event) {
		handleBottomOfPageKey()
	})
	ui.Handle("/sys/kbd/L", func(ui.Event) {
		handleLabelViewKey()
	})
	ui.Handle("/sys/kbd/<enter>", func(ui.Event) {
		handleSelectKey()
	})
	ui.Handle("/sys/kbd/<space>", func(ui.Event) {
		handlePageDownKey()
	})
	ui.Handle("/sys/kbd/C-f", func(ui.Event) {
		handlePageDownKey()
	})
	ui.Handle("/sys/kbd/C-b", func(ui.Event) {
		handlePageUpKey()
	})
	ui.Handle("/sys/kbd/E", func(ui.Event) {
		handleEditKey()
	})
	ui.Handle("/sys/kbd/C", func(ui.Event) {
		handleCommentKey()
	})
	ui.Handle("/sys/wnd/resize", func(ui.Event) {
		handleResize()
	})
}

func handleLabelViewKey() {
	switch page := currentPage.(type) {
	case *TicketListPage:
		q := new(LabelListPage)
		q.ActiveQuery = page.ActiveQuery
		currentPage = q
		changePage()
	}
	return
}

func handleRefreshKey() {
	if obj, ok := currentPage.(Refresher); ok {
		obj.Refresh()
	}
}

func handleEditKey() {
	if obj, ok := currentPage.(TicketEditer); ok {
		obj.EditTicket()
	}
}

func handleCommentKey() {
	if obj, ok := currentPage.(TicketCommenter); ok {
		obj.CommentTicket()
	}
}

func handleBackKey() {
	if obj, ok := currentPage.(GoBacker); ok {
		obj.GoBack()
	} else {
		ui.StopLoop()
		exitNow = true
	}
}

func handleResize() {
	changePage()
}

func handleSelectKey() {
	if obj, ok := currentPage.(ItemSelecter); ok {
		obj.SelectItem()
	}
}

func handleTopOfPageKey() {
	if obj, ok := currentPage.(PagePager); ok {
		obj.TopOfPage()
		obj.Update()
	}
}

func handleBottomOfPageKey() {
	if obj, ok := currentPage.(PagePager); ok {
		obj.BottomOfPage()
		obj.Update()
	}
}

func handleUpKey() {
	if obj, ok := currentPage.(PagePager); ok {
		obj.PreviousLine(1)
		obj.Update()
	}
}

func handleDownKey() {
	if obj, ok := currentPage.(PagePager); ok {
		obj.NextLine(1)
		obj.Update()
	}
}

func handlePageUpKey() {
	if obj, ok := currentPage.(PagePager); ok {
		obj.PreviousPage()
		obj.Update()
	}
}

func handlePageDownKey() {
	if obj, ok := currentPage.(PagePager); ok {
		obj.NextPage()
		obj.Update()
	}
}

func handleParaUpKey() {
	if obj, ok := currentPage.(PagePager); ok {
		obj.PreviousPara()
		obj.Update()
	}
}

func handleParaDownKey() {
	if obj, ok := currentPage.(PagePager); ok {
		obj.NextPara()
		obj.Update()
	}
}
