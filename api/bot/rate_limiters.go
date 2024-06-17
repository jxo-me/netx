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
	RateLimiterAdd      = "rateLimiterAdd"
	RateLimiterUpdate   = "rateLimiterUpdate"
	RLimiterExampleJson = `
{
  "name": "rlimiter-0",
  "limits": [
    "$ 100",
    "$$ 10",
    "192.168.1.1  50",
    "192.168.0.0/16  5"
  ]
}
`
)

func (h *hEvent) OnClickRateLimiters(c telebot.Context) error {
	var (
		msg string
	)

	cfg := config.Global()
	rowList := make([]telebot.Row, 0)
	btnList := make([]telebot.Btn, 0)
	selector := &telebot.ReplyMarkup{}
	for _, service := range cfg.RLimiters {
		btnList = append(btnList, selector.Data(fmt.Sprintf("@%s", service.Name), "detailRateLimiter", service.Name))
	}
	rowList = append(rowList, selector.Split(MaxCol, btnList)...)
	rowList = append(rowList, selector.Row(
		selector.Data("@添加速率限制器", "addRateLimiter", "addRateLimiter"),
		selector.Data("« 返回 服务列表", "backServices", "backServices"),
	))

	selector.Inline(
		rowList...,
	)
	msg = "RateLimiters List:\n"
	if c.Callback() != nil {
		return c.Edit(msg, &telebot.SendOptions{ReplyMarkup: selector, ParseMode: telebot.ModeMarkdownV2})
	}

	return c.Reply(msg, &telebot.SendOptions{ReplyMarkup: selector, ParseMode: telebot.ModeMarkdownV2})
}

func (h *hEvent) OnClickDetailRateLimiter(c telebot.Context) error {
	var (
		msg string
		str string
		err error
	)

	serviceName := c.Callback().Data
	cfg := config.Global()
	var srv *config.LimiterConfig
	for _, service := range cfg.RLimiters {
		if service.Name == serviceName {
			srv = service
		}
	}
	if srv != nil {
		str, err = ConvertJsonMsg(srv)
		if err != nil {
			return c.Reply("OnClickDetailRateLimiter ConvertJsonMsg err:", err.Error())
		}
		msg = fmt.Sprintf(CodeMsgTpl, msg, CodeStart, str, CodeEnd)
	}
	selector := &telebot.ReplyMarkup{}
	selector.Inline(
		selector.Row(
			selector.Data("@更新速率限制器", "updateRateLimiter", serviceName),
			selector.Data("@删除速率限制器", "delRateLimiter", serviceName),
		),
		selector.Row(
			selector.Data("« 返回 服务列表", "backServices", "backServices"),
		),
	)
	return c.Edit(msg, &telebot.SendOptions{ReplyMarkup: selector, ParseMode: telebot.ModeMarkdownV2})
}

func (h *hEvent) OnClickDelRateLimiter(c telebot.Context) error {
	serviceName := c.Callback().Data
	svc := app.Runtime.RateLimiterRegistry().Get(serviceName)
	if svc == nil {
		return c.Send(ErrNotFound)
	}

	app.Runtime.RateLimiterRegistry().Unregister(serviceName)

	_ = config.OnUpdate(func(c *config.Config) error {
		rateLimiteres := c.RLimiters
		c.RLimiters = nil
		for _, s := range rateLimiteres {
			if s.Name == serviceName {
				continue
			}
			c.RLimiters = append(c.RLimiters, s)
		}
		return nil
	})
	_ = c.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("%s 删除成功!", serviceName)})
	return h.OnClickRateLimiters(c)
}

func AddRateLimiterConversation(entry, cancel string) handlers.Conversation {
	return handlers.NewConversation(
		entry,
		telebot.HandlerFunc(startAddRateLimiterHandler), // 入口
		map[string][]telebot.IHandler{
			RateLimiterAdd: {telebot.HandlerFunc(addRateLimiterHandler)},
		}, // states状态
		&handlers.ConversationOpts{
			ExitName:     cancel,
			ExitHandler:  telebot.HandlerFunc(cancelRateLimiterHandler),
			AllowReEntry: true,
		}, // options
	)
}

func UpdateRateLimiterConversation(entry, cancel string) handlers.Conversation {
	return handlers.NewConversation(
		entry,
		telebot.HandlerFunc(startUpdateRateLimiterHandler), // 入口
		map[string][]telebot.IHandler{
			RateLimiterUpdate: {telebot.HandlerFunc(updateRateLimiterHandler)},
		}, // states状态
		&handlers.ConversationOpts{
			ExitName:     cancel,
			ExitHandler:  telebot.HandlerFunc(cancelRateLimiterHandler),
			AllowReEntry: true,
		}, // options
	)
}

func startAddRateLimiterHandler(ctx telebot.Context) error {
	example := fmt.Sprintf(CodeTpl, CodeStart, RLimiterExampleJson, CodeEnd)
	err := ctx.Send(fmt.Sprintf("请输入 速率限制器 配置?\nExample：%s\n您可以随时键入 /cancel 来取消该操作。", example), &telebot.SendOptions{ParseMode: telebot.ModeMarkdownV2})
	if err != nil {
		return fmt.Errorf("failed to send start message: %w", err)
	}

	// 设置当前用户下一个入口
	return handlers.NextConversationState(RateLimiterAdd)
}

func addRateLimiterHandler(ctx telebot.Context) error {
	var (
		data config.LimiterConfig
	)

	str := ctx.Message().Text
	err := json.Unmarshal([]byte(str), &data)
	if err != nil {
		return ctx.Reply("addRateLimiterHandler json.Unmarshal error:", err.Error())
	}
	v := parser.ParseRateLimiter(&data)
	if err = app.Runtime.RateLimiterRegistry().Register(data.Name, v); err != nil {
		return ctx.Reply(ErrDup)
	}

	_ = config.OnUpdate(func(c *config.Config) error {
		c.RLimiters = append(c.RLimiters, &data)
		return nil
	})

	_ = ctx.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("%s 添加成功!", data.Name)})
	_ = Event.OnClickRateLimiters(ctx)

	return handlers.EndConversation()
}

func startUpdateRateLimiterHandler(ctx telebot.Context) error {
	srvName := ctx.Callback().Data
	err := ctx.Bot().Store().UpdateData(ctx, RateLimiterUpdate, srvName)
	if err != nil {
		return fmt.Errorf("failed UpdateData message: %w", err)
	}
	example := fmt.Sprintf(CodeTpl, CodeStart, RLimiterExampleJson, CodeEnd)
	err = ctx.Send(fmt.Sprintf("请输入 速率限制器 配置?\nExample：%s\n您可以随时键入 /cancel 来取消该操作。", example), &telebot.SendOptions{ParseMode: telebot.ModeMarkdownV2})
	if err != nil {
		return fmt.Errorf("failed to send start message: %w", err)
	}

	// 设置当前用户下一个入口
	return handlers.NextConversationState(RateLimiterUpdate)
}

func updateRateLimiterHandler(ctx telebot.Context) error {
	var (
		data config.LimiterConfig
	)
	state, err := ctx.Bot().Store().Get(ctx)
	if err != nil {
		return ctx.Reply("updateRateLimiterHandler Store().Get error:", err.Error())
	}
	srvName := ""
	if sn, ok := state.Data[RateLimiterUpdate]; ok {
		srvName = gconv.String(sn)
	}
	fmt.Println(fmt.Sprintf("srvName :%s", srvName))
	if srvName == "" {
		return ctx.Reply(ErrInvalid)
	}
	str := ctx.Message().Text
	err = json.Unmarshal([]byte(str), &data)
	if err != nil {
		return ctx.Reply("updateRateLimiterHandler json.Unmarshal error:", err.Error())
	}

	if !app.Runtime.RateLimiterRegistry().IsRegistered(srvName) {
		return ctx.Reply(ErrNotFound)
	}

	data.Name = srvName

	v := parser.ParseRateLimiter(&data)

	app.Runtime.RateLimiterRegistry().Unregister(srvName)

	if err = app.Runtime.RateLimiterRegistry().Register(srvName, v); err != nil {
		return ctx.Reply(ErrDup)
	}

	_ = config.OnUpdate(func(c *config.Config) error {
		for i := range c.RLimiters {
			if c.RLimiters[i].Name == srvName {
				c.RLimiters[i] = &data
				break
			}
		}
		return nil
	})

	_ = ctx.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("%s 更新成功!", data.Name)})
	_ = Event.OnClickRateLimiters(ctx)

	return handlers.EndConversation()
}

func cancelRateLimiterHandler(ctx telebot.Context) error {
	err := ctx.Reply("添加 速率限制器 已被取消。 还有什么我可以为你做的吗？", &telebot.SendOptions{})
	if err != nil {
		return fmt.Errorf("failed to send cancelHandler message: %w", err)
	}
	//return handlers.EndConversation()
	return nil
}
