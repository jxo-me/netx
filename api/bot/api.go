package bot

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	telebot "github.com/jxo-me/gfbot"
	"github.com/jxo-me/netx/x/app"
)

var (
	Bot = hBot{}
)

type hBot struct {
	Path     string `json:"path" yaml:"path"`
	BasePath string `json:"basePath" yaml:"base_path"`
}

type HookReq struct {
	g.Meta `path:"/v1/bot/hook" method:"post" tags:"TGBot" summary:"Bot WebHook API"`
	//Name   string `v:"required" json:"nameHandler" description:"名称"`
}

type HookRes struct {
}

func (h *hBot) Hook(ctx context.Context, req *HookReq) (result *HookRes, err error) {
	r := g.RequestFromCtx(ctx)
	if app.ApiSrv != nil {
		app.ApiSrv.TGBot().Bot.Hook().Handler(r.Response.ResponseWriter, r.Request)
	}

	return &HookRes{}, nil
}

type ApiRes struct {
	Resp *telebot.Message
}

type BotMsg struct {
	Id      int64  `json:"id"        description:"Id"`
	Message string `json:"message"     description:"消息"`
}

type BotSendMsgReq struct {
	g.Meta `path:"/api/v1/bot/message" method:"post" tags:"TGBot" summary:"Telegram Bot Send message"`
	BotMsg
}

func (h *hBot) Send(ctx context.Context, req *BotSendMsgReq) (result *ApiRes, err error) {
	b := app.ApiSrv.TGBot()
	u := telebot.User{ID: req.Id}
	send, err := b.Bot.Send(&u, fmt.Sprintf("%s%s", "Mickey", "121.7.103.209"))
	if err != nil {
		return nil, err
	}
	glog.Debug(ctx, "Send Message to User ...")
	g.Dump(send)
	return &ApiRes{Resp: send}, nil
}
