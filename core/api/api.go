package api

import (
	"github.com/gogf/gf/v2/net/ghttp"
	telebot "github.com/jxo-me/gfbot"
	"github.com/jxo-me/netx/core/service"
	"net"
)

type TGBot struct {
	Bot    *telebot.Bot `json:"bot"`
	Domain string       `json:"domain"`
	Token  string       `json:"token"`
}

type Server struct {
	Srv      *ghttp.Server
	Listener net.Listener
	Bot      *TGBot
}

type IApi interface {
	service.IService
	TGBot() *TGBot
}
