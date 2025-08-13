Update the title and description of the current pull request based on recent changes.

Instructions:
1. First, check if there's an open PR for the current branch using `gh pr view`
2. Analyze recent commits and code changes to automatically generate:
   - PR title based on the main features/changes
   - PR description with a summary of changes, files modified, and test coverage
3. Use `gh pr edit` to update the PR with the generated title and description
4. Show the PR URL after updating

Requirements:
- GitHub CLI (gh) must be installed and authenticated
- Must be on a feature branch (not master/main)
- PR must exist for the current branch

Example workflow:
- Get PR number: `gh pr view --json number -q .number`
- Update title: `gh pr edit [PR_NUMBER] --title "New Title"`
- Update description: `gh pr edit [PR_NUMBER] --body "New Description"`
- Update both: `gh pr edit [PR_NUMBER] --title "New Title" --body "New Description"`