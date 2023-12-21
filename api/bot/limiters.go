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
	LimiterAdd         = "limiterAdd"
	LimiterUpdate      = "limiterUpdate"
	LimiterExampleJson = `
{
  "name": "limiter-0",
  "limits": [
    "$ 100MB  200MB",
    "$$ 10MB",
    "192.168.1.1  1MB 10MB",
    "192.168.0.0/16  512KB  1MB"
  ]
}
`
)

func (h *hEvent) OnClickLimiters(c telebot.IContext) error {
	var (
		msg string
	)

	cfg := config.Global()
	rowList := make([]telebot.Row, 0)
	btnList := make([]telebot.Btn, 0)
	selector := &telebot.ReplyMarkup{}
	for _, service := range cfg.Limiters {
		btnList = append(btnList, selector.Data(fmt.Sprintf("@%s", service.Name), "detailLimiter", service.Name))
	}
	rowList = append(rowList, selector.Split(MaxCol, btnList)...)
	rowList = append(rowList, selector.Row(
		selector.Data("@添加流量速率限制器", "addLimiter", "addLimiter"),
		selector.Data("« 返回 服务列表", "backServices", "backServices"),
	))

	selector.Inline(
		rowList...,
	)
	msg = "Limiters List:\n"
	if c.Callback() != nil {
		return c.Edit(msg, &telebot.SendOptions{ReplyMarkup: selector, ParseMode: telebot.ModeMarkdownV2})
	}

	return c.Reply(msg, &telebot.SendOptions{ReplyMarkup: selector, ParseMode: telebot.ModeMarkdownV2})
}

func (h *hEvent) OnClickDetailLimiter(c telebot.IContext) error {
	var (
		msg string
		str string
		err error
	)

	serviceName := c.Callback().Data
	cfg := config.Global()
	var srv *config.LimiterConfig
	for _, service := range cfg.Limiters {
		if service.Name == serviceName {
			srv = service
		}
	}
	if srv != nil {
		str, err = ConvertJsonMsg(srv)
		if err != nil {
			return c.Reply("OnClickDetailLimiter ConvertJsonMsg err:", err.Error())
		}
		msg = fmt.Sprintf(CodeMsgTpl, msg, CodeStart, str, CodeEnd)
	}
	selector := &telebot.ReplyMarkup{}
	selector.Inline(
		selector.Row(
			selector.Data("@更新流量速率限制器", "updateLimiter", serviceName),
			selector.Data("@删除流量速率限制器", "delLimiter", serviceName),
		),
		selector.Row(
			selector.Data("« 返回 服务列表", "backServices", "backServices"),
		),
	)
	return c.Edit(msg, &telebot.SendOptions{ReplyMarkup: selector, ParseMode: telebot.ModeMarkdownV2})
}

func (h *hEvent) OnClickDelLimiter(c telebot.IContext) error {
	serviceName := c.Callback().Data
	svc := app.Runtime.TrafficLimiterRegistry().Get(serviceName)
	if svc == nil {
		return c.Send(ErrNotFound)
	}

	app.Runtime.TrafficLimiterRegistry().Unregister(serviceName)

	_ = config.OnUpdate(func(c *config.Config) error {
		limiteres := c.Limiters
		c.Limiters = nil
		for _, s := range limiteres {
			if s.Name == serviceName {
				continue
			}
			c.Limiters = append(c.Limiters, s)
		}
		return nil
	})
	_ = c.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("%s 删除成功!", serviceName)})
	return h.OnClickLimiters(c)
}

func AddLimiterConversation(entry, cancel string) handlers.Conversation {
	return handlers.NewConversation(
		entry,
		telebot.HandlerFunc(startAddLimiterHandler), // 入口
		map[string][]telebot.IHandler{
			LimiterAdd: {telebot.HandlerFunc(addLimiterHandler)},
		}, // states状态
		&handlers.ConversationOpts{
			ExitName:     cancel,
			ExitHandler:  telebot.HandlerFunc(cancelLimiterHandler),
			AllowReEntry: true,
		}, // options
	)
}

func UpdateLimiterConversation(entry, cancel string) handlers.Conversation {
	return handlers.NewConversation(
		entry,
		telebot.HandlerFunc(startUpdateLimiterHandler), // 入口
		map[string][]telebot.IHandler{
			LimiterUpdate: {telebot.HandlerFunc(updateLimiterHandler)},
		}, // states状态
		&handlers.ConversationOpts{
			ExitName:     cancel,
			ExitHandler:  telebot.HandlerFunc(cancelLimiterHandler),
			AllowReEntry: true,
		}, // options
	)
}

func startAddLimiterHandler(ctx telebot.IContext) error {
	example := fmt.Sprintf(CodeTpl, CodeStart, LimiterExampleJson, CodeEnd)
	err := ctx.Send(fmt.Sprintf("请输入 流量速率限制器 配置?\nExample：%s\n您可以随时键入 /cancel 来取消该操作。", example), &telebot.SendOptions{ParseMode: telebot.ModeMarkdownV2})
	if err != nil {
		return fmt.Errorf("failed to send start message: %w", err)
	}

	// 设置当前用户下一个入口
	return handlers.NextConversationState(LimiterAdd)
}

func addLimiterHandler(ctx telebot.IContext) error {
	var (
		data config.LimiterConfig
	)

	str := ctx.Message().Text
	err := json.Unmarshal([]byte(str), &data)
	if err != nil {
		return ctx.Reply("addLimiterHandler json.Unmarshal error:", err.Error())
	}
	v := parser.ParseTrafficLimiter(&data)
	if err = app.Runtime.TrafficLimiterRegistry().Register(data.Name, v); err != nil {
		return ctx.Reply(ErrDup)
	}

	_ = config.OnUpdate(func(c *config.Config) error {
		c.Limiters = append(c.Limiters, &data)
		return nil
	})

	_ = ctx.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("%s 添加成功!", data.Name)})
	_ = Event.OnClickLimiters(ctx)

	return handlers.EndConversation()
}

func startUpdateLimiterHandler(ctx telebot.IContext) error {
	srvName := ctx.Callback().Data
	err := ctx.Bot().Store().UpdateData(ctx, LimiterUpdate, srvName)
	if err != nil {
		return fmt.Errorf("failed UpdateData message: %w", err)
	}
	example := fmt.Sprintf(CodeTpl, CodeStart, LimiterExampleJson, CodeEnd)
	err = ctx.Send(fmt.Sprintf("请输入 流量速率限制器 配置?\nExample：%s\n您可以随时键入 /cancel 来取消该操作。", example), &telebot.SendOptions{ParseMode: telebot.ModeMarkdownV2})
	if err != nil {
		return fmt.Errorf("failed to send start message: %w", err)
	}

	// 设置当前用户下一个入口
	return handlers.NextConversationState(LimiterUpdate)
}

func updateLimiterHandler(ctx telebot.IContext) error {
	var (
		data config.LimiterConfig
	)
	state, err := ctx.Bot().Store().Get(ctx)
	if err != nil {
		return ctx.Reply("updateLimiterHandler Store().Get error:", err.Error())
	}
	srvName := ""
	if sn, ok := state.Data[LimiterUpdate]; ok {
		srvName = gconv.String(sn)
	}
	fmt.Println(fmt.Sprintf("srvName :%s", srvName))
	if srvName == "" {
		return ctx.Reply(ErrInvalid)
	}
	str := ctx.Message().Text
	err = json.Unmarshal([]byte(str), &data)
	if err != nil {
		return ctx.Reply("updateLimiterHandler json.Unmarshal error:", err.Error())
	}

	if !app.Runtime.TrafficLimiterRegistry().IsRegistered(srvName) {
		return ctx.Reply(ErrNotFound)
	}

	data.Name = srvName

	v := parser.ParseTrafficLimiter(&data)

	app.Runtime.TrafficLimiterRegistry().Unregister(srvName)

	if err = app.Runtime.TrafficLimiterRegistry().Register(srvName, v); err != nil {
		return ctx.Reply(ErrDup)
	}

	_ = config.OnUpdate(func(c *config.Config) error {
		for i := range c.Limiters {
			if c.Limiters[i].Name == srvName {
				c.Limiters[i] = &data
				break
			}
		}
		return nil
	})

	_ = ctx.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("%s 更新成功!", data.Name)})
	_ = Event.OnClickLimiters(ctx)

	return handlers.EndConversation()
}

func cancelLimiterHandler(ctx telebot.IContext) error {
	err := ctx.Reply("添加 流量速率限制器 已被取消。 还有什么我可以为你做的吗？", &telebot.SendOptions{})
	if err != nil {
		return fmt.Errorf("failed to send cancelHandler message: %w", err)
	}
	//return handlers.EndConversation()
	return nil
}
