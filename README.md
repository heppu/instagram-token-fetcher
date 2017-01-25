# Instagram token fetcher
This repository offers module for Go project to fetch instagram access tokens programmatically.
It also has a command line tool for getting access tokens.

## Usage

In code:
```
go get github.com/heppu/instagram-token-fetcher
```

```go
package main

import (
    "log"

    "github.com/heppu/instagram-token-fetcher"
    "gopkg.in/alecthomas/kingpin.v2"
)

func main() {
    tokenClient := igtoken.NewClient(
        "my-app-id",
        "http://some-url.com:8888/token",
        "some-user",
        "some-passwd",
    )

    token, err := tokenClient.GetToken(igtoken.PUBLIC_CONTENT)
    if err != nil {
        log.Fatal(err)
        return
    }
    log.Print(token)
}
```

From command line

```
go install github.com/heppu/instagram-token-fetcher/cmd/ig-token
ig-token --help
```