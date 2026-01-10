
## PROJECT CONVENTIONS

- **Git repo cache**: All cloned repositories go in `./.cache` directory (NOT /tmp)
  - Example: `.cache/github.com/open-telemetry/opentelemetry-android`
  - The GitFetcher handles this automatically when using adapters
  - For manual exploration, clone to `.cache/` subdirectory
- **Database**: SQLite at `./data/metric-library.db`
- **Binary**: `./bin/metric-library`

## WAY OF WORKING

- Start with minimal functionality and verify it works before adding complexity
- For all compiled languages please compile after each change.
- Do not leave code with compile errors.
- Once you are done making a change, kindly run linting and fix any errors.
- Follow Test Driven Development
   -  Make sure to add test cases before you make a change.
   -  Be kind and run a failing test, fix it and then run test again.
   -  commit when tests are passing.
- Please make sure to run tests before committing.
- Please make sure after Compilation and Linting that the tests are passing before reporting any success.
- Kindly avoid stating things like "it works", if you want to show, show me green tests.
- I cannot request this enough, please make sure to run tests after every change.
- I prefer trunk based development, and git for version control.
- Prefer latest version of libs unless there is a reason not to.
- Use dbmate when any table is added or changed.
- Use Makefile for all build commands

## building and running the project
- Always use Makefile for building and running the project
- use make commands like `make build`, `make run`, `make test`, etc.
- do not add "Generated with" and "Co-Authored-By" lines to commit messages

## CI
- Use GitHub Actions for CI/CD
- Location of CI configuration files: `.github/workflows`

## Golang specific instructions

1. Please use go mod for dependency management
2. Please use go fmt for code formatting
3. Please use go test for running tests
4. Please use golangci-lint for linting
5. Please use Makefile for all commands, and for ci
6. Please default to Go style guide https://google.github.io/styleguide/go for
   naming and formatting, and best practices -
   https://google.github.io/styleguide/go/best-practices
7. Avoid writing comments unless necessary
8. Use https://github.com/uber-go/dig for dependency injection
9. Use dbmate for SQL migrations
10. Follow Test Driven Development - write tests first, see them fail, then
    implement functionality to make them pass
    
### Development Practices
- Test Driven Development with failing tests first
- Makefile for all build/test/run commands
- Go modules for dependencies
- golangci-lint for code quality
- dbmate for migrations (Memgraph schema setup)
- Uber's dig for dependency injection
- Minimal functionality first, then add complexity
- Compile after each change, no compile errors allowed
- Run tests before any success reporting
