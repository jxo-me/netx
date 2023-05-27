package service

import (
	"net"
)

type IService interface {
	Serve() error
	Addr() net.Addr
	Close() error
}
