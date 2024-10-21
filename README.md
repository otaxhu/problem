# HTTP Problem Details for Golang

Library for parsing and creating HTTP Problem Details, following RFC 9457 document, for ensuring compatibility between other libraries.

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
        for _, ac := range mp["bank_accounts"] {
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
