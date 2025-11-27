package history

import (
	"path/filepath"
	"sort"
	"strings"
)

// ApplyContextBoost updates commands with context-aware scores and sorts them.
func ApplyContextBoost(commands []Command, current CommandContext) []Command {
	if len(commands) == 0 {
		return commands
	}

	for i := range commands {
		commands[i].ContextBoost = calculateContextBoost(commands[i].Context, current)
	}

	sortByContext(commands)
	return commands
}

func calculateContextBoost(cmdCtx, current CommandContext) int {
	score := 0

	score += directoryBoost(cmdCtx.Directory, current.Directory)
	score += repoBoost(cmdCtx.Repo, current.Repo)
	score += branchBoost(cmdCtx.Branch, current.Branch)
	score += workspaceBoost(cmdCtx.Workspace, current.Workspace)

	return score
}

func directoryBoost(cmdDir, currentDir string) int {
	if cmdDir == "" || currentDir == "" {
		return 0
	}
	if cmdDir == currentDir {
		return 400
	}
	if strings.HasPrefix(cmdDir, currentDir+string(filepath.Separator)) {
		return 250
	}
	if strings.HasPrefix(currentDir, cmdDir+string(filepath.Separator)) {
		return 150
	}
	return 0
}

func repoBoost(cmdRepo, currentRepo string) int {
	if cmdRepo == "" || currentRepo == "" {
		return 0
	}
	if cmdRepo == currentRepo {
		return 200
	}
	return 0
}

func branchBoost(cmdBranch, currentBranch string) int {
	if cmdBranch == "" || currentBranch == "" {
		return 0
	}
	if cmdBranch == currentBranch {
		return 150
	}
	return 0
}

func workspaceBoost(cmdWorkspace, currentWorkspace string) int {
	if cmdWorkspace == "" || currentWorkspace == "" {
		return 0
	}
	if cmdWorkspace == currentWorkspace {
		return 100
	}
	return 0
}

func sortByContext(commands []Command) {
	if len(commands) < 2 {
		return
	}
	sort.SliceStable(commands, func(i, j int) bool {
		if commands[i].ContextBoost == commands[j].ContextBoost {
			return commands[i].Index > commands[j].Index
		}
		return commands[i].ContextBoost > commands[j].ContextBoost
	})
}

