package main

// ============================================================
//  C18 – Pull Request Workflow (Enterprise Style)
//  Objective : Work via Pull Requests like a real company team.
//
//  Full workflow executed by this script:
//
//    PHASE 1 – Bootstrap the repo on main
//      1.  Create README.md
//      2.  git init
//      3.  git add .               ← stages ALL files (README + main.go)
//      4.  git commit -m "first commit"
//      5.  git branch -M main
//      6.  git remote add origin   (removes old remote first if exists)
//      7.  git push -u origin main
//
//    PHASE 2 – Feature branch
//      8.  git checkout -b feature/add-description
//      9.  Add description.md (the feature work)
//      10. git add .
//      11. git commit -m "feat: add project description"
//      12. git push origin feature/add-description
//
//    PHASE 3 – Pull Request via GitHub CLI  (gh)
//      13. gh pr create  (open PR with title + body)
//      14. gh pr merge   (merge PR into main, delete branch)
//
//    PHASE 4 – Sync local main
//      15. git checkout main
//      16. git pull origin main
//
//  Pre-requisites:
//    - Git installed             (git --version)
//    - GitHub CLI installed      (gh --version)
//    - GitHub CLI authenticated  (gh auth login)
//    - Remote repo C18 already created as PRIVATE on GitHub
//
//  Usage:
//    go run automate.go
// ============================================================

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// ─────────────────────────────────────────────
//  CONFIGURATION  –  edit only these constants
// ─────────────────────────────────────────────

const (
	repoName  = "C18"
	remoteURL = "https://github.com/3D2Y-BK/C18.git"

	// Branch that carries the new feature
	featureBranch = "feature/add-description"

	// Pull Request title shown in the GitHub PR list
	prTitle = "feat: add project description file"

	// Pull Request body — explains WHAT changed and WHY
	prBody = "## What this PR does\n\n" +
		"- Adds `description.md` with the project purpose\n" +
		"- Demonstrates the full PR workflow: branch → push → PR → merge\n\n" +
		"## Why\n\n" +
		"Working via Pull Requests lets the team review code before it\n" +
		"reaches `main`, exactly like in a professional environment.\n\n" +
		"## Checklist\n" +
		"- [x] Branch created from main\n" +
		"- [x] Change committed with a clear message\n" +
		"- [x] PR description explains the change\n" +
		"- [x] Ready to merge"
)

// ─────────────────────────────────────────────
//  ANSI colours
// ─────────────────────────────────────────────

const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorRed    = "\033[31m"
	colorBold   = "\033[1m"
)

// ─────────────────────────────────────────────
//  Terminal helpers
// ─────────────────────────────────────────────

func phase(n int, title string) {
	fmt.Printf("\n%s┌─ PHASE %d · %s %s\n", colorCyan, n, title, colorReset)
}

var stepN int

func step(label string) {
	stepN++
	fmt.Printf("%s│  [%d] %s%s\n", colorCyan, stepN, label, colorReset)
}

func ok(msg string)   { fmt.Printf("%s│   ✔  %s%s\n", colorGreen, msg, colorReset) }
func info(msg string) { fmt.Printf("%s│   ➜  %s%s\n", colorYellow, msg, colorReset) }
func die(ctx string, err error) {
	fmt.Printf("%s✖  %s: %v%s\n", colorRed, ctx, err, colorReset)
	os.Exit(1)
}

// ─────────────────────────────────────────────
//  Shell helpers
// ─────────────────────────────────────────────

// run executes a command in `dir` and streams output live to the terminal.
// Prints the full command first so every action is visible to the professor.
// Exits the script immediately if the command fails.
func run(dir string, args ...string) {
	info("$ " + strings.Join(args, " "))
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		die("command failed: "+strings.Join(args, " "), err)
	}
	ok("success")
}

// runIgnoreError is like run() but silently ignores failure.
// Used for cleanup steps that are safe to skip if nothing exists yet
// (e.g. removing a remote that may not be set).
func runIgnoreError(dir string, args ...string) {
	info("$ " + strings.Join(args, " ") + "  (skip if fails)")
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run() // error intentionally ignored
}

// writeFile creates or overwrites a file with the given text content.
func writeFile(path, content string) {
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		die("cannot write "+path, err)
	}
	ok("file written: " + path)
}

// ─────────────────────────────────────────────
//  File contents
// ─────────────────────────────────────────────

const readmeContent = "# C18\n"

const descriptionContent = `# Project Description

## C18 – Pull Request Workflow

This project demonstrates how to collaborate on GitHub
using Pull Requests, exactly as teams do in companies.

## Steps Covered

| Step | Action                                       |
|------|----------------------------------------------|
| 1    | Create a feature branch from main            |
| 2    | Commit a change on that branch               |
| 3    | Push the branch to GitHub                    |
| 4    | Open a Pull Request with a clear description |
| 5    | Merge the PR into main                       |

## Why Pull Requests?

- Team members can **review** code before it reaches production.
- Each change is **documented** and traceable in the PR history.
- Conflicts are **resolved** before merging, not after.
`

// ─────────────────────────────────────────────
//  MAIN
// ─────────────────────────────────────────────

func main() {
	dir, err := os.Getwd()
	if err != nil {
		die("cannot get working directory", err)
	}

	fmt.Printf("\n%s══════════════════════════════════════════════%s\n", colorBold, colorReset)
	fmt.Printf("%s  C18 – Pull Request Workflow (Enterprise Style)%s\n", colorBold, colorReset)
	fmt.Printf("%s══════════════════════════════════════════════%s\n", colorBold, colorReset)
	fmt.Printf("  Directory  : %s%s%s\n", colorCyan, dir, colorReset)
	fmt.Printf("  Remote     : %s%s%s\n", colorCyan, remoteURL, colorReset)
	fmt.Printf("  PR branch  : %s%s%s\n", colorCyan, featureBranch, colorReset)

	// ════════════════════════════════════════════════════════
	//  PHASE 1 – Bootstrap repo on main
	// ════════════════════════════════════════════════════════
	phase(1, "Bootstrap repo on main")

	// Step 1 · Create README.md  (echo "# C18" >> README.md)
	step("Create README.md")
	writeFile(dir+"/README.md", readmeContent)

	// Step 2 · Initialise a local Git repository.
	// Safe to run on an existing repo — it simply reinitialises it.
	step("git init")
	run(dir, "git", "init")

	// Step 3 · Stage ALL files in the folder (README.md + main.go + any others).
	// Using "git add ." instead of "git add README.md" ensures there is always
	// something new to commit even when README.md content has not changed.
	step("git add .  (stage all files, not just README)")
	run(dir, "git", "add", ".")

	// Step 4 · Create the first commit — snapshot of the entire project
	step(`git commit -m "first commit"`)
	run(dir, "git", "commit", "-m", "first commit")

	// Step 5 · Rename the default branch to "main"
	// -M = force rename even if branch named main already exists
	step("git branch -M main")
	run(dir, "git", "branch", "-M", "main")

	// Step 6 · Set the remote to point to our GitHub repository.
	// We remove any existing "origin" first (ignored if it doesn't exist)
	// to avoid the "remote origin already exists" error.
	step("git remote add origin  (removes old remote first if any)")
	runIgnoreError(dir, "git", "remote", "remove", "origin") // safe cleanup
	run(dir, "git", "remote", "add", "origin", remoteURL)

	// Step 7 · Push main to GitHub for the first time.
	// -u = set upstream so plain "git push" works in future
	step("git push -u origin main")
	run(dir, "git", "push", "-u", "origin", "main")

	// ════════════════════════════════════════════════════════
	//  PHASE 2 – Feature branch
	//  In a team, nobody commits directly to main.
	//  Each change lives on its own branch until it is reviewed.
	// ════════════════════════════════════════════════════════
	phase(2, "Create feature branch and commit the change")

	// Step 8 · Create the feature branch and switch to it immediately.
	// Convention: feature/<short-description>
	step("git checkout -b " + featureBranch)
	run(dir, "git", "checkout", "-b", featureBranch)

	// Step 9 · Write description.md — the actual change this PR introduces
	step("Add description.md  ← the change that goes into the PR")
	writeFile(dir+"/description.md", descriptionContent)

	// Step 10 · Stage the new file
	step("git add .")
	run(dir, "git", "add", ".")

	// Step 11 · Commit using the Conventional Commits format.
	// "feat:" prefix signals a new feature to code reviewers.
	step(`git commit -m "feat: add project description"`)
	run(dir, "git", "commit", "-m", "feat: add project description")

	// Step 12 · Push the feature branch to GitHub so a PR can be opened against it
	step("git push origin " + featureBranch)
	run(dir, "git", "push", "origin", featureBranch)

	// ════════════════════════════════════════════════════════
	//  PHASE 3 – Pull Request  (GitHub CLI "gh")
	//  gh uses credentials from "gh auth login" — no token needed.
	// ════════════════════════════════════════════════════════
	phase(3, "Open and merge Pull Request on GitHub")

	// Step 13 · Open the PR: featureBranch → main
	//   --title  short headline visible in the GitHub PR list
	//   --body   full description of what changed and why
	//   --base   branch we are merging INTO  (main)
	//   --head   branch that holds our change (featureBranch)
	step("gh pr create  →  " + featureBranch + " → main")
	run(dir,
		"gh", "pr", "create",
		"--title", prTitle,
		"--body", prBody,
		"--base", "main",
		"--head", featureBranch,
	)

	// Step 14 · Merge the PR.
	//   --merge          standard merge commit (full history preserved)
	//   --delete-branch  delete the feature branch after merge (clean repo)
	step("gh pr merge  (merge commit + auto-delete branch)")
	run(dir,
		"gh", "pr", "merge",
		"--merge",
		"--delete-branch",
	)

	// ════════════════════════════════════════════════════════
	//  PHASE 4 – Sync local main
	//  The PR merged on GitHub; local main is behind by one commit.
	//  Pulling brings local main up to date with remote.
	// ════════════════════════════════════════════════════════
	phase(4, "Sync local main with the merged commit")

	// Step 15 · Switch back to main
	step("git checkout main")
	run(dir, "git", "checkout", "main")

	// Step 16 · Download the merge commit — local main now equals remote main
	step("git pull origin main")
	run(dir, "git", "pull", "origin", "main")

	// ── Summary ───────────────────────────────────────────
	fmt.Printf("\n%s══════════════════════════════════════════════%s\n", colorBold, colorReset)
	fmt.Printf("%s🎉  PR Workflow complete!%s\n\n", colorGreen, colorReset)
	fmt.Printf("  %-22s %s%s%s\n", "Repo:", colorCyan, remoteURL, colorReset)
	fmt.Printf("  %-22s %s%s%s\n", "Branch merged:", colorGreen, featureBranch, colorReset)
	fmt.Printf("  %-22s %sfeat: add project description → main%s\n", "Change:", colorYellow, colorReset)
	fmt.Printf("%s══════════════════════════════════════════════%s\n\n", colorBold, colorReset)
}
