package recorder

import (
	"context"

	"github.com/jxo-me/netx/core/logger"
	"github.com/jxo-me/netx/core/recorder"
	"github.com/jxo-me/netx/plugin/recorder/proto"
	xlogger "github.com/jxo-me/netx/x/logger"
	"google.golang.org/grpc"
)

type pluginOptions struct {
	client *grpc.ClientConn
	logger logger.ILogger
}

type PluginOption func(opts *pluginOptions)

func PluginConnOption(c *grpc.ClientConn) PluginOption {
	return func(opts *pluginOptions) {
		opts.client = c
	}
}

func LoggerOption(logger logger.ILogger) PluginOption {
	return func(opts *pluginOptions) {
		opts.logger = logger
	}
}

type pluginRecorder struct {
	client  proto.RecorderClient
	options pluginOptions
}

// NewPluginRecorder creates a plugin recorder.
func NewPluginRecorder(opts ...PluginOption) recorder.Recorder {
	var options pluginOptions
	for _, opt := range opts {
		opt(&options)
	}
	if options.logger == nil {
		options.logger = xlogger.Nop()
	}

	p := &pluginRecorder{
		options: options,
	}
	if options.client != nil {
		p.client = proto.NewRecorderClient(options.client)
	}
	return p
}

func (p *pluginRecorder) Record(ctx context.Context, b []byte) error {
	if p.client == nil {
		return nil
	}

	_, err := p.client.Record(context.Background(),
		&proto.RecordRequest{
			Data: b,
		})
	if err != nil {
		p.options.logger.Error(err)
		return err
	}
	return nil
}

func (p *pluginRecorder) Close() error {
	if p.options.client != nil {
		return p.options.client.Close()
	}
	return nil
}
