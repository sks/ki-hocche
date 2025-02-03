# Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.

## Code of Conduct

[ ]: #
This project and everyone participating in it is governed by the [Ki-Hocche Code of Conduct](https://go.dev/conduct). By participating, you are expected to uphold this code. Please report unacceptable behavior in github.

### Prerequisites

- [Go](https://go.dev/doc/install)
- [Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)
- [GitHub account](https://github.com)
- [golangci-lint](https://golangci-lint.run/usage/install/) for linting
- [gofmt](https://pkg.go.dev/cmd/gofmt) for formatting

### Getting Started

1. Fork the repository on GitHub.
2. Clone your fork:

   `git clone https://github.com/sks/kihocche`

3. Create a new branch for your changes:

   `git checkout -b my-feature-branch`

4. Make your changes.

5. Commit your changes:

   `git commit -m "Add my feature"`

6. Push to your fork:

   `git push origin my-feature-branch`

7. Open a pull request on GitHub.

### Running Tests

To run the tests, use the following command:

    go test ./...

### Linting

To lint the code, use the following command:

    golangci-lint run

### Formatting

To format the code, use the following command:

    gofmt -s -w .

## License

By contributing to Ki-Hocche, you agree that your contributions will be licensed under its [Apache License 2.0](https://www.apache.org/licenses/LICENSE-2.0).

## Contact

For any questions or suggestions, please open an issue on the [GitHub repository](https://github.com/sks/kihocche/issues).

Thank you for contributing to Ki-Hocche!
