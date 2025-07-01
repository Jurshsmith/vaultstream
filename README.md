# VaultStream

## Quick Start

This may not work depending on your system configuration and pre-installed software. If you encounter any issues, refer to the [Prerequisites](#prerequisites) section to ensure that all necessary dependencies are installed before re-running this command.

Quickly set up and start VaultStream with:

    make quick-start

Adjust the `BATCH_SIZE` in the generated `.env` to observe the efficient of the systems architecture at different batch sizes. If you encounter issues even after installing the prerequisites, use `make help` to see additional debugging commands.

### Available tables in the PostgresDB

- `records`
- `signatures`

### Available streams in the NATS Jetstream Setup

- `records.>`
- `keys.>`

## Prerequisites

Before running the project, ensure you have the following tools installed. For detailed installation instructions, please refer to each tool’s official documentation:

- **Golang**
  - Install from [golang.org/dl](https://golang.org/dl/)
  - Minimum supported Go version is `v1.18.0`
- **Make**
  - **macOS:** Install via Xcode Command Line Tools. See [Apple’s documentation](https://developer.apple.com/xcode/)
  - **Linux:** Install via your package manager. See [GNU Make](https://www.gnu.org/software/make/) for details
- **Docker (v20.10 or later)**
  - Follow the instructions at [Docker’s Get Docker](https://docs.docker.com/get-docker/)

## Future Enhancements (Non-Exhaustive)

- **Implement Robust Retry Strategies**  
  Incorporate retries with jitter and exponential backoff to gracefully handle transient errors and unsignable records.

- **Enhance Observability**  
  Integrate comprehensive monitoring, logging, and alerting solutions to improve system visibility and facilitate proactive diagnostics.

- **Deploy with Kubernetes & Helm**  
  Develop and maintain Helm charts for deploying the system on Kubernetes, to further demonstrate how well `signing-service` autoscales

- **Expand Unit Test Coverage**  
  Increase test coverage for critical modules to ensure reliability and detect issues early in the development cycle.

- **Strengthen Integration Testing**  
  Enhance end-to-end integration tests to validate overall system functionality, performance, and interoperability.

- **Improve readability**  
  Break down logical units of the services to improve readability and reusability

- **Ensure Graceful Shutdown Procedures**  
  Implement robust shutdown mechanisms to maintain system stability during deployments and unexpected terminations.

## Architecture Diagram (Quick Draft)

Below is a high-level overview of the systems architecture:

![VaultStream Architecture Diagram](./docs/vaultstream-architecture.png "VaultStream Architecture Diagram")
