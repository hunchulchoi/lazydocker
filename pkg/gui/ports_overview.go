package gui

import (
	"github.com/jesseduffield/gocui"
	"github.com/hunchulchoi/lazydocker/pkg/commands"
)

func (gui *Gui) handleShowAllPorts(g *gocui.Gui, v *gocui.View) error {
	if gui.isPopupPanel(v.Name()) {
		return nil
	}

	return gui.WithWaitingStatus(gui.Tr.LoadingPorts, func() error {
		containers, err := gui.DockerCommand.GetContainers(nil)
		if err != nil {
			return gui.createErrorPanel(err.Error())
		}

		if err := gui.DockerCommand.LoadAllContainerDetails(containers); err != nil {
			gui.Log.Error(err)
		}

		rows := commands.CollectAllPorts(containers)
		if len(rows) == 0 {
			return gui.createErrorPanel(gui.Tr.NothingToDisplay)
		}

		return gui.showPortsOverviewPanel(rows)
	})
}
