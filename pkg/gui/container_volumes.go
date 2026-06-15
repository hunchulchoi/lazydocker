package gui

import (
	"github.com/hunchulchoi/lazydocker/pkg/commands"
	"github.com/hunchulchoi/lazydocker/pkg/tasks"
)

func (gui *Gui) renderContainerVolumes(container *commands.Container) tasks.TaskFunc {
	return gui.NewSimpleRenderStringTask(func() string {
		sizes := gui.containerMountSizes()
		output, err := container.RenderVolumes(gui.Tr, sizes)
		if err != nil {
			gui.Log.Error(err)
			return err.Error()
		}
		return output
	})
}

func (gui *Gui) renderServiceContainerVolumes(service *commands.Service) tasks.TaskFunc {
	if service.Container == nil {
		return gui.NewSimpleRenderStringTask(func() string { return gui.Tr.NoContainer })
	}

	return gui.renderContainerVolumes(service.Container)
}

func (gui *Gui) containerMountSizes() commands.MountSizes {
	volumeSizes, err := gui.DockerCommand.VolumeSizesByName()
	if err != nil {
		gui.Log.Error(err)
	}

	return commands.MountSizes{
		VolumeSizes: volumeSizes,
		OSCommand:   gui.OSCommand,
	}
}
