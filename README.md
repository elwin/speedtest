# speedtest on SCION using QUIC

This speedtest is used to evaluate the available bandwidth when using QUIC for the transport layer.

## Installation
```go
go get github.com/elwin/speedtest
go install github.com/elwin/speedtest/speedtest_server
go install github.com/elwin/speedtest/speedtest_client
```

## Server
```
$ speedtest_server -h
Usage of speedtest_server:
  -local string
    	Local address (with Port)

```

## Client
```
$ speedtest_client -h
Usage of speedtest_client:
  -local string
    	Local address (without Port)
  -remote string
    	Remote address (with Port)
```