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
	AdmissionAdd    = "admissionAdd"
	AdmissionUpdate = "admissionUpdate"
)

func (h *hEvent) OnClickAdmissions(c telebot.IContext) error {
	var (
		msg string
	)

	cfg := config.Global()
	rowList := make([]telebot.Row, 0)
	btnList := make([]telebot.Btn, 0)
	selector := &telebot.ReplyMarkup{}
	for i, service := range cfg.Admissions {
		btnList = append(btnList, selector.Data(fmt.Sprintf("@%s", service.Name), "detailAdmission", service.Name))
		if i%3 == 0 {
			rowList = append(rowList, selector.Row(btnList...))
			btnList = make([]telebot.Btn, 0)
		}
	}
	rowList = append(rowList, selector.Row(
		selector.Data("@添加准入控制器", "addAdmission", "addAdmission"),
		selector.Data("« 返回 服务列表", "backServices", "backServices"),
	))

	selector.Inline(
		rowList...,
	)
	msg = "Admissions List:\n"
	if c.Callback() != nil {
		return c.Edit(msg, &telebot.SendOptions{ReplyMarkup: selector, ParseMode: telebot.ModeMarkdownV2})
	}

	return c.Reply(msg, &telebot.SendOptions{ReplyMarkup: selector, ParseMode: telebot.ModeMarkdownV2})
}

func (h *hEvent) OnClickDetailAdmission(c telebot.IContext) error {
	var (
		msg string
		str string
		err error
	)

	serviceName := c.Callback().Data
	cfg := config.Global()
	var srv *config.AdmissionConfig
	for _, service := range cfg.Admissions {
		if service.Name == serviceName {
			srv = service
		}
	}
	if srv != nil {
		str, err = ConvertJsonMsg(srv)
		if err != nil {
			return c.Reply("OnClickDetailAdmission ConvertJsonMsg err:", err.Error())
		}
		msg = fmt.Sprintf(CodeMsgTpl, msg, CodeStart, str, CodeEnd)
	}
	selector := &telebot.ReplyMarkup{}
	selector.Inline(
		selector.Row(
			selector.Data("@更新准入控制器", "updateAdmission", serviceName),
			selector.Data("@删除准入控制器", "delAdmission", serviceName),
		),
		selector.Row(
			selector.Data("« 返回 服务列表", "backServices", "backServices"),
		),
	)
	return c.Edit(msg, &telebot.SendOptions{ReplyMarkup: selector, ParseMode: telebot.ModeMarkdownV2})
}

func (h *hEvent) OnClickDelAdmission(c telebot.IContext) error {
	serviceName := c.Callback().Data
	svc := app.Runtime.AdmissionRegistry().Get(serviceName)
	if svc == nil {
		return c.Send(ErrNotFound)
	}

	app.Runtime.AdmissionRegistry().Unregister(serviceName)

	_ = config.OnUpdate(func(c *config.Config) error {
		admissiones := c.Admissions
		c.Admissions = nil
		for _, s := range admissiones {
			if s.Name == serviceName {
				continue
			}
			c.Admissions = append(c.Admissions, s)
		}
		return nil
	})
	_ = c.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("%s 删除成功!", serviceName)})
	return h.OnClickAdmissions(c)
}

func AddAdmissionConversation(entry, cancel string) handlers.Conversation {
	return handlers.NewConversation(
		entry,
		telebot.HandlerFunc(startAddAdmissionHandler), // 入口
		map[string][]telebot.IHandler{
			AdmissionAdd: {telebot.HandlerFunc(addAdmissionHandler)},
		}, // states状态
		&handlers.ConversationOpts{
			ExitName:     cancel,
			ExitHandler:  telebot.HandlerFunc(cancelAdmissionHandler),
			AllowReEntry: true,
		}, // options
	)
}

func UpdateAdmissionConversation(entry, cancel string) handlers.Conversation {
	return handlers.NewConversation(
		entry,
		telebot.HandlerFunc(startUpdateAdmissionHandler), // 入口
		map[string][]telebot.IHandler{
			AdmissionUpdate: {telebot.HandlerFunc(updateAdmissionHandler)},
		}, // states状态
		&handlers.ConversationOpts{
			ExitName:     cancel,
			ExitHandler:  telebot.HandlerFunc(cancelAdmissionHandler),
			AllowReEntry: true,
		}, // options
	)
}

func startAddAdmissionHandler(ctx telebot.IContext) error {
	err := ctx.Send(fmt.Sprintf("你好, @%s.\n请输入准入控制器JSON配置?\n您可以随时键入 /cancel 来取消该操作。", ctx.Sender().Username), &telebot.SendOptions{})
	if err != nil {
		return fmt.Errorf("failed to send start message: %w", err)
	}

	// 设置当前用户下一个入口
	return handlers.NextConversationState(AdmissionAdd)
}

func addAdmissionHandler(ctx telebot.IContext) error {
	var (
		data config.AdmissionConfig
	)

	str := ctx.Message().Text
	err := json.Unmarshal([]byte(str), &data)
	if err != nil {
		return ctx.Reply("addAdmissionHandler json.Unmarshal error:", err.Error())
	}
	v := parsing.ParseAdmission(&data)
	if err = app.Runtime.AdmissionRegistry().Register(data.Name, v); err != nil {
		return ctx.Reply(ErrDup)
	}

	_ = config.OnUpdate(func(c *config.Config) error {
		c.Admissions = append(c.Admissions, &data)
		return nil
	})

	_ = ctx.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("%s 添加成功!", data.Name)})
	_ = Event.OnClickAdmissions(ctx)

	return handlers.EndConversation()
}

func startUpdateAdmissionHandler(ctx telebot.IContext) error {
	srvName := ctx.Callback().Data
	err := ctx.Bot().Store().UpdateData(ctx, AdmissionUpdate, srvName)
	if err != nil {
		return fmt.Errorf("failed UpdateData message: %w", err)
	}

	err = ctx.Send(fmt.Sprintf("你好, @%s.\n请输入服务JSON配置?\n您可以随时键入 /cancel 来取消该操作。", ctx.Sender().Username), &telebot.SendOptions{})
	if err != nil {
		return fmt.Errorf("failed to send start message: %w", err)
	}

	// 设置当前用户下一个入口
	return handlers.NextConversationState(AdmissionUpdate)
}

func updateAdmissionHandler(ctx telebot.IContext) error {
	var (
		data config.AdmissionConfig
	)
	state, err := ctx.Bot().Store().Get(ctx)
	if err != nil {
		return ctx.Reply("updateAdmissionHandler Store().Get error:", err.Error())
	}
	srvName := ""
	if sn, ok := state.Data[AdmissionUpdate]; ok {
		srvName = gconv.String(sn)
	}
	fmt.Println(fmt.Sprintf("srvName :%s", srvName))
	if srvName == "" {
		return ctx.Reply(ErrInvalid)
	}
	str := ctx.Message().Text
	err = json.Unmarshal([]byte(str), &data)
	if err != nil {
		return ctx.Reply("updateAdmissionHandler json.Unmarshal error:", err.Error())
	}

	if !app.Runtime.AdmissionRegistry().IsRegistered(srvName) {
		return ctx.Reply(ErrNotFound)
	}

	data.Name = srvName

	v := parsing.ParseAdmission(&data)

	app.Runtime.AdmissionRegistry().Unregister(srvName)

	if err = app.Runtime.AdmissionRegistry().Register(srvName, v); err != nil {
		return ctx.Reply(ErrDup)
	}

	_ = config.OnUpdate(func(c *config.Config) error {
		for i := range c.Admissions {
			if c.Admissions[i].Name == srvName {
				c.Admissions[i] = &data
				break
			}
		}
		return nil
	})

	_ = ctx.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("%s 更新成功!", data.Name)})
	_ = Event.OnClickAdmissions(ctx)

	return handlers.EndConversation()
}

func cancelAdmissionHandler(ctx telebot.IContext) error {
	err := ctx.Reply("添加 准入控制器 已被取消。 还有什么我可以为你做的吗？", &telebot.SendOptions{})
	if err != nil {
		return fmt.Errorf("failed to send cancelHandler message: %w", err)
	}
	//return handlers.EndConversation()
	return nil
}
