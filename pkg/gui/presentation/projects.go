package presentation

import "github.com/hunchulchoi/lazydocker/pkg/commands"

func GetProjectDisplayStrings(project *commands.Project) []string {
	return []string{project.Name}
}
