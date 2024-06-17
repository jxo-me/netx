package bot

import (
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/v2/util/gconv"
	telebot "github.com/jxo-me/gfbot"
	"github.com/jxo-me/gfbot/handlers"
	"github.com/jxo-me/netx/x/app"
	"github.com/jxo-me/netx/x/config"
	parser "github.com/jxo-me/netx/x/config/parsing/auth"
)

const (
	AutherAdd         = "autherAdd"
	AutherUpdate      = "autherUpdate"
	AutherExampleJson = `
{
  "name": "auther-0",
  "auths": [
    {
      "username": "user1",
      "password": "pass1"
    }
  ]
}
`
)

func (h *hEvent) OnClickAuthers(c telebot.Context) error {
	var (
		msg string
	)

	cfg := config.Global()
	rowList := make([]telebot.Row, 0)
	btnList := make([]telebot.Btn, 0)
	selector := &telebot.ReplyMarkup{}
	for _, service := range cfg.Authers {
		btnList = append(btnList, selector.Data(fmt.Sprintf("@%s", service.Name), "detailAuther", service.Name))
	}
	rowList = append(rowList, selector.Split(MaxCol, btnList)...)
	rowList = append(rowList, selector.Row(
		selector.Data("@添加认证器", "addAuther", "addAuther"),
		selector.Data("« 返回 服务列表", "backServices", "backServices"),
	))

	selector.Inline(
		rowList...,
	)
	msg = "Authers List:\n"
	if c.Callback() != nil {
		return c.Edit(msg, &telebot.SendOptions{ReplyMarkup: selector, ParseMode: telebot.ModeMarkdownV2})
	}

	return c.Reply(msg, &telebot.SendOptions{ReplyMarkup: selector, ParseMode: telebot.ModeMarkdownV2})
}

func (h *hEvent) OnClickDetailAuther(c telebot.Context) error {
	var (
		msg string
		str string
		err error
	)

	serviceName := c.Callback().Data
	cfg := config.Global()
	var srv *config.AutherConfig
	for _, service := range cfg.Authers {
		if service.Name == serviceName {
			srv = service
		}
	}
	if srv != nil {
		str, err = ConvertJsonMsg(srv)
		if err != nil {
			return c.Reply("OnClickDetailAuther ConvertJsonMsg err:", err.Error())
		}
		msg = fmt.Sprintf(CodeMsgTpl, msg, CodeStart, str, CodeEnd)
	}
	selector := &telebot.ReplyMarkup{}
	selector.Inline(
		selector.Row(
			selector.Data("@更新认证器", "updateAuther", serviceName),
			selector.Data("@删除认证器", "delAuther", serviceName),
		),
		selector.Row(
			selector.Data("« 返回 服务列表", "backServices", "backServices"),
		),
	)
	return c.Edit(msg, &telebot.SendOptions{ReplyMarkup: selector, ParseMode: telebot.ModeMarkdownV2})
}

func (h *hEvent) OnClickDelAuther(c telebot.Context) error {
	serviceName := c.Callback().Data
	svc := app.Runtime.AutherRegistry().Get(serviceName)
	if svc == nil {
		return c.Send(ErrNotFound)
	}

	app.Runtime.AutherRegistry().Unregister(serviceName)

	_ = config.OnUpdate(func(c *config.Config) error {
		authers := c.Authers
		c.Authers = nil
		for _, s := range authers {
			if s.Name == serviceName {
				continue
			}
			c.Authers = append(c.Authers, s)
		}
		return nil
	})
	_ = c.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("%s 删除成功!", serviceName)})
	return h.OnClickAuthers(c)
}

func AddAutherConversation(entry, cancel string) handlers.Conversation {
	return handlers.NewConversation(
		entry,
		telebot.HandlerFunc(startAddAutherHandler), // 入口
		map[string][]telebot.IHandler{
			AutherAdd: {telebot.HandlerFunc(addAutherHandler)},
		}, // states状态
		&handlers.ConversationOpts{
			ExitName:     cancel,
			ExitHandler:  telebot.HandlerFunc(cancelAutherHandler),
			AllowReEntry: true,
		}, // options
	)
}

func UpdateAutherConversation(entry, cancel string) handlers.Conversation {
	return handlers.NewConversation(
		entry,
		telebot.HandlerFunc(startUpdateAutherHandler), // 入口
		map[string][]telebot.IHandler{
			AutherUpdate: {telebot.HandlerFunc(updateAutherHandler)},
		}, // states状态
		&handlers.ConversationOpts{
			ExitName:     cancel,
			ExitHandler:  telebot.HandlerFunc(cancelAutherHandler),
			AllowReEntry: true,
		}, // options
	)
}

func startAddAutherHandler(ctx telebot.Context) error {
	example := fmt.Sprintf(CodeTpl, CodeStart, AutherExampleJson, CodeEnd)
	err := ctx.Send(fmt.Sprintf("请输入 认证器 配置?\nExample：%s\n您可以随时键入 /cancel 来取消该操作。", example), &telebot.SendOptions{ParseMode: telebot.ModeMarkdownV2})
	if err != nil {
		return fmt.Errorf("failed to send start message: %w", err)
	}

	// 设置当前用户下一个入口
	return handlers.NextConversationState(AutherAdd)
}

func addAutherHandler(ctx telebot.Context) error {
	var (
		data config.AutherConfig
	)

	str := ctx.Message().Text
	err := json.Unmarshal([]byte(str), &data)
	if err != nil {
		return ctx.Reply("addAutherHandler json.Unmarshal error:", err.Error())
	}
	v := parser.ParseAuther(&data)
	if err = app.Runtime.AutherRegistry().Register(data.Name, v); err != nil {
		return ctx.Reply(ErrDup)
	}

	_ = config.OnUpdate(func(c *config.Config) error {
		c.Authers = append(c.Authers, &data)
		return nil
	})

	_ = ctx.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("%s 添加成功!", data.Name)})
	_ = Event.OnClickAuthers(ctx)

	return handlers.EndConversation()
}

func startUpdateAutherHandler(ctx telebot.Context) error {
	srvName := ctx.Callback().Data
	err := ctx.Bot().Store().UpdateData(ctx, AutherUpdate, srvName)
	if err != nil {
		return fmt.Errorf("failed UpdateData message: %w", err)
	}
	example := fmt.Sprintf(CodeTpl, CodeStart, AutherExampleJson, CodeEnd)
	err = ctx.Send(fmt.Sprintf("请输入 认证器 配置?\nExample：%s\n您可以随时键入 /cancel 来取消该操作。", example), &telebot.SendOptions{ParseMode: telebot.ModeMarkdownV2})
	if err != nil {
		return fmt.Errorf("failed to send start message: %w", err)
	}

	// 设置当前用户下一个入口
	return handlers.NextConversationState(AutherUpdate)
}

func updateAutherHandler(ctx telebot.Context) error {
	var (
		data config.AutherConfig
	)
	state, err := ctx.Bot().Store().Get(ctx)
	if err != nil {
		return ctx.Reply("updateAutherHandler Store().Get error:", err.Error())
	}
	srvName := ""
	if sn, ok := state.Data[AutherUpdate]; ok {
		srvName = gconv.String(sn)
	}
	fmt.Println(fmt.Sprintf("srvName :%s", srvName))
	if srvName == "" {
		return ctx.Reply(ErrInvalid)
	}
	str := ctx.Message().Text
	err = json.Unmarshal([]byte(str), &data)
	if err != nil {
		return ctx.Reply("updateAutherHandler json.Unmarshal error:", err.Error())
	}

	if !app.Runtime.AutherRegistry().IsRegistered(srvName) {
		return ctx.Reply(ErrNotFound)
	}

	data.Name = srvName

	v := parser.ParseAuther(&data)

	app.Runtime.AutherRegistry().Unregister(srvName)

	if err = app.Runtime.AutherRegistry().Register(srvName, v); err != nil {
		return ctx.Reply(ErrDup)
	}

	_ = config.OnUpdate(func(c *config.Config) error {
		for i := range c.Authers {
			if c.Authers[i].Name == srvName {
				c.Authers[i] = &data
				break
			}
		}
		return nil
	})

	_ = ctx.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("%s 更新成功!", data.Name)})
	_ = Event.OnClickAuthers(ctx)

	return handlers.EndConversation()
}

func cancelAutherHandler(ctx telebot.Context) error {
	err := ctx.Reply("添加 认证器 已被取消。 还有什么我可以为你做的吗？", &telebot.SendOptions{})
	if err != nil {
		return fmt.Errorf("failed to send cancelHandler message: %w", err)
	}
	//return handlers.EndConversation()
	return nil
}
