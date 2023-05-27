module github.com/jxo-me/netx/sdk

go 1.20

replace (
	github.com/jxo-me/netx/core => ../core
	github.com/jxo-me/netx/gosocks4 => ../gosocks4
	github.com/jxo-me/netx/gosocks5 => ../gosocks5
	github.com/jxo-me/netx/plugin => ../plugin
	github.com/jxo-me/netx/relay => ../relay
	github.com/jxo-me/netx/tls-dissector => ../tls-dissector
	github.com/jxo-me/netx/x => ../x
)
