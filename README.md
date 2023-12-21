# shrimpy-binance

`Shrimpy-binance` is a lightweight Go wrapper for some of Binance's [REST](https://binance-docs.github.io/apidocs/#change-log) and [Websocket](https://binance-docs.github.io/apidocs/#change-log) APIs.
It simplifies interaction with REST and Websocket endpoints by providing streamlined interfaces and useful features.
Endpoints are implemented on an "as needed" basis.

## Features

 - [x] **Websocket Streams**
 - [x] **REST Services**
 - [x] **Response Parsing**: all responses from websocket streams and rest services are parsed into structs.
 - [x] **Rate Limit Management**: keep track of rate limits across diffenent endpoint types and rate limit types.
 - [x] **Server Time Synchronization**: default implementation for synchronizing shrimpy-binance  client with server time.
 - [ ] **Orderbook Management**: default implementation for maintaining local copies of orderbooks.

## Installation

to install shrimpy-binance use go get:

```sh
go get github.com/svdro/shrimpy-binance
```

## Usage

### REST API Usage

```golang
import (
    "context"
    "fmt"
	shrimpy "github.com/svdro/shrimpy-binance"
)

client := shrimpy.BinanceClient("apiKey", "apiSecret")
if resp, err := client.NewSpotMarginServerTimeService().Do(context.Background()); err == nil {
    fmt.Printf("%+v\n", resp)
}
```

## Websocket API Usage

```golang
import (
    "context"
    "fmt"
	shrimpy "github.com/svdro/shrimpy-binance"
)

client := binance.BinanceClient("apiKey", "apiSecret")
stream := client.NewSpotMarginDiffDepth100Stream().SetSymbol("BTCUSDT")
go func() {
    for {
        select {
        case event := <-stream.Handler.EventChan:
            fmt.Printf("%v+\n", event)
        case err := <-stream.Handler.ErrChan:
            fmt.Println(err)
            return
        }
    }
}()
stream.Run(context.Background())

```

## TODOs:
 * [ ] Make RLH and TH accessible for users of client.
 * [ ] Monitorine: Client should be able to register MetaData channels 
 * [ ] Test RateLimitHandler in live setting
 * [ ] Services
 * [ ] Fix sapi rate limit issues
   * [ ] ExchangeInfoService
   * [ ] OrderService

## Issues with SAPI Rate Limits

Currently rate limits for sapi endpoints are not handled correctly. 
Until this is fixed set sapi rate limits to -1 in client options and hope for 
the best :(.

This is because the logic for counting `sapi` rate limits is different from that 
for counting other rate limits (e.g. `api`, `fapi`, `dapi`, etc).
- usually all endpoints of the same type (e.g. /api/v3/serverTime, /api/v3/ping)
  contribute to a shared rate limit for that endpoint type (`api` in this case).
- for `sapi` endpoints binance uses a distinct rate limit for each endpoint.
  This means that different `sapi` endpoints (e.g. /sapi/v1/system/status,
  /sapi/v1/userDataStream) each have their own rate limits.

To fix this, the logic for sapi rate limit counters needs to be changed fundamentally.
- rate limits for `sapi` cannot defined in client options, but should logically 
  be defined in the service definition for each `sapi` service instead, because 
  different sapi endpoints can have different max rate limits.
- `sapi` rate limit counters can no longer created the client options when the 
  client is initialized. Instead `sapi` rate limit counters could be created 
  dynamically the first time a `sapi` endpoint is called.
