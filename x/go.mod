module github.com/jxo-me/netx/x

go 1.20

replace (
	github.com/jxo-me/netx/api => ../api
	github.com/jxo-me/netx/core => ../core
	github.com/jxo-me/netx/gosocks4 => ../gosocks4
	github.com/jxo-me/netx/gosocks5 => ../gosocks5
	github.com/jxo-me/netx/plugin => ../plugin
	github.com/jxo-me/netx/relay => ../relay
	github.com/jxo-me/netx/tls-dissector => ../tls-dissector
)

require (
	github.com/alecthomas/units v0.0.0-20211218093645-b94a6e3cc137
	github.com/asaskevich/govalidator v0.0.0-20230301143203-a9d515a09cc2
	github.com/gin-contrib/cors v1.4.0
	github.com/gin-gonic/gin v1.9.0
	github.com/go-redis/redis/v8 v8.11.5
	github.com/gobwas/glob v0.2.3
	github.com/golang/snappy v0.0.4
	github.com/google/uuid v1.3.0
	github.com/gorilla/websocket v1.5.0
	github.com/jxo-me/netx/api v0.0.0-20230601103646-f32dd55b6fe3
	github.com/jxo-me/netx/core v0.0.0-20230531025546-78c9020abc9b
	github.com/jxo-me/netx/gosocks4 v0.0.1
	github.com/jxo-me/netx/gosocks5 v0.3.0
	github.com/jxo-me/netx/plugin v0.0.0-00010101000000-000000000000
	github.com/jxo-me/netx/relay v0.4.0
	github.com/jxo-me/netx/tls-dissector v0.0.1
	github.com/miekg/dns v1.1.54
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/pion/dtls/v2 v2.2.7
	github.com/pires/go-proxyproto v0.7.0
	github.com/prometheus/client_golang v1.15.1
	github.com/quic-go/quic-go v0.34.0
	github.com/rs/xid v1.5.0
	github.com/shadowsocks/go-shadowsocks2 v0.1.5
	github.com/shadowsocks/shadowsocks-go v0.0.0-20200409064450-3e585ff90601
	github.com/sirupsen/logrus v1.9.2
	github.com/songgao/water v0.0.0-20200317203138-2b4b6d7c09d8
	github.com/spf13/viper v1.15.0
	github.com/vishvananda/netlink v1.1.0
	github.com/xtaci/kcp-go/v5 v5.6.2
	github.com/xtaci/smux v1.5.24
	github.com/xtaci/tcpraw v1.2.25
	github.com/yl2chen/cidranger v1.0.2
	golang.org/x/crypto v0.9.0
	golang.org/x/net v0.10.0
	golang.org/x/sys v0.8.0
	golang.org/x/time v0.3.0
	golang.zx2c4.com/wireguard v0.0.0-20220703234212-c31a7b1ab478
	google.golang.org/grpc v1.55.0
	google.golang.org/protobuf v1.30.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/BurntSushi/toml v1.1.0 // indirect
	github.com/aead/chacha20 v0.0.0-20180709150244-8b13a72661da // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bytedance/sonic v1.8.0 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/chenzhuoyu/base64x v0.0.0-20221115062448-fe3a3abad311 // indirect
	github.com/clbanning/mxj/v2 v2.5.5 // indirect
	github.com/coreos/go-iptables v0.6.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/fatih/color v1.13.0 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.11.2 // indirect
	github.com/go-task/slim-sprig v0.0.0-20210107165309-348f09dbbbc0 // indirect
	github.com/goccy/go-json v0.10.0 // indirect
	github.com/gogf/gf/v2 v2.4.1 // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/gopacket v1.1.19 // indirect
	github.com/google/pprof v0.0.0-20210407192527-94a9f03dee38 // indirect
	github.com/grokify/html-strip-tags-go v0.0.1 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/jxo-me/gfbot v0.1.14 // indirect
	github.com/klauspost/cpuid/v2 v2.0.14 // indirect
	github.com/klauspost/reedsolomon v1.10.0 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/onsi/ginkgo/v2 v2.2.0 // indirect
	github.com/pelletier/go-toml/v2 v2.0.6 // indirect
	github.com/pion/logging v0.2.2 // indirect
	github.com/pion/transport/v2 v2.2.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/client_model v0.3.0 // indirect
	github.com/prometheus/common v0.42.0 // indirect
	github.com/prometheus/procfs v0.9.0 // indirect
	github.com/quic-go/qpack v0.4.0 // indirect
	github.com/quic-go/qtls-go1-19 v0.3.2 // indirect
	github.com/quic-go/qtls-go1-20 v0.2.2 // indirect
	github.com/riobard/go-bloom v0.0.0-20200614022211-cdc8013cb5b3 // indirect
	github.com/spf13/afero v1.9.3 // indirect
	github.com/spf13/cast v1.5.0 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/subosito/gotenv v1.4.2 // indirect
	github.com/templexxx/cpu v0.0.9 // indirect
	github.com/templexxx/xorsimd v0.4.1 // indirect
	github.com/tjfoc/gmsm v1.4.1 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.9 // indirect
	github.com/vishvananda/netns v0.0.0-20191106174202-0a2b9b5464df // indirect
	go.opentelemetry.io/otel v1.7.0 // indirect
	go.opentelemetry.io/otel/sdk v1.7.0 // indirect
	go.opentelemetry.io/otel/trace v1.7.0 // indirect
	golang.org/x/arch v0.0.0-20210923205945-b76863e36670 // indirect
	golang.org/x/exp v0.0.0-20221205204356-47842c84f3db // indirect
	golang.org/x/mod v0.8.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	golang.org/x/tools v0.6.0 // indirect
	golang.zx2c4.com/wintun v0.0.0-20230126152724-0fa3db229ce2 // indirect
	google.golang.org/genproto v0.0.0-20230306155012-7f2fa6fef1f4 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
)
