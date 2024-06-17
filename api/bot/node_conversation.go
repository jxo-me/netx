package bot

import (
	"fmt"
	"github.com/gogf/gf/v2/util/gconv"
	telebot "github.com/jxo-me/gfbot"
	"github.com/jxo-me/gfbot/handlers"
	"html"
)

const (
	NODENAME   = "nodeName"
	NODETOKEN  = "nodeToken"
	NODEDOMAIN = "nodeDomain"
)

func AddNodeConversation(entry, cancel string) handlers.Conversation {
	return handlers.NewConversation(
		entry,
		telebot.HandlerFunc(nodeStartHandler), // 入口
		map[string][]telebot.IHandler{
			NODENAME:   {telebot.HandlerFunc(nodeNameHandler)},
			NODETOKEN:  {telebot.HandlerFunc(nodeTokenHandler)},
			NODEDOMAIN: {telebot.HandlerFunc(nodeDomainHandler)},
		}, // states状态
		&handlers.ConversationOpts{
			ExitName:     cancel,
			ExitHandler:  telebot.HandlerFunc(cancelNodeHandler),
			AllowReEntry: true,
		}, // options
	)
}

func nodeStartHandler(ctx telebot.Context) error {
	err := ctx.Reply(fmt.Sprintf("你好, @%s.\n请输入新增节点名称?\n您可以随时键入 /cancel 来取消该过程。", ctx.Sender().Username), &telebot.SendOptions{})
	if err != nil {
		return fmt.Errorf("failed to send start message: %w", err)
	}

	// 设置当前用户下一个入口
	return handlers.NextConversationState(NODENAME)
}

// cancelNodeHandler cancels the conversation.
func cancelNodeHandler(ctx telebot.Context) error {
	err := ctx.Reply("命令 AddNode 已被取消。 还有什么我可以为你做的吗？", &telebot.SendOptions{})
	if err != nil {
		return fmt.Errorf("failed to send cancelHandler message: %w", err)
	}
	//return handlers.EndConversation()
	return nil
}

// nodeNameHandler gets the user's nodeNameHandler
func nodeNameHandler(ctx telebot.Context) error {
	inputName := ctx.Message().Text
	err := ctx.Reply(fmt.Sprintf("新增机器人名称:%s!\n\n 请输入机器人Token?", html.EscapeString(inputName)), &telebot.SendOptions{})
	if err != nil {
		return fmt.Errorf("failed to send name message: %w", err)
	}

	_ = ctx.Bot().Store().UpdateData(ctx, NODENAME, inputName)

	return handlers.NextConversationState(NODETOKEN)
}

// nodeTokenHandler gets the user's nodeTokenHandler
func nodeTokenHandler(ctx telebot.Context) error {
	inputToken := ctx.Message().Text
	if len(inputToken) != 46 {
		// If the number is not valid, try again!
		_ = ctx.Reply(fmt.Sprintf("输入的Token格式错误！请重新输入?"), &telebot.SendOptions{})
		// We try the ageHandler handler again
		return handlers.NextConversationState(NODETOKEN)
	}

	err := ctx.Reply(fmt.Sprintf("Token %s\n 请输入绑定的域名?", inputToken), &telebot.SendOptions{})
	if err != nil {
		return fmt.Errorf("failed to send ageHandler message: %w", err)
	}
	_ = ctx.Bot().Store().UpdateData(ctx, NODETOKEN, inputToken)
	return handlers.NextConversationState(NODEDOMAIN)
}

func nodeDomainHandler(ctx telebot.Context) error {
	inputDomain := ctx.Message().Text
	nodeName := ""
	nodeToken := ""
	s, err := ctx.Bot().Store().Get(ctx)
	if err != nil {
		return err
	}
	fmt.Println("ctx data:", s.Data)
	if n, ok := s.Data[NODENAME]; ok {
		nodeName = gconv.String(n)
	}
	if n, ok := s.Data[NODETOKEN]; ok {
		nodeToken = gconv.String(n)
	}
	err = ctx.Reply(fmt.Sprintf("节点名称:%s\nToken: %s\n域名: %s", nodeName, nodeToken, html.EscapeString(inputDomain)), &telebot.SendOptions{})
	if err != nil {
		return fmt.Errorf("failed to send name message: %w", err)
	}

	_ = ctx.Bot().Store().UpdateData(ctx, NODEDOMAIN, inputDomain)
	// 完成
	return handlers.EndConversation()
}
