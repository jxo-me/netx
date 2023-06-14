package bot

import (
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/v2/util/gconv"
	telebot "github.com/jxo-me/gfbot"
	"github.com/jxo-me/gfbot/handlers"
	"github.com/jxo-me/netx/x/app"
	"github.com/jxo-me/netx/x/config"
	"github.com/jxo-me/netx/x/config/parsing"
)

const (
	HostAdd          = "hostAdd"
	HostUpdate       = "hostUpdate"
	HostsExampleJson = `
{
  "name": "hosts-0",
  "mappings": [
    {
      "ip": "127.0.0.1",
      "hostname": "localhost"
    },
    {
      "ip": "192.168.1.10",
      "hostname": "foo.mydomain.org",
      "aliases": [
        "foo"
      ]
    },
    {
      "ip": "192.168.1.13",
      "hostname": "bar.mydomain.org",
      "aliases": [
        "bar",
        "baz"
      ]
    }
  ]
}
`
)

func (h *hEvent) OnClickHosts(c telebot.IContext) error {
	var (
		msg string
	)

	cfg := config.Global()
	rowList := make([]telebot.Row, 0)
	btnList := make([]telebot.Btn, 0)
	selector := &telebot.ReplyMarkup{}
	for _, service := range cfg.Hosts {
		btnList = append(btnList, selector.Data(fmt.Sprintf("@%s", service.Name), "detailHosts", service.Name))
	}
	rowList = append(rowList, selector.Split(MaxCol, btnList)...)
	rowList = append(rowList, selector.Row(
		selector.Data("@添加主机映射器", "addHosts", "addHosts"),
		selector.Data("« 返回 服务列表", "backServices", "backServices"),
	))

	selector.Inline(
		rowList...,
	)
	msg = "Hosts List:\n"
	if c.Callback() != nil {
		return c.Edit(msg, &telebot.SendOptions{ReplyMarkup: selector, ParseMode: telebot.ModeMarkdownV2})
	}

	return c.Reply(msg, &telebot.SendOptions{ReplyMarkup: selector, ParseMode: telebot.ModeMarkdownV2})
}

func (h *hEvent) OnClickDetailHosts(c telebot.IContext) error {
	var (
		msg string
		str string
		err error
	)

	serviceName := c.Callback().Data
	cfg := config.Global()
	var srv *config.HostsConfig
	for _, service := range cfg.Hosts {
		if service.Name == serviceName {
			srv = service
		}
	}
	if srv != nil {
		str, err = ConvertJsonMsg(srv)
		if err != nil {
			return c.Reply("OnClickDetailHost ConvertJsonMsg err:", err.Error())
		}
		msg = fmt.Sprintf(CodeMsgTpl, msg, CodeStart, str, CodeEnd)
	}
	selector := &telebot.ReplyMarkup{}
	selector.Inline(
		selector.Row(
			selector.Data("@更新主机映射器", "updateHosts", serviceName),
			selector.Data("@删除主机映射器", "delHosts", serviceName),
		),
		selector.Row(
			selector.Data("« 返回 服务列表", "backServices", "backServices"),
		),
	)
	return c.Edit(msg, &telebot.SendOptions{ReplyMarkup: selector, ParseMode: telebot.ModeMarkdownV2})
}

func (h *hEvent) OnClickDelHosts(c telebot.IContext) error {
	serviceName := c.Callback().Data
	svc := app.Runtime.HostsRegistry().Get(serviceName)
	if svc == nil {
		return c.Send(ErrNotFound)
	}

	app.Runtime.HostsRegistry().Unregister(serviceName)

	_ = config.OnUpdate(func(c *config.Config) error {
		hostes := c.Hosts
		c.Hosts = nil
		for _, s := range hostes {
			if s.Name == serviceName {
				continue
			}
			c.Hosts = append(c.Hosts, s)
		}
		return nil
	})
	_ = c.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("%s 删除成功!", serviceName)})
	return h.OnClickHosts(c)
}

func AddHostsConversation(entry, cancel string) handlers.Conversation {
	return handlers.NewConversation(
		entry,
		telebot.HandlerFunc(startAddHostHandler), // 入口
		map[string][]telebot.IHandler{
			HostAdd: {telebot.HandlerFunc(addHostHandler)},
		}, // states状态
		&handlers.ConversationOpts{
			ExitName:     cancel,
			ExitHandler:  telebot.HandlerFunc(cancelHostHandler),
			AllowReEntry: true,
		}, // options
	)
}

func UpdateHostsConversation(entry, cancel string) handlers.Conversation {
	return handlers.NewConversation(
		entry,
		telebot.HandlerFunc(startUpdateHostHandler), // 入口
		map[string][]telebot.IHandler{
			HostUpdate: {telebot.HandlerFunc(updateHostHandler)},
		}, // states状态
		&handlers.ConversationOpts{
			ExitName:     cancel,
			ExitHandler:  telebot.HandlerFunc(cancelHostHandler),
			AllowReEntry: true,
		}, // options
	)
}

func startAddHostHandler(ctx telebot.IContext) error {
	example := fmt.Sprintf(CodeTpl, CodeStart, HostsExampleJson, CodeEnd)
	err := ctx.Send(fmt.Sprintf("请输入 主机映射器 配置?\nExample：%s\n您可以随时键入 /cancel 来取消该操作。", example), &telebot.SendOptions{ParseMode: telebot.ModeMarkdownV2})
	if err != nil {
		return fmt.Errorf("failed to send start message: %w", err)
	}

	// 设置当前用户下一个入口
	return handlers.NextConversationState(HostAdd)
}

func addHostHandler(ctx telebot.IContext) error {
	var (
		data config.HostsConfig
	)

	str := ctx.Message().Text
	err := json.Unmarshal([]byte(str), &data)
	if err != nil {
		return ctx.Reply("addHostHandler json.Unmarshal error:", err.Error())
	}
	v := parsing.ParseHosts(&data)
	if err = app.Runtime.HostsRegistry().Register(data.Name, v); err != nil {
		return ctx.Reply(ErrDup)
	}

	_ = config.OnUpdate(func(c *config.Config) error {
		c.Hosts = append(c.Hosts, &data)
		return nil
	})

	_ = ctx.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("%s 添加成功!", data.Name)})
	_ = Event.OnClickHosts(ctx)

	return handlers.EndConversation()
}

func startUpdateHostHandler(ctx telebot.IContext) error {
	srvName := ctx.Callback().Data
	err := ctx.Bot().Store().UpdateData(ctx, HostUpdate, srvName)
	if err != nil {
		return fmt.Errorf("failed UpdateData message: %w", err)
	}
	example := fmt.Sprintf(CodeTpl, CodeStart, HostsExampleJson, CodeEnd)
	err = ctx.Send(fmt.Sprintf("请输入 主机映射器 配置?\nExample：%s\n您可以随时键入 /cancel 来取消该操作。", example), &telebot.SendOptions{ParseMode: telebot.ModeMarkdownV2})
	if err != nil {
		return fmt.Errorf("failed to send start message: %w", err)
	}

	// 设置当前用户下一个入口
	return handlers.NextConversationState(HostUpdate)
}

func updateHostHandler(ctx telebot.IContext) error {
	var (
		data config.HostsConfig
	)
	state, err := ctx.Bot().Store().Get(ctx)
	if err != nil {
		return ctx.Reply("updateHostHandler Store().Get error:", err.Error())
	}
	srvName := ""
	if sn, ok := state.Data[HostUpdate]; ok {
		srvName = gconv.String(sn)
	}
	fmt.Println(fmt.Sprintf("srvName :%s", srvName))
	if srvName == "" {
		return ctx.Reply(ErrInvalid)
	}
	str := ctx.Message().Text
	err = json.Unmarshal([]byte(str), &data)
	if err != nil {
		return ctx.Reply("updateHostHandler json.Unmarshal error:", err.Error())
	}

	if !app.Runtime.HostsRegistry().IsRegistered(srvName) {
		return ctx.Reply(ErrNotFound)
	}

	data.Name = srvName

	v := parsing.ParseHosts(&data)

	app.Runtime.HostsRegistry().Unregister(srvName)

	if err = app.Runtime.HostsRegistry().Register(srvName, v); err != nil {
		return ctx.Reply(ErrDup)
	}

	_ = config.OnUpdate(func(c *config.Config) error {
		for i := range c.Hosts {
			if c.Hosts[i].Name == srvName {
				c.Hosts[i] = &data
				break
			}
		}
		return nil
	})

	_ = ctx.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("%s 更新成功!", data.Name)})
	_ = Event.OnClickHosts(ctx)

	return handlers.EndConversation()
}

func cancelHostHandler(ctx telebot.IContext) error {
	err := ctx.Reply("添加 主机映射器 已被取消。 还有什么我可以为你做的吗？", &telebot.SendOptions{})
	if err != nil {
		return fmt.Errorf("failed to send cancelHandler message: %w", err)
	}
	//return handlers.EndConversation()
	return nil
}
