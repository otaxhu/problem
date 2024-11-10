# HTTP Problem Details for Golang

Library for parsing and creating HTTP Problem Details, following [RFC 9457 document](https://www.rfc-editor.org/rfc/rfc9457.html), for ensuring compatibility between other libraries.

<div>
  <a href="https://pkg.go.dev/github.com/otaxhu/problem" style="text-decoration:none;">
    <img src="https://pkg.go.dev/badge/github.com/otaxhu/problem" alt="Go Reference">
  </a>
  <a href="https://coveralls.io/github/otaxhu/problem?branch=main" style="text-decoration:none;">
    <img src="https://coveralls.io/repos/github/otaxhu/problem/badge.svg?branch=main" alt="Coverage Status">
  </a>
  <a href="https://goreportcard.com/report/github.com/otaxhu/problem" style="text-decoration:none;">
    <img src="https://goreportcard.com/badge/github.com/otaxhu/problem" alt="Go Report Card">
  </a>
  <a href="https://github.com/otaxhu/problem/actions/workflows/ci.yml">
    <img src="https://github.com/otaxhu/problem/actions/workflows/ci.yml/badge.svg?branch=main" alt="CI Status">
  </a>
</div>

## Installation
```sh
$ go get github.com/otaxhu/problem
```

## Features

- ### JSON and XML support:

  You can encode/decode Problem Details in JSON and XML format using standard packages `encoding/json` and `encoding/xml`.

  XML encoding/decoding has some limitations regarding `MapProblem` implementation, please read the [docs](https://pkg.go.dev/github.com/otaxhu/problem) to know more details.

- ### HTTP Client APIs:

  As an HTTP client, you can parse HTTP Problem Details responses using `ParseResponse()` and `ParseResponseCustom()` functions.

- ### HTTP Server APIs:

  As an HTTP server, you can respond to clients with Problem Details responses, using any of the available `Problem` interface implementations, encoding it using `ServeJSON()` or `ServeXML()` helpers.

- ### Polymorphic Problem Details and Easy extension members:

  You can embed `RegisteredProblem` struct in your own struct, and extend it with any members you want, as allowed by [RFC 9457 Section 3.2](https://www.rfc-editor.org/rfc/rfc9457.html#name-extension-members)

## Quick Usage:

### Client code:

```go
import (
    "net/http"

    "github.com/otaxhu/problem"
)

func main() {
    res, _ := http.Get("https://example.org/endpoint")

    if res.StatusCode != http.StatusOK {
        p, _ := problem.ParseResponse(res)
        fmt.Println(p.GetDetail())
        // You can also retrieve extension members by doing the following
        mp := p.(*MapProblem)

        // bank_accounts is a list of accounts
        for _, ac := range mp["bank_accounts"].([]any) {
            // Do something...
        }
        os.Exit(1)
    }

    fmt.Println("Everything is OK")
}
```

### Server code:

```go
import (
    "net/http"

    "github.com/otaxhu/problem"
)

func main() {
    http.HandleFunc("GET /endpoint", func(w http.ResponseWriter, r *http.Request) {
        if r.Header.Get("X-Application-Specific") == "bad foo" {
            p := problem.NewMap(http.StatusBadRequest, "bad parameter, please try again")

            // Setting a extension member
            p["more_info"] = []string{"foo", "bar", "baz"}
            problem.ServeJSON(p).ServeHTTP(w, r)
            return
        }

        w.WriteHeader(http.StatusOK)
    })

    http.ListenAndServe(":8080", nil)
}
```

### Client Code (Custom Problem):

```go
import (
    "net/http"

    "github.com/otaxhu/problem"
)

type CustomProblem struct {
    problem.RegisteredProblem
    ExtensionMember string `json:"extension_member" xml:"extension_member"`
}

func main() {
    res, _ := http.Get("https://example.org/endpoint")

    if res.StatusCode != http.StatusOK {
        var p CustomProblem
        _ = problem.ParseResponseCustom(res, &p)

        // You get p populated with the Problem Details
        fmt.Println(p.Detail, p.ExtensionMember)

        os.Exit(1)
    }

    fmt.Println("Everything is OK")
}
```

### Server Code (Custom Problem):

```go
import (
    "net/http"

    "github.com/otaxhu/problem"
)

type CustomProblem struct {
    problem.RegisteredProblem
    ExtensionMember string `json:"extension_member" xml:"extension_member"`
}

func main() {
    http.HandleFunc("GET /endpoint", func(w http.ResponseWriter, r *http.Request) {
        if r.Header.Get("X-Application-Specific") == "bad foo" {
            p := &CustomProblem{
                // NewRegistered Returns a pointer, you need to derefence it using *
                RegisteredProblem: *problem.NewRegistered(http.StatusBadRequest, "bad parameter"),
                ExtensionMember:   "bar",
            }

            // Using JSON:
            problem.ServeJSON(p).ServeHTTP(w, r)

            // Using XML:
            //
            // problem.ServeXML(p).ServeHTTP(w, r)
            return
        }

        w.WriteHeader(http.StatusOK)
    })

    http.ListenAndServe(":8080", nil)
}
```

## Contributing

Any kind of contribution is totally welcome, thank you very much for any contribution to this project.

Any pull request that fixes little problems like "typos", or enhances documentation; does not need to open an issue in advance.

Otherwise, if the pull request adds or modifies features, please make sure there is a open related issue. If there is no issue, then the pull request will likely be rejected.

To start contributing, make your own fork of the repository, make sure to pull all of the commits from `main`, create a new branch from `main`, commit your changes, push changes to your fork, and start the pull request against `main`. Reviews of pull request may take a time, please be patient :wink:.
