## shrimpy-binance

A relatively light weight wrapper for some [binance REST APIs](https://binance-docs.github.io/apidocs/#change-log) and some websocket APIs.
Endpoints are implemented on a "as needed" basis.

### RESTClient

RESTClient implements `services` that return `Responses`

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
 * [x] Implement authentication
 * [ ] Implement RateLimitHandling
 * [ ] Implement TimeHandler
 * [ ] Handle status codes
 * [ ] MetaDataREST
   * [ ] What should this be? 
   * [ ] I don't like the name.
   * [ ] Should I scrap responseHeader, and include it's fields here?
 * [ ] Client should be able to register MetaData channels 
   * [ ] RestMetaData -> Holds relevant data pertaining to requests 
   * [ ] WsMetaData -> Holds relevant data pertaining to ws streams
 * [ ] Implement WSClient
   * [ ] Implement Websocket Market Streams (Spot/Margin, Futures)
   * [ ] Market data requests (WebsocketAPI)
 * [x] Move ratelimit handler from RESTClient to Client 
       (WebSocket API, and Spot/Margin API share rate limits :( )
