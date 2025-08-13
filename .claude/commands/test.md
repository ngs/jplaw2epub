Run all tests and ensure they pass with at least 70% code coverage.

Instructions:
1. Run the test suite with coverage reporting
2. If tests fail, analyze the failures and fix them
3. If coverage is below 70%, add more tests to increase coverage
4. Repeat until all tests pass and coverage is at least 70%
5. Show the final test results and coverage report

Requirements:
- All tests must pass (exit code 0)
- Code coverage must be at least 70%
- Fix any failing tests by updating the code or test cases as appropriate
- Add new tests if coverage is insufficient

Workflow:
1. Run: `go test -v -coverprofile=coverage.out ./...`
2. Check coverage: `go tool cover -func=coverage.out`
3. If failures or coverage < 70%, fix issues and repeat
4. Continue until all requirements are met

Notes:
- Prioritize fixing test failures before improving coverage
- When adding tests, focus on uncovered critical paths
- Ensure new tests are meaningful and not just for coverage