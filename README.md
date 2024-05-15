My personal notes regarding the codebase can be found in the `notes/` folder. <br />
I intend to extend the Go blueprint to learn more about the standard http library.

# Extensions

- [ ] Log a message when the server is started
- [ ] Create an API point that serves data from the database
- [ ] Create an API point to `INSERT`s data in the database

# Project go-blueprint

One Paragraph of project description goes here

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

## MakeFile

run all make commands with clean tests

```bash
make all build
```

build the application

```bash
make build
```

run the application

```bash
make run
```

Create DB container

```bash
make docker-run
```

Shutdown DB container

```bash
make docker-down
```

live reload the application

```bash
make watch
```

run the test suite

```bash
make test
```

clean up binary from the last build

```bash
make clean
```
