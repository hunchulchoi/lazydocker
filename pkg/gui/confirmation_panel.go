// lots of this has been directly ported from one of the example files, will brush up later

// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"math"
	"strings"

	"github.com/fatih/color"
	"github.com/jesseduffield/gocui"
)

func (gui *Gui) wrappedConfirmationFunction(function func(*gocui.Gui, *gocui.View) error) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		if err := gui.closeConfirmationPrompt(); err != nil {
			return err
		}

		if function != nil {
			if err := function(g, v); err != nil {
				return err
			}
		}

		return nil
	}
}

func (gui *Gui) closeConfirmationPrompt() error {
	if err := gui.returnFocus(); err != nil {
		return err
	}
	gui.g.DeleteViewKeybindings("confirmation")
	gui.Views.Confirmation.Visible = false
	gui.State.confirmationIsInfo = false
	return nil
}

func (gui *Gui) getMessageHeight(wrap bool, message string, width int) int {
	lines := strings.Split(message, "\n")
	lineCount := 0
	// if we need to wrap, calculate height to fit content within view's width
	if wrap {
		for _, line := range lines {
			lineCount += len(line)/width + 1
		}
	} else {
		lineCount = len(lines)
	}
	return lineCount
}

func (gui *Gui) getConfirmationPanelDimensions(wrap bool, prompt string) (int, int, int, int) {
	if gui.State.confirmationIsInfo {
		return gui.getInfoPanelDimensions()
	}

	width, height := gui.g.Size()
	panelWidth := width / 2
	panelHeight := gui.getMessageHeight(wrap, prompt, panelWidth)
	return width/2 - panelWidth/2,
		height/2 - panelHeight/2 - panelHeight%2 - 1,
		width/2 + panelWidth/2,
		height/2 + panelHeight/2
}

func (gui *Gui) getInfoPanelDimensions() (int, int, int, int) {
	width, height := gui.g.Size()
	panelWidth := width * 4 / 5
	panelHeight := height * 7 / 10
	return (width - panelWidth) / 2,
		(height - panelHeight) / 2,
		(width + panelWidth) / 2,
		(height + panelHeight) / 2
}

func (gui *Gui) createPromptPanel(title string, handleConfirm func(*gocui.Gui, *gocui.View) error) error {
	gui.onNewPopupPanel()
	err := gui.prepareConfirmationPanel(title, "")
	if err != nil {
		return err
	}
	gui.Views.Confirmation.Editable = true
	return gui.setKeyBindings(gui.g, handleConfirm, nil)
}

func (gui *Gui) prepareConfirmationPanel(title, prompt string) error {
	x0, y0, x1, y1 := gui.getConfirmationPanelDimensions(true, prompt)
	confirmationView := gui.Views.Confirmation
	_, err := gui.g.SetView("confirmation", x0, y0, x1, y1, 0)
	if err != nil {
		return err
	}
	confirmationView.Title = title
	confirmationView.Visible = true
	gui.g.Update(func(g *gocui.Gui) error {
		return gui.switchFocus(confirmationView)
	})
	return nil
}

func (gui *Gui) onNewPopupPanel() {
	gui.Views.Menu.Visible = false
	gui.Views.PortsOverview.Visible = false
	gui.Views.Confirmation.Visible = false
}

// It is very important that within this function we never include the original prompt in any error messages, because it may contain e.g. a user password.
// The golangcilint unparam linter complains that handleClose is alwans nil but one day it won't be nil.
// nolint:unparam
func (gui *Gui) createConfirmationPanel(title, prompt string, handleConfirm, handleClose func(*gocui.Gui, *gocui.View) error) error {
	return gui.createPopupPanel(title, prompt, handleConfirm, handleClose, false)
}

func (gui *Gui) createPopupPanel(title, prompt string, handleConfirm, handleClose func(*gocui.Gui, *gocui.View) error, isInfo bool) error {
	gui.onNewPopupPanel()
	gui.g.Update(func(g *gocui.Gui) error {
		if gui.currentViewName() == "confirmation" {
			if err := gui.closeConfirmationPrompt(); err != nil {
				gui.Log.Error(err.Error())
			}
		}
		gui.State.confirmationIsInfo = isInfo
		err := gui.prepareConfirmationPanel(title, prompt)
		if err != nil {
			return err
		}
		gui.Views.Confirmation.Editable = false
		if err := gui.renderString(g, "confirmation", prompt); err != nil {
			return err
		}
		if gui.State.confirmationIsInfo {
			return gui.setInfoKeyBindings(g)
		}
		return gui.setKeyBindings(g, handleConfirm, handleClose)
	})
	return nil
}

func (gui *Gui) createInfoPanel(title, prompt string) error {
	return gui.createPopupPanel(title, prompt, nil, nil, true)
}

func (gui *Gui) setKeyBindings(g *gocui.Gui, handleConfirm, handleClose func(*gocui.Gui, *gocui.View) error) error {
	// would use a loop here but because the function takes an interface{} and slices of interfaces require even more boilerplate
	if err := g.SetKeybinding("confirmation", gocui.KeyEnter, gocui.ModNone, gui.wrappedConfirmationFunction(handleConfirm)); err != nil {
		return err
	}
	if err := g.SetKeybinding("confirmation", 'y', gocui.ModNone, gui.wrappedConfirmationFunction(handleConfirm)); err != nil {
		return err
	}

	if err := g.SetKeybinding("confirmation", gocui.KeyEsc, gocui.ModNone, gui.wrappedConfirmationFunction(handleClose)); err != nil {
		return err
	}
	if err := g.SetKeybinding("confirmation", 'n', gocui.ModNone, gui.wrappedConfirmationFunction(handleClose)); err != nil {
		return err
	}

	return nil
}

func (gui *Gui) setInfoKeyBindings(g *gocui.Gui) error {
	closeHandler := gui.wrappedConfirmationFunction(nil)
	scrollUp := func(g *gocui.Gui, v *gocui.View) error {
		return gui.scrollUpConfirmation()
	}
	scrollDown := func(g *gocui.Gui, v *gocui.View) error {
		return gui.scrollDownConfirmation()
	}

	bindings := []struct {
		key     interface{}
		mod     gocui.Modifier
		handler func(*gocui.Gui, *gocui.View) error
	}{
		{gocui.KeyEsc, gocui.ModNone, closeHandler},
		{gocui.KeyEnter, gocui.ModNone, closeHandler},
		{gocui.KeyPgup, gocui.ModNone, scrollUp},
		{gocui.KeyPgdn, gocui.ModNone, scrollDown},
		{gocui.KeyCtrlU, gocui.ModNone, scrollUp},
		{gocui.KeyCtrlD, gocui.ModNone, scrollDown},
		{'k', gocui.ModNone, scrollUp},
		{'j', gocui.ModNone, scrollDown},
	}

	for _, binding := range bindings {
		if err := g.SetKeybinding("confirmation", binding.key, binding.mod, binding.handler); err != nil {
			return err
		}
	}

	return nil
}

func (gui *Gui) scrollUpConfirmation() error {
	view := gui.Views.Confirmation
	ox, oy := view.Origin()
	newOy := int(math.Max(0, float64(oy-gui.Config.UserConfig.Gui.ScrollHeight)))
	return view.SetOrigin(ox, newOy)
}

func (gui *Gui) scrollDownConfirmation() error {
	view := gui.Views.Confirmation
	ox, oy := view.Origin()
	_, sizeY := view.Size()

	totalLines := view.ViewLinesHeight()
	if oy+sizeY >= totalLines {
		return nil
	}

	return view.SetOrigin(ox, oy+gui.Config.UserConfig.Gui.ScrollHeight)
}

func (gui *Gui) createErrorPanel(message string) error {
	colorFunction := color.New(color.FgRed).SprintFunc()
	coloredMessage := colorFunction(strings.TrimSpace(message))
	return gui.createConfirmationPanel(gui.Tr.ErrorTitle, coloredMessage, nil, nil)
}

func (gui *Gui) renderConfirmationOptions() error {
	if gui.State.confirmationIsInfo {
		return gui.renderInfoOptions()
	}

	optionsMap := map[string]string{
		"n/esc":   gui.Tr.No,
		"y/enter": gui.Tr.Yes,
	}
	return gui.renderOptionsMap(optionsMap)
}

func (gui *Gui) renderInfoOptions() error {
	optionsMap := map[string]string{
		"esc/enter": gui.Tr.Close,
		"↑ ↓":       gui.Tr.Navigate,
	}
	return gui.renderOptionsMap(optionsMap)
}
