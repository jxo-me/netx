package bot

import (
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/v2/util/gconv"
	telebot "github.com/jxo-me/gfbot"
	"github.com/jxo-me/gfbot/handlers"
	"github.com/jxo-me/netx/x/app"
	"github.com/jxo-me/netx/x/config"
	parser "github.com/jxo-me/netx/x/config/parsing/router"
)

const (
	RouterAdd         = "RouterAdd"
	RouterUpdate      = "RouterUpdate"
	RouterExampleJson = `
{
  "name": "Router-0",
  "nameservers": [
    {
      "addr": "udp://8.8.8.8:53",
      "chain": "chain-0",
      "prefer": "ipv4",
      "clientIP": "1.2.3.4",
      "ttl": 60,
      "timeout": 30
    },
    {
      "addr": "tcp://1.1.1.1:53"
    },
    {
      "addr": "tls://1.1.1.1:853"
    },
    {
      "addr": "https://1.0.0.1/dns-query",
      "hostname": "cloudflare-dns.com"
    }
  ]
}
`
)

func (h *hEvent) OnClickRouters(c telebot.IContext) error {
	var (
		msg string
	)

	cfg := config.Global()
	rowList := make([]telebot.Row, 0)
	btnList := make([]telebot.Btn, 0)
	selector := &telebot.ReplyMarkup{}
	for _, service := range cfg.Routers {
		btnList = append(btnList, selector.Data(fmt.Sprintf("@%s", service.Name), "detailRouter", service.Name))
	}
	rowList = append(rowList, selector.Split(MaxCol, btnList)...)
	rowList = append(rowList, selector.Row(
		selector.Data("@添加路由", "addRouter", "addRouter"),
		selector.Data("« 返回 服务列表", "backServices", "backServices"),
	))

	selector.Inline(
		rowList...,
	)
	msg = "Routers List:\n"
	if c.Callback() != nil {
		return c.Edit(msg, &telebot.SendOptions{ReplyMarkup: selector, ParseMode: telebot.ModeMarkdownV2})
	}

	return c.Reply(msg, &telebot.SendOptions{ReplyMarkup: selector, ParseMode: telebot.ModeMarkdownV2})
}

func (h *hEvent) OnClickDetailRouter(c telebot.IContext) error {
	var (
		msg string
		str string
		err error
	)

	serviceName := c.Callback().Data
	cfg := config.Global()
	var srv *config.RouterConfig
	for _, service := range cfg.Routers {
		if service.Name == serviceName {
			srv = service
		}
	}
	if srv != nil {
		str, err = ConvertJsonMsg(srv)
		if err != nil {
			return c.Reply("OnClickDetailRouter ConvertJsonMsg err:", err.Error())
		}
		msg = fmt.Sprintf(CodeMsgTpl, msg, CodeStart, str, CodeEnd)
	}
	selector := &telebot.ReplyMarkup{}
	selector.Inline(
		selector.Row(
			selector.Data("@更新域名解析器", "updateRouter", serviceName),
			selector.Data("@删除域名解析器", "delRouter", serviceName),
		),
		selector.Row(
			selector.Data("« 返回 服务列表", "backServices", "backServices"),
		),
	)
	return c.Edit(msg, &telebot.SendOptions{ReplyMarkup: selector, ParseMode: telebot.ModeMarkdownV2})
}

func (h *hEvent) OnClickDelRouter(c telebot.IContext) error {
	serviceName := c.Callback().Data
	svc := app.Runtime.RouterRegistry().Get(serviceName)
	if svc == nil {
		return c.Send(ErrNotFound)
	}

	app.Runtime.RouterRegistry().Unregister(serviceName)

	_ = config.OnUpdate(func(c *config.Config) error {
		Routeres := c.Routers
		c.Routers = nil
		for _, s := range Routeres {
			if s.Name == serviceName {
				continue
			}
			c.Routers = append(c.Routers, s)
		}
		return nil
	})
	_ = c.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("%s 删除成功!", serviceName)})
	return h.OnClickRouters(c)
}

func AddRouterConversation(entry, cancel string) handlers.Conversation {
	return handlers.NewConversation(
		entry,
		telebot.HandlerFunc(startAddRouterHandler), // 入口
		map[string][]telebot.IHandler{
			RouterAdd: {telebot.HandlerFunc(addRouterHandler)},
		}, // states状态
		&handlers.ConversationOpts{
			ExitName:     cancel,
			ExitHandler:  telebot.HandlerFunc(cancelRouterHandler),
			AllowReEntry: true,
		}, // options
	)
}

func UpdateRouterConversation(entry, cancel string) handlers.Conversation {
	return handlers.NewConversation(
		entry,
		telebot.HandlerFunc(startUpdateRouterHandler), // 入口
		map[string][]telebot.IHandler{
			RouterUpdate: {telebot.HandlerFunc(updateRouterHandler)},
		}, // states状态
		&handlers.ConversationOpts{
			ExitName:     cancel,
			ExitHandler:  telebot.HandlerFunc(cancelRouterHandler),
			AllowReEntry: true,
		}, // options
	)
}

func startAddRouterHandler(ctx telebot.IContext) error {
	example := fmt.Sprintf(CodeTpl, CodeStart, RouterExampleJson, CodeEnd)
	err := ctx.Send(fmt.Sprintf("请输入 域名解析器 配置?\nExample：%s\n您可以随时键入 /cancel 来取消该操作。", example), &telebot.SendOptions{ParseMode: telebot.ModeMarkdownV2})
	if err != nil {
		return fmt.Errorf("failed to send start message: %w", err)
	}

	// 设置当前用户下一个入口
	return handlers.NextConversationState(RouterAdd)
}

func addRouterHandler(ctx telebot.IContext) error {
	var (
		data config.RouterConfig
	)

	str := ctx.Message().Text
	err := json.Unmarshal([]byte(str), &data)
	if err != nil {
		return ctx.Reply("addRouterHandler json.Unmarshal error:", err.Error())
	}
	v := parser.ParseRouter(&data)
	if err = app.Runtime.RouterRegistry().Register(data.Name, v); err != nil {
		return ctx.Reply(ErrDup)
	}

	_ = config.OnUpdate(func(c *config.Config) error {
		c.Routers = append(c.Routers, &data)
		return nil
	})

	_ = ctx.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("%s 添加成功!", data.Name)})
	_ = Event.OnClickRouters(ctx)

	return handlers.EndConversation()
}

func startUpdateRouterHandler(ctx telebot.IContext) error {
	srvName := ctx.Callback().Data
	err := ctx.Bot().Store().UpdateData(ctx, RouterUpdate, srvName)
	if err != nil {
		return fmt.Errorf("failed UpdateData message: %w", err)
	}
	example := fmt.Sprintf(CodeTpl, CodeStart, RouterExampleJson, CodeEnd)
	err = ctx.Send(fmt.Sprintf("请输入 域名解析器 配置?\nExample：%s\n您可以随时键入 /cancel 来取消该操作。", example), &telebot.SendOptions{ParseMode: telebot.ModeMarkdownV2})
	if err != nil {
		return fmt.Errorf("failed to send start message: %w", err)
	}

	// 设置当前用户下一个入口
	return handlers.NextConversationState(RouterUpdate)
}

func updateRouterHandler(ctx telebot.IContext) error {
	var (
		data config.RouterConfig
	)
	state, err := ctx.Bot().Store().Get(ctx)
	if err != nil {
		return ctx.Reply("updateRouterHandler Store().Get error:", err.Error())
	}
	srvName := ""
	if sn, ok := state.Data[RouterUpdate]; ok {
		srvName = gconv.String(sn)
	}
	fmt.Println(fmt.Sprintf("srvName :%s", srvName))
	if srvName == "" {
		return ctx.Reply(ErrInvalid)
	}
	str := ctx.Message().Text
	err = json.Unmarshal([]byte(str), &data)
	if err != nil {
		return ctx.Reply("updateRouterHandler json.Unmarshal error:", err.Error())
	}

	if !app.Runtime.RouterRegistry().IsRegistered(srvName) {
		return ctx.Reply(ErrNotFound)
	}

	data.Name = srvName

	v := parser.ParseRouter(&data)
	app.Runtime.RouterRegistry().Unregister(srvName)

	if err = app.Runtime.RouterRegistry().Register(srvName, v); err != nil {
		return ctx.Reply(ErrDup)
	}

	_ = config.OnUpdate(func(c *config.Config) error {
		for i := range c.Routers {
			if c.Routers[i].Name == srvName {
				c.Routers[i] = &data
				break
			}
		}
		return nil
	})

	_ = ctx.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("%s 更新成功!", data.Name)})
	_ = Event.OnClickRouters(ctx)

	return handlers.EndConversation()
}

func cancelRouterHandler(ctx telebot.IContext) error {
	err := ctx.Reply("添加 域名解析器 已被取消。 还有什么我可以为你做的吗？", &telebot.SendOptions{})
	if err != nil {
		return fmt.Errorf("failed to send cancelHandler message: %w", err)
	}
	//return handlers.EndConversation()
	return nil
}
