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
 * [ ] Monitoring: Client should be able to register MetaData channels 
 * [ ] Test RateLimitHandler in live setting
