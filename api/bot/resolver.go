package bot

import (
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/v2/util/gconv"
	telebot "github.com/jxo-me/gfbot"
	"github.com/jxo-me/gfbot/handlers"
	"github.com/jxo-me/netx/x/app"
	"github.com/jxo-me/netx/x/config"
	parser "github.com/jxo-me/netx/x/config/parsing/resolver"
)

const (
	ResolverAdd         = "resolverAdd"
	ResolverUpdate      = "resolverUpdate"
	ResolverExampleJson = `
{
  "name": "resolver-0",
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

func (h *hEvent) OnClickResolvers(c telebot.IContext) error {
	var (
		msg string
	)

	cfg := config.Global()
	rowList := make([]telebot.Row, 0)
	btnList := make([]telebot.Btn, 0)
	selector := &telebot.ReplyMarkup{}
	for _, service := range cfg.Resolvers {
		btnList = append(btnList, selector.Data(fmt.Sprintf("@%s", service.Name), "detailResolver", service.Name))
	}
	rowList = append(rowList, selector.Split(MaxCol, btnList)...)
	rowList = append(rowList, selector.Row(
		selector.Data("@添加域名解析器", "addResolver", "addResolver"),
		selector.Data("« 返回 服务列表", "backServices", "backServices"),
	))

	selector.Inline(
		rowList...,
	)
	msg = "Resolvers List:\n"
	if c.Callback() != nil {
		return c.Edit(msg, &telebot.SendOptions{ReplyMarkup: selector, ParseMode: telebot.ModeMarkdownV2})
	}

	return c.Reply(msg, &telebot.SendOptions{ReplyMarkup: selector, ParseMode: telebot.ModeMarkdownV2})
}

func (h *hEvent) OnClickDetailResolver(c telebot.IContext) error {
	var (
		msg string
		str string
		err error
	)

	serviceName := c.Callback().Data
	cfg := config.Global()
	var srv *config.ResolverConfig
	for _, service := range cfg.Resolvers {
		if service.Name == serviceName {
			srv = service
		}
	}
	if srv != nil {
		str, err = ConvertJsonMsg(srv)
		if err != nil {
			return c.Reply("OnClickDetailResolver ConvertJsonMsg err:", err.Error())
		}
		msg = fmt.Sprintf(CodeMsgTpl, msg, CodeStart, str, CodeEnd)
	}
	selector := &telebot.ReplyMarkup{}
	selector.Inline(
		selector.Row(
			selector.Data("@更新域名解析器", "updateResolver", serviceName),
			selector.Data("@删除域名解析器", "delResolver", serviceName),
		),
		selector.Row(
			selector.Data("« 返回 服务列表", "backServices", "backServices"),
		),
	)
	return c.Edit(msg, &telebot.SendOptions{ReplyMarkup: selector, ParseMode: telebot.ModeMarkdownV2})
}

func (h *hEvent) OnClickDelResolver(c telebot.IContext) error {
	serviceName := c.Callback().Data
	svc := app.Runtime.ResolverRegistry().Get(serviceName)
	if svc == nil {
		return c.Send(ErrNotFound)
	}

	app.Runtime.ResolverRegistry().Unregister(serviceName)

	_ = config.OnUpdate(func(c *config.Config) error {
		resolveres := c.Resolvers
		c.Resolvers = nil
		for _, s := range resolveres {
			if s.Name == serviceName {
				continue
			}
			c.Resolvers = append(c.Resolvers, s)
		}
		return nil
	})
	_ = c.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("%s 删除成功!", serviceName)})
	return h.OnClickResolvers(c)
}

func AddResolverConversation(entry, cancel string) handlers.Conversation {
	return handlers.NewConversation(
		entry,
		telebot.HandlerFunc(startAddResolverHandler), // 入口
		map[string][]telebot.IHandler{
			ResolverAdd: {telebot.HandlerFunc(addResolverHandler)},
		}, // states状态
		&handlers.ConversationOpts{
			ExitName:     cancel,
			ExitHandler:  telebot.HandlerFunc(cancelResolverHandler),
			AllowReEntry: true,
		}, // options
	)
}

func UpdateResolverConversation(entry, cancel string) handlers.Conversation {
	return handlers.NewConversation(
		entry,
		telebot.HandlerFunc(startUpdateResolverHandler), // 入口
		map[string][]telebot.IHandler{
			ResolverUpdate: {telebot.HandlerFunc(updateResolverHandler)},
		}, // states状态
		&handlers.ConversationOpts{
			ExitName:     cancel,
			ExitHandler:  telebot.HandlerFunc(cancelResolverHandler),
			AllowReEntry: true,
		}, // options
	)
}

func startAddResolverHandler(ctx telebot.IContext) error {
	example := fmt.Sprintf(CodeTpl, CodeStart, ResolverExampleJson, CodeEnd)
	err := ctx.Send(fmt.Sprintf("请输入 域名解析器 配置?\nExample：%s\n您可以随时键入 /cancel 来取消该操作。", example), &telebot.SendOptions{ParseMode: telebot.ModeMarkdownV2})
	if err != nil {
		return fmt.Errorf("failed to send start message: %w", err)
	}

	// 设置当前用户下一个入口
	return handlers.NextConversationState(ResolverAdd)
}

func addResolverHandler(ctx telebot.IContext) error {
	var (
		data config.ResolverConfig
	)

	str := ctx.Message().Text
	err := json.Unmarshal([]byte(str), &data)
	if err != nil {
		return ctx.Reply("addResolverHandler json.Unmarshal error:", err.Error())
	}
	v, err := parser.ParseResolver(&data)
	if err != nil {
		return ctx.Reply(ErrCreate)
	}
	if err = app.Runtime.ResolverRegistry().Register(data.Name, v); err != nil {
		return ctx.Reply(ErrDup)
	}

	_ = config.OnUpdate(func(c *config.Config) error {
		c.Resolvers = append(c.Resolvers, &data)
		return nil
	})

	_ = ctx.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("%s 添加成功!", data.Name)})
	_ = Event.OnClickResolvers(ctx)

	return handlers.EndConversation()
}

func startUpdateResolverHandler(ctx telebot.IContext) error {
	srvName := ctx.Callback().Data
	err := ctx.Bot().Store().UpdateData(ctx, ResolverUpdate, srvName)
	if err != nil {
		return fmt.Errorf("failed UpdateData message: %w", err)
	}
	example := fmt.Sprintf(CodeTpl, CodeStart, ResolverExampleJson, CodeEnd)
	err = ctx.Send(fmt.Sprintf("请输入 域名解析器 配置?\nExample：%s\n您可以随时键入 /cancel 来取消该操作。", example), &telebot.SendOptions{ParseMode: telebot.ModeMarkdownV2})
	if err != nil {
		return fmt.Errorf("failed to send start message: %w", err)
	}

	// 设置当前用户下一个入口
	return handlers.NextConversationState(ResolverUpdate)
}

func updateResolverHandler(ctx telebot.IContext) error {
	var (
		data config.ResolverConfig
	)
	state, err := ctx.Bot().Store().Get(ctx)
	if err != nil {
		return ctx.Reply("updateResolverHandler Store().Get error:", err.Error())
	}
	srvName := ""
	if sn, ok := state.Data[ResolverUpdate]; ok {
		srvName = gconv.String(sn)
	}
	fmt.Println(fmt.Sprintf("srvName :%s", srvName))
	if srvName == "" {
		return ctx.Reply(ErrInvalid)
	}
	str := ctx.Message().Text
	err = json.Unmarshal([]byte(str), &data)
	if err != nil {
		return ctx.Reply("updateResolverHandler json.Unmarshal error:", err.Error())
	}

	if !app.Runtime.ResolverRegistry().IsRegistered(srvName) {
		return ctx.Reply(ErrNotFound)
	}

	data.Name = srvName

	v, err := parser.ParseResolver(&data)
	if err != nil {
		return ctx.Reply(ErrCreate)
	}
	app.Runtime.ResolverRegistry().Unregister(srvName)

	if err = app.Runtime.ResolverRegistry().Register(srvName, v); err != nil {
		return ctx.Reply(ErrDup)
	}

	_ = config.OnUpdate(func(c *config.Config) error {
		for i := range c.Resolvers {
			if c.Resolvers[i].Name == srvName {
				c.Resolvers[i] = &data
				break
			}
		}
		return nil
	})

	_ = ctx.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("%s 更新成功!", data.Name)})
	_ = Event.OnClickResolvers(ctx)

	return handlers.EndConversation()
}

func cancelResolverHandler(ctx telebot.IContext) error {
	err := ctx.Reply("添加 域名解析器 已被取消。 还有什么我可以为你做的吗？", &telebot.SendOptions{})
	if err != nil {
		return fmt.Errorf("failed to send cancelHandler message: %w", err)
	}
	//return handlers.EndConversation()
	return nil
}
