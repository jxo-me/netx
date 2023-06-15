package bot

import (
	"fmt"
	telebot "github.com/jxo-me/gfbot"
	"github.com/jxo-me/gfbot/handlers"
)

const (
	MenuSet = "menuSet"
)

func SetMenuConversation(entry, cancel string) handlers.Conversation {
	return handlers.NewConversation(
		entry,
		telebot.HandlerFunc(startSetMenuHandler), // 入口
		map[string][]telebot.IHandler{
			MenuSet: {telebot.HandlerFunc(setMenuHandler)},
		}, // states状态
		&handlers.ConversationOpts{
			ExitName:     cancel,
			ExitHandler:  telebot.HandlerFunc(cancelSetMenuHandler),
			AllowReEntry: true,
		}, // options
	)
}

func startSetMenuHandler(ctx telebot.IContext) error {
	menu := &telebot.ReplyMarkup{ResizeKeyboard: true, OneTimeKeyboard: true}
	disable := menu.Text("Disable")
	enable := menu.Text("Enable")
	menu.Reply(menu.Row(disable, enable))
	msg := `'Enable' - your bot will only receive messages that either start with the '/' symbol or mention the bot by username.
'Disable' - your bot will receive all messages that people send to groups.
Current status is: DISABLED`
	_ = ctx.Send(msg, &telebot.SendOptions{ReplyMarkup: menu})
	// 设置当前用户下一个入口
	return handlers.NextConversationState(MenuSet)
}

func setMenuHandler(ctx telebot.IContext) error {
	str := ctx.Message().Text
	fmt.Println("set Menu status str:", str)
	msg := `Success! The new status is: ENABLED. /help`
	_ = ctx.Send(msg, &telebot.SendOptions{ReplyMarkup: &telebot.ReplyMarkup{RemoveKeyboard: true}}) // 移除自定义键盘
	return handlers.EndConversation()
}

func cancelSetMenuHandler(ctx telebot.IContext) error {
	err := ctx.Reply("当前 /menu 操作已被取消。 还有什么我可以为你做的吗？", &telebot.SendOptions{})
	if err != nil {
		return fmt.Errorf("failed to send cancelHandler message: %w", err)
	}
	//return handlers.EndConversation()
	return nil
}
