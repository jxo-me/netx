module github.com/jxo-me/netx/plugin

go 1.20

require (
	google.golang.org/grpc v1.55.0
	google.golang.org/protobuf v1.30.0
)

replace (
	github.com/jxo-me/netx/core => ../core
	github.com/jxo-me/netx/gosocks4 => ../gosocks4
	github.com/jxo-me/netx/gosocks5 => ../gosocks5
	github.com/jxo-me/netx/plugin => ../plugin
	github.com/jxo-me/netx/relay => ../relay
	github.com/jxo-me/netx/tls-dissector => ../tls-dissector
)

require (
	github.com/golang/protobuf v1.5.3 // indirect
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/text v0.8.0 // indirect
	google.golang.org/genproto v0.0.0-20230306155012-7f2fa6fef1f4 // indirect
)
