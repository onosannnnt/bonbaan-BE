# Testing Structure

This directory contains the test helper and utilities for running unit tests in the Bonbaan-BE project.

## Structure

- `test_helper.go`: Contains setup and teardown functions for the test environment, including loading the `.env.test` file.

## Running Tests

To run the tests locally:

1. Ensure you have a `.env.test` file in the root directory with the necessary environment variables.
2. Run the following command from the project root:

```bash
go test -v ./...
```

This will run all tests in the project, including those that use the test helper.

## Writing Tests

When writing new tests:

1. Import the `tests` package in your test file:

```go
import (
    "testing"
    "github.com/onosannnnt/bonbaan-BE/tests"
)
```

2. Use the `TestMain` function to set up the test environment:

```go
func TestMain(m *testing.M) {
    tests.TestMain(m)
}
```

3. Write your test functions as usual.

## Environment Variables

The `SetupTestEnvironment` function in `test_helper.go` loads the `.env.test` file. Ensure that all necessary environment variables for testing are included in this file.

## Continuous Integration

The project uses GitHub Actions for continuous integration. The workflow is defined in `.github/workflows/go-test.yml`. It runs tests, generates coverage reports, and uploads them to Codecov.
