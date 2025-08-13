Run linters and fix all issues until the code passes all lint checks.

Instructions:
1. Run all available linters for the project
2. If any linting issues are found, fix them automatically or manually
3. Re-run linters after fixes to verify all issues are resolved
4. Repeat until all linters pass with no errors or warnings
5. Show the final clean lint results

Requirements:
- All linters must pass with exit code 0
- No errors or warnings should remain
- Code must follow project style guidelines
- All auto-fixable issues should be resolved

Workflow for Go projects:
1. Run: `gofmt -l -w .` to format code
2. Run: `go vet ./...` to check for suspicious constructs
3. Run: `golint ./...` or `golangci-lint run` if available
4. Run: `go mod tidy` to clean up dependencies
5. Fix any reported issues
6. Repeat until all checks pass

Common fixes:
- Format issues: Use `gofmt -w` to auto-fix
- Import ordering: Use `goimports -w` if available
- Unused variables/imports: Remove or use them
- Missing comments: Add documentation for exported functions
- Error handling: Ensure all errors are properly handled

Notes:
- Apply auto-fixes first when available
- For manual fixes, understand the issue before making changes
- Ensure fixes don't break existing functionality