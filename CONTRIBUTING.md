# Contributing to Heartbeat Operator

First off, thank you for considering contributing to Heartbeat Operator! It's people like you that make Heartbeat Operator an awesome Kubernetes tool.

## Where to Start?

1.  **Issues:** Check the issue tracker for bugs or requested features.
2.  **Discussions:** If you have an idea, feel free to open an issue or discussion first so we can align on the implementation approach.
3.  **Fork:** Ensure you fork the repo and create your branch from `main`.

## Local Development

We have built-in Make targets to get you started quickly:

*   `make build` - Compiles the operator binary locally.
*   `make run` - Runs the operator against your currently configured Kubernetes context (e.g. your Minikube, Kind, or Docker Desktop cluster).
*   `make test` - Runs the unit tests.
*   `make lint` - Runs `golangci-lint` to ensure code aligns with our style guidelines.

See `Makefile` for more details.

## Pull Request Guidelines

1. Please make sure that `make verify` passes before submitting your PR.
2. If your PR introduces a new feature or fixes a bug, please ensure that you add adequate test coverage.
3. Update the `README.md` if any user-facing APIs or behaviors change.
4. Keep commit messages clear, descriptive, and reference associated issues if applicable.

We look forward to reviewing your PRs!
