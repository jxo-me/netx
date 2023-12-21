package boot

import (
	"fmt"
	"github.com/jxo-me/netx/core/app"
	"github.com/jxo-me/netx/core/connector"
	"github.com/jxo-me/netx/core/dialer"
	"github.com/jxo-me/netx/core/handler"
	"github.com/jxo-me/netx/core/listener"
)

var (
	insBoot = Boot{}
)

func Boots(a app.IRuntime) *Boot {
	insBoot.App = a
	insBoot.Connectors = Connectors
	insBoot.Dialers = Dialers
	insBoot.Handlers = Handlers
	insBoot.Listeners = Listeners
	// Register connectors
	err := insBoot.InitConnector()
	if err != nil {
		panic(fmt.Sprintf("InitConnector error: %s", err.Error()))
		return nil
	}
	// Register dialers
	err = insBoot.InitDialer()
	if err != nil {
		panic(fmt.Sprintf("InitDialer error: %s", err.Error()))
		return nil
	}
	// Register handlers
	err = insBoot.InitHandler()
	if err != nil {
		panic(fmt.Sprintf("InitHandler error: %s", err.Error()))
		return nil
	}
	// Register listeners
	err = insBoot.InitListener()
	if err != nil {
		panic(fmt.Sprintf("InitListener error: %s", err.Error()))
		return nil
	}

	return &insBoot
}

type Boot struct {
	App        app.IRuntime
	Connectors map[string]connector.NewConnector
	Dialers    map[string]dialer.NewDialer
	Handlers   map[string]handler.NewHandler
	Listeners  map[string]listener.NewListener
}

func (b *Boot) InitConnector() (err error) {
	// connector
	for name, connector := range b.Connectors {
		//fmt.Println("Register Connector type:", name)
		err = b.App.ConnectorRegistry().Register(name, connector)
		if err != nil {
			return
		}
	}
	return err
}

func (b *Boot) InitDialer() (err error) {
	// dialer
	for name, dialer := range b.Dialers {
		//fmt.Println("Register Dialer type:", name)
		err = b.App.DialerRegistry().Register(name, dialer)
		if err != nil {
			return err
		}
	}
	return err
}

func (b *Boot) InitHandler() (err error) {
	// handler
	for name, handle := range b.Handlers {
		//fmt.Println("Register Handler type:", name)
		err = b.App.HandlerRegistry().Register(name, handle)
		if err != nil {
			return err
		}
	}
	return err
}

func (b *Boot) InitListener() (err error) {
	// listener
	for name, listener := range b.Listeners {
		//fmt.Println("Register Listener type:", name)
		err = b.App.ListenerRegistry().Register(name, listener)
		if err != nil {
			return err
		}
	}
	return err
}
