# okrforjira

![](https://github.com/grandper/okrforjira/workflows/Test/badge.svg)

A golang client for [OKR for Jira](https://digitaltoucan.com/products/okr-for-jira).
This client implements the API described in
* [OKR API documentation](https://intercom.help/okr-for-jira/en/articles/6178256-okr-api-documentation)
* [API query methods](https://intercom.help/okr-for-jira/en/articles/6178378-api-query-methods)
* [API update methods](https://intercom.help/okr-for-jira/en/articles/6252250-api-update-methods)

Digital Toucan is a company that continuously improves *OKR for Jira*.
Although good care has been provided to this package, it is possible that at some point the API evolves and this package becomes outdated.
Finally, this package is provided to you without any guarantee. Use it at your _own risk_.

# Install

`go get github.com/grandper/okrforjira`

## Usage

```go
c := okrforjira.NewClient(nil, token)

expand := []string{"TEAMS"}
resp1, err := c.ObjectivesByDate(ctx, startDate, deadline, expand)
resp2, err := c.ObjectivesByIDs(ctx, objectiveIDs, expand)
resp3, err := c.KeyResultsByDate(ctx, startDate, deadline, expand)
resp4, err := c.KeyResultsByIDs(ctx, keyResultIDs, expand)
```

## Example

```go
package main

import (
    "context"
    "flag"
    "fmt"
    "os"
    "time"

    "github.com/grandper/okrforjira"
)

func main() {
    token := flag.String("token", "", "token to access your OKR data")
    flag.Parse()

    ctx := context.Background()
    c := okrforjira.NewClient(nil, *token)

    startDate := time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC)
    deadline := time.Date(2022, time.June, 30, 0, 0, 0, 0, time.UTC)
    expand := []string{"TEAMS"}
    resp, err := c.ObjectivesByDate(ctx, startDate, deadline, expand)
    if err != nil {
        fmt.Printf("failed to get objectives by date: %s\n", err.Error())
        os.Exit(1)
    }
    fmt.Printf("%#v", resp)
}
```
## Note

**okrforjir** is using Go 1.18
