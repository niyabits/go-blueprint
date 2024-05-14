# Exploration

I start by exploring different files and the project structure.
I am mostly writing about my thought process of exploring the project and syntax/technqiues that intrigue my interest.

## Server

The obvious Go file to start exploring is `cmd/api/main.go` as it is the entry point of this module.

The `main.go` file calls the `server/server.go` package.

## `server.go`

### Blank Identifier

On the top of `server.go` there is a [blank identifier](https://go.dev/ref/spec#Blank_identifier) for the autoload library, I looked it up in the Go specs, and the blank identifier tells Go that this module is loaded to just call the `init` function and this module won't be called in the code.

### net/http Library

I try to understand the different fields in the `&http.Server` starting with `Addr`.

`Addr` takes in a `string` in the form of 'host:port', interestingly in the code `:8080` is passed as a value.

The [`net/http`](https://pkg.go.dev/net/http#Server) docs for `Addr` says that the format of the string is specified by [`net/Dial`](https://pkg.go.dev/net#Dial).
Which tells me if the host is empty then the local system is assumed.

The `handler` value is a method named `RegisterRoutes()`, the method is defined on the `Server` struct and in a different file called `routes.go`

I found the code a bit jarring to read because it was not explicit where the `RegisterRoutes` method is defined.

**Anti-Pattern**: The method is described in the `routes.go` file and not in the same file as the type.

The Type, Factory Functions, and Methods should be defined together in the chronological order in the same file, the methods should not be scattered around.

This is something that I learnt in the Ultimate Go course by Ardan Labs. Here are details of where I learnt this: <br />
**Timestamp**: [8:24](https://courses.ardanlabs.com/courses/take/ultimate-go/lessons/7419439-4-1-1-methods-part-1-value-pointer-semantics)<br />
**Video**: 4.1.1 - Methods-Part 1 (Value & Pointer Semantics)

### `http.Handler`

The [`Handler`](https://pkg.go.dev/net/http#Handler) type responds to an HTTP request. <br/>
A `ServeMux` can be used as an `Handler` because it implements [`ServerHTTP`](https://pkg.go.dev/net/http#ServeMux.ServeHTTP).

`ServerMux` basically maps the incoming request's URL patterns with the list of registered patterns and calls the appropriate function.

Now we can use this `ServerMux` to execute functions according to the route on which we recieved a request. <br />
The functions we execute can use the data from the request and write a response back. The signature of these functions look like: `handler func(http.ResponseWriter, *http.Request)` as specified by the `mux.HandleFunc()`.

## Database

The [`pgx`](https://pkg.go.dev/github.com/jackc/pgx/v5@v5.5.5#section-readme) (Postgres Go Driver) and `godotenv` are imported with a blank identifier.

Using the `database/sql` interfaces for connection is an interesting choice, the documentation says the `pgx` interface is faster and is recommended when –

- The application only targets PostgreSQL.
- No other libraries that require database/sql are in use.

The Go module does not seem to use `database/sql` for any other libraries and is only targetting PostgreSQL, I think using `database/sql`'s interface guarantees compatibility with other databases too and also future proofing in case we want to use a library that uses `database/sql`

The `init` function of `pgx` implements the [`database/sql/driver`](https://pkg.go.dev/database/sql/driver) and [loads it](https://github.com/jackc/pgx/blob/523411a3fbcb76daf4eb3bc60326e2a68115e92f/stdlib/sql.go#L94). <br /> I could not find this information in the `pgx` documentation, so I thought it might be helpful to understand why we do an import with a blank identifier.

```go
ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
```

I don't completely understand `context` in Go, specifically why `context.Background()` is used instead of `context.Context` as in the [documentation of the `PingContext`](https://pkg.go.dev/database/sql#DB.PingContext) methood. It's worth exploring the [Go concurreny patern `context`](https://go.dev/blog/context).

The `Health` and `Close` functions are pretty straightforward, I can see how struct types can be utilized to create abstractions over an existing API.
