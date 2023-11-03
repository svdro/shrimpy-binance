## shrimpy-binance

A relatively light weight wrapper for some [binance REST APIs](https://binance-docs.github.io/apidocs/#change-log) and some websocket APIs.
Endpoints are implemented on a "as needed" basis.

### RESTClient

RESTClient implements `services` that return `Responses`.

#### Endpoints
 * [ ] `/api/v3/ping`
 * [x] `/api/v3/time` 

### WSClient

WSClient implements `streams` that generate `"Events"`

### WebSocketAPI

WebSocketAPI is a hybrid between a binance rest api and a websocket api.
Effectively I want to use this like a RESTClient, but use an underlying
websocket connection.


#### TODO:
 * [ ] Client should be able to register MetaData channels 
   * [ ] RestMetaData -> Holds relevant data pertaining to requests 
   * [ ] WsMetaData -> Holds relevant data pertaining to ws streams
 * [ ] Implement RateLimitHandling
 * [ ] Implement TimeHandler
 * [ ] Implement RestClient
   * [x] ServiceMeta
   * [x] Implement authentication
   * [ ] Implement Handle status codes
 * [ ] Implement WSClient
   * [ ] Implement Websocket Market Streams (Spot/Margin, Futures)
   * [ ] Market data requests (WebsocketAPI)
 * [ ] Implement WSAPIClient
   * [ ] Hybrid between RESTClient and WSClient 
   * [ ] Only really need to implement market data endpoints (serverTime).
 * [x] Move ratelimit handler from RESTClient to Client 
       (WebSocket API, and Spot/Margin API share rate limits :( )
