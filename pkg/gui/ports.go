package gui

import (
	"github.com/jesseduffield/lazydocker/pkg/commands"
	"github.com/jesseduffield/lazydocker/pkg/tasks"
)

func (gui *Gui) renderContainerPorts(container *commands.Container) tasks.TaskFunc {
	return gui.NewSimpleRenderStringTask(func() string {
		output, err := container.RenderPorts(gui.Tr)
		if err != nil {
			gui.Log.Error(err)
			return err.Error()
		}
		return output
	})
}

func (gui *Gui) renderServiceContainerPorts(service *commands.Service) tasks.TaskFunc {
	if service.Container == nil {
		return gui.NewSimpleRenderStringTask(func() string { return gui.Tr.NoContainer })
	}

	return gui.renderContainerPorts(service.Container)
}
