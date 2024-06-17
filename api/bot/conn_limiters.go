package bot

import (
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/v2/util/gconv"
	telebot "github.com/jxo-me/gfbot"
	"github.com/jxo-me/gfbot/handlers"
	"github.com/jxo-me/netx/x/app"
	"github.com/jxo-me/netx/x/config"
	parser "github.com/jxo-me/netx/x/config/parsing/limiter"
)

const (
	ConnLimiterAdd      = "connLimiterAdd"
	ConnLimiterUpdate   = "connLimiterUpdate"
	CLimiterExampleJson = `
{
  "name": "climiter-0",
  "limits": [
    "$ 1000",
    "$$ 100",
    "192.168.1.1  200",
    "192.168.0.0/16  50"
  ]
}
`
)

func (h *hEvent) OnClickConnLimiters(c telebot.Context) error {
	var (
		msg string
	)

	cfg := config.Global()
	rowList := make([]telebot.Row, 0)
	btnList := make([]telebot.Btn, 0)
	selector := &telebot.ReplyMarkup{}
	for _, service := range cfg.CLimiters {
		btnList = append(btnList, selector.Data(fmt.Sprintf("@%s", service.Name), "detailConnLimiter", service.Name))
	}
	rowList = append(rowList, selector.Split(MaxCol, btnList)...)
	rowList = append(rowList, selector.Row(
		selector.Data("@添加并发连接数限制器", "addConnLimiter", "addConnLimiter"),
		selector.Data("« 返回 服务列表", "backServices", "backServices"),
	))

	selector.Inline(
		rowList...,
	)
	msg = "ConnLimiters List:\n"
	if c.Callback() != nil {
		return c.Edit(msg, &telebot.SendOptions{ReplyMarkup: selector, ParseMode: telebot.ModeMarkdownV2})
	}

	return c.Reply(msg, &telebot.SendOptions{ReplyMarkup: selector, ParseMode: telebot.ModeMarkdownV2})
}

func (h *hEvent) OnClickDetailConnLimiter(c telebot.Context) error {
	var (
		msg string
		str string
		err error
	)

	serviceName := c.Callback().Data
	cfg := config.Global()
	var srv *config.LimiterConfig
	for _, service := range cfg.CLimiters {
		if service.Name == serviceName {
			srv = service
		}
	}
	if srv != nil {
		str, err = ConvertJsonMsg(srv)
		if err != nil {
			return c.Reply("OnClickDetailConnLimiter ConvertJsonMsg err:", err.Error())
		}
		msg = fmt.Sprintf(CodeMsgTpl, msg, CodeStart, str, CodeEnd)
	}
	selector := &telebot.ReplyMarkup{}
	selector.Inline(
		selector.Row(
			selector.Data("@更新并发连接数限制器", "updateConnLimiter", serviceName),
			selector.Data("@删除并发连接数限制器", "delConnLimiter", serviceName),
		),
		selector.Row(
			selector.Data("« 返回 服务列表", "backServices", "backServices"),
		),
	)
	return c.Edit(msg, &telebot.SendOptions{ReplyMarkup: selector, ParseMode: telebot.ModeMarkdownV2})
}

func (h *hEvent) OnClickDelConnLimiter(c telebot.Context) error {
	serviceName := c.Callback().Data
	svc := app.Runtime.ConnLimiterRegistry().Get(serviceName)
	if svc == nil {
		return c.Send(ErrNotFound)
	}

	app.Runtime.ConnLimiterRegistry().Unregister(serviceName)

	_ = config.OnUpdate(func(c *config.Config) error {
		connLimiteres := c.CLimiters
		c.CLimiters = nil
		for _, s := range connLimiteres {
			if s.Name == serviceName {
				continue
			}
			c.CLimiters = append(c.CLimiters, s)
		}
		return nil
	})
	_ = c.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("%s 删除成功!", serviceName)})
	return h.OnClickConnLimiters(c)
}

func AddConnLimiterConversation(entry, cancel string) handlers.Conversation {
	return handlers.NewConversation(
		entry,
		telebot.HandlerFunc(startAddConnLimiterHandler), // 入口
		map[string][]telebot.IHandler{
			ConnLimiterAdd: {telebot.HandlerFunc(addConnLimiterHandler)},
		}, // states状态
		&handlers.ConversationOpts{
			ExitName:     cancel,
			ExitHandler:  telebot.HandlerFunc(cancelConnLimiterHandler),
			AllowReEntry: true,
		}, // options
	)
}

func UpdateConnLimiterConversation(entry, cancel string) handlers.Conversation {
	return handlers.NewConversation(
		entry,
		telebot.HandlerFunc(startUpdateConnLimiterHandler), // 入口
		map[string][]telebot.IHandler{
			ConnLimiterUpdate: {telebot.HandlerFunc(updateConnLimiterHandler)},
		}, // states状态
		&handlers.ConversationOpts{
			ExitName:     cancel,
			ExitHandler:  telebot.HandlerFunc(cancelConnLimiterHandler),
			AllowReEntry: true,
		}, // options
	)
}

func startAddConnLimiterHandler(ctx telebot.Context) error {
	example := fmt.Sprintf(CodeTpl, CodeStart, CLimiterExampleJson, CodeEnd)
	err := ctx.Send(fmt.Sprintf("请输入 并发限制器 配置?\nExample：%s\n您可以随时键入 /cancel 来取消该操作。", example), &telebot.SendOptions{ParseMode: telebot.ModeMarkdownV2})
	if err != nil {
		return fmt.Errorf("failed to send start message: %w", err)
	}

	// 设置当前用户下一个入口
	return handlers.NextConversationState(ConnLimiterAdd)
}

func addConnLimiterHandler(ctx telebot.Context) error {
	var (
		data config.LimiterConfig
	)

	str := ctx.Message().Text
	err := json.Unmarshal([]byte(str), &data)
	if err != nil {
		return ctx.Reply("addConnLimiterHandler json.Unmarshal error:", err.Error())
	}
	v := parser.ParseConnLimiter(&data)
	if err = app.Runtime.ConnLimiterRegistry().Register(data.Name, v); err != nil {
		return ctx.Reply(ErrDup)
	}

	_ = config.OnUpdate(func(c *config.Config) error {
		c.CLimiters = append(c.CLimiters, &data)
		return nil
	})

	_ = ctx.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("%s 添加成功!", data.Name)})
	_ = Event.OnClickConnLimiters(ctx)

	return handlers.EndConversation()
}

func startUpdateConnLimiterHandler(ctx telebot.Context) error {
	srvName := ctx.Callback().Data
	err := ctx.Bot().Store().UpdateData(ctx, ConnLimiterUpdate, srvName)
	if err != nil {
		return fmt.Errorf("failed UpdateData message: %w", err)
	}
	example := fmt.Sprintf(CodeTpl, CodeStart, CLimiterExampleJson, CodeEnd)
	err = ctx.Send(fmt.Sprintf("请输入 并发限制器 配置?\nExample：%s\n您可以随时键入 /cancel 来取消该操作。", example), &telebot.SendOptions{ParseMode: telebot.ModeMarkdownV2})
	if err != nil {
		return fmt.Errorf("failed to send start message: %w", err)
	}

	// 设置当前用户下一个入口
	return handlers.NextConversationState(ConnLimiterUpdate)
}

func updateConnLimiterHandler(ctx telebot.Context) error {
	var (
		data config.LimiterConfig
	)
	state, err := ctx.Bot().Store().Get(ctx)
	if err != nil {
		return ctx.Reply("updateConnLimiterHandler Store().Get error:", err.Error())
	}
	srvName := ""
	if sn, ok := state.Data[ConnLimiterUpdate]; ok {
		srvName = gconv.String(sn)
	}
	fmt.Println(fmt.Sprintf("srvName :%s", srvName))
	if srvName == "" {
		return ctx.Reply(ErrInvalid)
	}
	str := ctx.Message().Text
	err = json.Unmarshal([]byte(str), &data)
	if err != nil {
		return ctx.Reply("updateConnLimiterHandler json.Unmarshal error:", err.Error())
	}

	if !app.Runtime.ConnLimiterRegistry().IsRegistered(srvName) {
		return ctx.Reply(ErrNotFound)
	}

	data.Name = srvName

	v := parser.ParseConnLimiter(&data)

	app.Runtime.ConnLimiterRegistry().Unregister(srvName)

	if err = app.Runtime.ConnLimiterRegistry().Register(srvName, v); err != nil {
		return ctx.Reply(ErrDup)
	}

	_ = config.OnUpdate(func(c *config.Config) error {
		for i := range c.CLimiters {
			if c.CLimiters[i].Name == srvName {
				c.CLimiters[i] = &data
				break
			}
		}
		return nil
	})

	_ = ctx.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("%s 更新成功!", data.Name)})
	_ = Event.OnClickConnLimiters(ctx)

	return handlers.EndConversation()
}

func cancelConnLimiterHandler(ctx telebot.Context) error {
	err := ctx.Reply("添加 并发连接数限制器 已被取消。 还有什么我可以为你做的吗？", &telebot.SendOptions{})
	if err != nil {
		return fmt.Errorf("failed to send cancelHandler message: %w", err)
	}
	//return handlers.EndConversation()
	return nil
}
