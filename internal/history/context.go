package history

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// CurrentContext inspects the running environment and returns metadata for boosting.
func CurrentContext() CommandContext {
	cwd, err := os.Getwd()
	if err != nil {
		return CommandContext{}
	}
	return ResolveContext(cwd)
}

// ResolveContext builds a CommandContext for an arbitrary directory.
func ResolveContext(dir string) CommandContext {
	if dir == "" {
		return CommandContext{}
	}

	absDir, err := filepath.Abs(dir)
	if err != nil {
		absDir = dir
	}
	absDir = filepath.Clean(absDir)

	ctx := CommandContext{
		Directory: absDir,
	}

	if repoRoot := gitOutput(absDir, "rev-parse", "--show-toplevel"); repoRoot != "" {
		ctx.Workspace = filepath.Clean(repoRoot)
		ctx.Repo = filepath.Base(ctx.Workspace)
		ctx.Branch = gitOutput(absDir, "rev-parse", "--abbrev-ref", "HEAD")
	}

	return ctx
}

func gitOutput(dir string, args ...string) string {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GIT_OPTIONAL_LOCKS=0")

	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

