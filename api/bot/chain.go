package bot

import (
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/v2/util/gconv"
	telebot "github.com/jxo-me/gfbot"
	"github.com/jxo-me/gfbot/handlers"
	"github.com/jxo-me/netx/core/logger"
	"github.com/jxo-me/netx/x/app"
	"github.com/jxo-me/netx/x/config"
	parser "github.com/jxo-me/netx/x/config/parsing/chain"
)

const (
	ChainAdd         = "chainAdd"
	ChainUpdate      = "chainUpdate"
	ChainExampleJson = `
{
  "name": "chain-0",
  "hops": [
    {
      "name": "hop-0",
      "nodes": [
        {
          "name": "node-0",
          "addr": "192.168.1.1:8080",
          "connector": {
            "type": "http"
          },
          "dialer": {
            "type": "tls"
          }
        }
      ]
    }
  ]
}
`
)

func (h *hEvent) OnClickChains(c telebot.Context) error {
	var (
		msg string
	)

	cfg := config.Global()
	rowList := make([]telebot.Row, 0)
	btnList := make([]telebot.Btn, 0)
	selector := &telebot.ReplyMarkup{}
	for _, service := range cfg.Chains {
		btnList = append(btnList, selector.Data(fmt.Sprintf("@%s", service.Name), "detailChain", service.Name))
	}
	rowList = append(rowList, selector.Split(MaxCol, btnList)...)
	rowList = append(rowList, selector.Row(
		selector.Data("@添加转发链", "addChain", "addChain"),
		selector.Data("« 返回 服务列表", "backServices", "backServices"),
	))

	selector.Inline(
		rowList...,
	)
	msg = "Chains List:\n"
	if c.Callback() != nil {
		return c.Edit(msg, &telebot.SendOptions{ReplyMarkup: selector, ParseMode: telebot.ModeMarkdownV2})
	}

	return c.Reply(msg, &telebot.SendOptions{ReplyMarkup: selector, ParseMode: telebot.ModeMarkdownV2})
}

func (h *hEvent) OnClickDetailChain(c telebot.Context) error {
	var (
		msg string
		str string
		err error
	)

	serviceName := c.Callback().Data
	cfg := config.Global()
	var srv *config.ChainConfig
	for _, service := range cfg.Chains {
		if service.Name == serviceName {
			srv = service
		}
	}
	if srv != nil {
		str, err = ConvertJsonMsg(srv)
		if err != nil {
			return c.Reply("OnClickDetailChain ConvertJsonMsg err:", err.Error())
		}
		msg = fmt.Sprintf(CodeMsgTpl, msg, CodeStart, str, CodeEnd)
	}
	selector := &telebot.ReplyMarkup{}
	selector.Inline(
		selector.Row(
			selector.Data("@更新转发链", "updateChain", serviceName),
			selector.Data("@删除转发链", "delChain", serviceName),
		),
		selector.Row(
			selector.Data("« 返回 服务列表", "backServices", "backServices"),
		),
	)
	return c.Edit(msg, &telebot.SendOptions{ReplyMarkup: selector, ParseMode: telebot.ModeMarkdownV2})
}

func (h *hEvent) OnClickDelChain(c telebot.Context) error {
	serviceName := c.Callback().Data
	svc := app.Runtime.ChainRegistry().Get(serviceName)
	if svc == nil {
		return c.Send(ErrNotFound)
	}

	app.Runtime.ChainRegistry().Unregister(serviceName)

	_ = config.OnUpdate(func(c *config.Config) error {
		chaines := c.Chains
		c.Chains = nil
		for _, s := range chaines {
			if s.Name == serviceName {
				continue
			}
			c.Chains = append(c.Chains, s)
		}
		return nil
	})
	_ = c.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("%s 删除成功!", serviceName)})
	return h.OnClickChains(c)
}

func AddChainConversation(entry, cancel string) handlers.Conversation {
	return handlers.NewConversation(
		entry,
		telebot.HandlerFunc(startAddChainHandler), // 入口
		map[string][]telebot.IHandler{
			ChainAdd: {telebot.HandlerFunc(addChainHandler)},
		}, // states状态
		&handlers.ConversationOpts{
			ExitName:     cancel,
			ExitHandler:  telebot.HandlerFunc(cancelChainHandler),
			AllowReEntry: true,
		}, // options
	)
}

func UpdateChainConversation(entry, cancel string) handlers.Conversation {
	return handlers.NewConversation(
		entry,
		telebot.HandlerFunc(startUpdateChainHandler), // 入口
		map[string][]telebot.IHandler{
			ChainUpdate: {telebot.HandlerFunc(updateChainHandler)},
		}, // states状态
		&handlers.ConversationOpts{
			ExitName:     cancel,
			ExitHandler:  telebot.HandlerFunc(cancelChainHandler),
			AllowReEntry: true,
		}, // options
	)
}

func startAddChainHandler(ctx telebot.Context) error {
	example := fmt.Sprintf(CodeTpl, CodeStart, ChainExampleJson, CodeEnd)
	err := ctx.Send(fmt.Sprintf("请输入 转发链 配置?\nExample：%s\n您可以随时键入 /cancel 来取消该操作。", example), &telebot.SendOptions{ParseMode: telebot.ModeMarkdownV2})
	if err != nil {
		return fmt.Errorf("failed to send start message: %w", err)
	}

	// 设置当前用户下一个入口
	return handlers.NextConversationState(ChainAdd)
}

func addChainHandler(ctx telebot.Context) error {
	var (
		data config.ChainConfig
	)

	str := ctx.Message().Text
	err := json.Unmarshal([]byte(str), &data)
	if err != nil {
		return ctx.Reply("addChainHandler json.Unmarshal error:", err.Error())
	}
	v, err := parser.ParseChain(&data, logger.Default())
	if err != nil {
		return ctx.Reply(ErrCreate)
	}
	if err = app.Runtime.ChainRegistry().Register(data.Name, v); err != nil {
		return ctx.Reply(ErrDup)
	}

	_ = config.OnUpdate(func(c *config.Config) error {
		c.Chains = append(c.Chains, &data)
		return nil
	})

	_ = ctx.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("%s 添加成功!", data.Name)})
	_ = Event.OnClickChains(ctx)

	return handlers.EndConversation()
}

func startUpdateChainHandler(ctx telebot.Context) error {
	srvName := ctx.Callback().Data
	err := ctx.Bot().Store().UpdateData(ctx, ChainUpdate, srvName)
	if err != nil {
		return fmt.Errorf("failed UpdateData message: %w", err)
	}
	example := fmt.Sprintf(CodeTpl, CodeStart, ChainExampleJson, CodeEnd)
	err = ctx.Send(fmt.Sprintf("请输入 转发链 配置?\nExample：%s\n您可以随时键入 /cancel 来取消该操作。", example), &telebot.SendOptions{ParseMode: telebot.ModeMarkdownV2})
	if err != nil {
		return fmt.Errorf("failed to send start message: %w", err)
	}

	// 设置当前用户下一个入口
	return handlers.NextConversationState(ChainUpdate)
}

func updateChainHandler(ctx telebot.Context) error {
	var (
		data config.ChainConfig
	)
	state, err := ctx.Bot().Store().Get(ctx)
	if err != nil {
		return ctx.Reply("updateChainHandler Store().Get error:", err.Error())
	}
	srvName := ""
	if sn, ok := state.Data[ChainUpdate]; ok {
		srvName = gconv.String(sn)
	}
	fmt.Println(fmt.Sprintf("srvName :%s", srvName))
	if srvName == "" {
		return ctx.Reply(ErrInvalid)
	}
	str := ctx.Message().Text
	err = json.Unmarshal([]byte(str), &data)
	if err != nil {
		return ctx.Reply("updateChainHandler json.Unmarshal error:", err.Error())
	}

	if !app.Runtime.ChainRegistry().IsRegistered(srvName) {
		return ctx.Reply(ErrNotFound)
	}

	data.Name = srvName

	v, err := parser.ParseChain(&data, logger.Default())
	if err != nil {
		return ctx.Reply(ErrCreate)
	}
	app.Runtime.ChainRegistry().Unregister(srvName)

	if err = app.Runtime.ChainRegistry().Register(srvName, v); err != nil {
		return ctx.Reply(ErrDup)
	}

	_ = config.OnUpdate(func(c *config.Config) error {
		for i := range c.Chains {
			if c.Chains[i].Name == srvName {
				c.Chains[i] = &data
				break
			}
		}
		return nil
	})

	_ = ctx.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("%s 更新成功!", data.Name)})
	_ = Event.OnClickChains(ctx)

	return handlers.EndConversation()
}

func cancelChainHandler(ctx telebot.Context) error {
	err := ctx.Reply("添加 转发链 已被取消。 还有什么我可以为你做的吗？", &telebot.SendOptions{})
	if err != nil {
		return fmt.Errorf("failed to send cancelHandler message: %w", err)
	}
	//return handlers.EndConversation()
	return nil
}
