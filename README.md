# Bonbaan-BE

Bonbaan-BE is the backend service for the Bonbaan project.

## Project Structure

The project follows a clean architecture pattern with the following main directories:

- `src/`: Contains the main application code
  - `adepters/`: Adapters for external services and databases
  - `Config/`: Configuration files
  - `Constance/`: Constant values used across the project
  - `entities/`: Domain entities
  - `model/`: Data models
  - `routers/`: HTTP route handlers
  - `usecases/`: Business logic
  - `utils/`: Utility functions
- `tests/`: Contains test helpers and utilities

## Getting Started

1. Clone the repository
2. Install dependencies: `go mod download`
3. Set up your `.env` file based on the `.env.example`
4. Run the application: `go run main.go`

## Testing

We use Go's built-in testing framework for unit tests. Our testing structure is designed to be easy to use and maintain.

### Running Tests

To run all tests:

```bash
go test -v ./...
```

For more details on our testing structure and how to write tests, please refer to the [tests/README.md](tests/README.md) file.

### Continuous Integration

We use GitHub Actions for continuous integration. The workflow is defined in `.github/workflows/go-test.yml`. It runs tests, generates coverage reports, and uploads them to Codecov.

## Contributing

1. Fork the repository
2. Create your feature branch: `git checkout -b feature/my-new-feature`
3. Commit your changes: `git commit -am 'Add some feature'`
4. Push to the branch: `git push origin feature/my-new-feature`
5. Submit a pull request

## License
.
[MIT License](LICENSE)
