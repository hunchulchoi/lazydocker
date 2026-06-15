package presentation

import (
	"github.com/fatih/color"
	"github.com/hunchulchoi/lazydocker/pkg/utils"
	"github.com/samber/lo"
)

func displayStringsMutedIf(strs []string, muted bool) []string {
	if !muted {
		return strs
	}

	return lo.Map(strs, func(str string, _ int) string {
		return utils.ColoredString(str, color.FgHiBlack)
	})
}
