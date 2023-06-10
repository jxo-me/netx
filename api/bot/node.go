package bot

import (
	"fmt"
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

func nodeStartHandler(ctx telebot.IContext) error {
	err := ctx.Reply(fmt.Sprintf("你好, @%s.\n请输入节点名称?.", ctx.Sender().Username), &telebot.SendOptions{})
	if err != nil {
		return fmt.Errorf("failed to send start message: %w", err)
	}
	// 设置当前用户下一个入口
	return handlers.NextConversationState(handlers.Entry, NODENAME)
}

// cancelNodeHandler cancels the conversation.
func cancelNodeHandler(ctx telebot.IContext) error {
	err := ctx.Reply("命令 AddNode 已被取消。 还有什么我可以为你做的吗？", &telebot.SendOptions{})
	if err != nil {
		return fmt.Errorf("failed to send cancelHandler message: %w", err)
	}
	//return handlers.EndConversation()
	return nil
}

// nodeNameHandler gets the user's nodeNameHandler
func nodeNameHandler(ctx telebot.IContext) error {
	inputName := ctx.Message().Text
	err := ctx.Reply(fmt.Sprintf("新增机器人名称:%s!\n\n 请输入机器人Token?", html.EscapeString(inputName)), &telebot.SendOptions{})
	if err != nil {
		return fmt.Errorf("failed to send name message: %w", err)
	}
	return handlers.NextConversationStateWithData(NODENAME, NODETOKEN, inputName)
}

// nodeTokenHandler gets the user's nodeTokenHandler
func nodeTokenHandler(ctx telebot.IContext) error {
	inputToken := ctx.Message().Text
	if len(inputToken) != 46 {
		// If the number is not valid, try again!
		_ = ctx.Reply(fmt.Sprintf("输入的Token格式错误！请重新输入?"), &telebot.SendOptions{})
		// We try the ageHandler handler again
		return handlers.NextConversationState(NODENAME, NODETOKEN)
	}

	err := ctx.Reply(fmt.Sprintf("Token %s\n 请输入绑定的域名?", inputToken), &telebot.SendOptions{})
	if err != nil {
		return fmt.Errorf("failed to send ageHandler message: %w", err)
	}
	return handlers.NextConversationStateWithData(NODETOKEN, NODEDOMAIN, inputToken)
}

func nodeDomainHandler(ctx telebot.IContext) error {
	inputDomain := ctx.Message().Text
	nodeName := ""
	nodeToken := ""
	err := ctx.Reply(fmt.Sprintf("节点名称:%s\nToken: %s\n域名: %s", nodeName, nodeToken, html.EscapeString(inputDomain)), &telebot.SendOptions{})
	if err != nil {
		return fmt.Errorf("failed to send name message: %w", err)
	}

	// 完成
	return handlers.EndConversation()
}
