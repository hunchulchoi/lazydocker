package presentation

import (
	"github.com/hunchulchoi/lazydocker/pkg/commands"
	"github.com/hunchulchoi/lazydocker/pkg/utils"
)

func GetImageDisplayStrings(image *commands.Image) []string {
	return displayStringsMutedIf([]string{
		image.Name,
		image.Tag,
		utils.FormatDecimalBytes(int(image.Image.Size)),
	}, image.IsDangling())
}
