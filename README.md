My personal notes regarding the codebase can be found in the `notes/` folder. <br />
I intend to extend the Go blueprint to learn more about the standard http library.

# Extensions

- [x] Log Message when the Server Starts
- [x] SQL File to Create Tables
- [x] GET: get all albums
- [x] GET: get an album by ID
- [x] POST: add an album
- [x] DELETE: delete an album
- [ ] Gracefully Handle server errors by sending appropriate responses
- [ ] Add tests for all the new endpoints

## Connect Database with `psql`

Since `psql` hostname defaults to local socket we need to explicitly pass in the hostname as 'localhost'.
Passing the username as `melkey` will give root access. The password is in the `.env`

These details are specified in `.env` and passed to `docker-compose.yml`

```bash
psql -h localhost blueprint melkey
```

## Execute `create-tables.sql` with `psql`

`\i` is an internal `psql`command that can be used to execute SQL in files that are in the same filesystem as the client.

> Aside: This is in contrast to SQL's [`COPY` command](https://www.postgresql.org/docs/16/sql-copy.html)

```
\i create-tables.sql
```

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
