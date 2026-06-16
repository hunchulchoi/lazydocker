package gui

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/jesseduffield/gocui"
	"github.com/hunchulchoi/lazydocker/pkg/commands"
	"github.com/hunchulchoi/lazydocker/pkg/gui/panels"
	"github.com/hunchulchoi/lazydocker/pkg/utils"
	"github.com/samber/lo"
)

type PortsOverviewPanel struct {
	*panels.SideListPanel[commands.PortOverviewRow]
	gui *Gui
}

const portsOverviewHeaderLines = 1

func (panel *PortsOverviewPanel) Refocus() {
	panel.gui.FocusY(
		panel.SelectedIdx+portsOverviewHeaderLines,
		panel.List.Len()+portsOverviewHeaderLines,
		panel.View,
	)
}

func (panel *PortsOverviewPanel) HandleClick() error {
	itemCount := panel.List.Len()
	selectedLine := &panel.SelectedIdx

	if err := panel.gui.handleClickAux(
		panel.View,
		itemCount+portsOverviewHeaderLines,
		selectedLine,
		func(g *gocui.Gui, v *gocui.View) error {
			if *selectedLine < portsOverviewHeaderLines {
				*selectedLine = 0
			} else {
				*selectedLine -= portsOverviewHeaderLines
			}
			return panel.HandleSelect()
		},
	); err != nil {
		return err
	}

	return nil
}

func (panel *PortsOverviewPanel) HandleSelect() error {
	_, err := panel.GetSelectedItem()
	if err != nil {
		if err.Error() != panel.NoItemsMessage {
			return err
		}
		if panel.NoItemsMessage != "" {
			panel.gui.NewSimpleRenderStringTask(func() string { return panel.NoItemsMessage })
		}
		return nil
	}

	panel.Refocus()
	return nil
}

func (panel *PortsOverviewPanel) HandleNextLine() error {
	panel.SelectNextLine()
	return panel.HandleSelect()
}

func (panel *PortsOverviewPanel) HandlePrevLine() error {
	panel.SelectPrevLine()
	return panel.HandleSelect()
}

func (panel *PortsOverviewPanel) RerenderList() error {
	return panel.gui.rerenderPortsOverviewList()
}

func (gui *Gui) getPortsOverviewPanel() *PortsOverviewPanel {
	sideListPanel := &panels.SideListPanel[commands.PortOverviewRow]{
		ListPanel: panels.ListPanel[commands.PortOverviewRow]{
			List: panels.NewFilteredList[commands.PortOverviewRow](),
			View: gui.Views.PortsOverview,
		},
		NoItemsMessage: gui.Tr.NothingToDisplay,
		Gui:            gui.intoInterface(),
		Sort: func(a, b commands.PortOverviewRow) bool {
			less := commands.ComparePortOverviewRows(a, b, gui.State.portsOverviewSortColumn)
			if gui.State.portsOverviewSortAsc {
				return less
			}
			return !less
		},
		GetTableCells: commands.PortOverviewDisplayCells,
		OnRerender: func() error {
			return gui.resizePopupPanel(gui.Views.PortsOverview)
		},
	}

	return &PortsOverviewPanel{
		SideListPanel: sideListPanel,
		gui:           gui,
	}
}

func (gui *Gui) showPortsOverviewPanel(rows []commands.PortOverviewRow) error {
	gui.onNewPopupPanel()
	gui.State.portsOverviewSortColumn = commands.PortSortExternal
	gui.State.portsOverviewSortAsc = true

	if err := gui.preparePortsOverviewPanel(); err != nil {
		return err
	}

	gui.Panels.PortsOverview.SetItems(rows)
	if err := gui.Panels.PortsOverview.RerenderList(); err != nil {
		return err
	}

	gui.Views.PortsOverview.Title = gui.Tr.AllPortsTitle
	gui.Views.PortsOverview.Visible = true

	return gui.switchFocus(gui.Views.PortsOverview)
}

func (gui *Gui) preparePortsOverviewPanel() error {
	x0, y0, x1, y1 := gui.getInfoPanelDimensions()
	view := gui.Views.PortsOverview
	_, err := gui.g.SetView("portsOverview", x0, y0, x1, y1, 0)
	if err != nil {
		return err
	}
	view.Highlight = true
	return nil
}

func (gui *Gui) rerenderPortsOverviewList() error {
	panel := gui.Panels.PortsOverview.SideListPanel
	panel.FilterAndSort()

	headers := commands.PortOverviewHeaders(
		commands.PortOverviewHeaderLabels(gui.Tr),
		gui.State.portsOverviewSortColumn,
		gui.State.portsOverviewSortAsc,
	)
	headerTable, err := utils.RenderTable([][]string{headers})
	if err != nil {
		return err
	}

	bodyRows := lo.Map(panel.List.GetItems(), func(row commands.PortOverviewRow, _ int) []string {
		return portOverviewDisplayRow(row)
	})

	bodyTable := ""
	if len(bodyRows) == 0 {
		bodyTable = gui.Tr.NothingToDisplay
	} else {
		bodyTable, err = utils.RenderTable(bodyRows)
		if err != nil {
			return err
		}
	}

	gui.Update(func() error {
		panel.View.Clear()
		fmt.Fprint(panel.View, headerTable+"\n"+bodyTable)
		panel.Refocus()
		return nil
	})

	return nil
}

func portOverviewDisplayRow(row commands.PortOverviewRow) []string {
	cells := commands.PortOverviewDisplayCells(row)
	if !row.Conflict {
		return cells
	}

	for i, cell := range cells {
		cells[i] = utils.ColoredString(cell, color.FgRed)
	}

	return cells
}

func (gui *Gui) handlePortsOverviewClose() error {
	gui.Views.PortsOverview.Visible = false

	if gui.State.Filter.panel == gui.Panels.PortsOverview {
		if err := gui.clearFilter(); err != nil {
			return err
		}
		gui.removeViewFromStack(gui.Views.Filter)
	}

	return gui.returnFocus()
}

func (gui *Gui) handlePortsOverviewSort(column commands.PortOverviewSortColumn) error {
	if gui.State.portsOverviewSortColumn == column {
		gui.State.portsOverviewSortAsc = !gui.State.portsOverviewSortAsc
	} else {
		gui.State.portsOverviewSortColumn = column
		gui.State.portsOverviewSortAsc = true
	}

	return gui.Panels.PortsOverview.RerenderList()
}

func (gui *Gui) renderPortsOverviewOptions() error {
	optionsMap := map[string]string{
		"esc/q": gui.Tr.Close,
		"/":     gui.Tr.LcFilter,
		"↑ ↓":   gui.Tr.Navigate,
		"1-7":   gui.Tr.PortsSortByColumn,
	}
	return gui.renderOptionsMap(optionsMap)
}
