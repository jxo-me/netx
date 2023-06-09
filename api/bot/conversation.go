package bot

import (
	"fmt"
	telebot "github.com/jxo-me/gfbot"
	"github.com/jxo-me/gfbot/handlers"
	"html"
	"strconv"
)

const (
	NAME     = "nameHandler"
	LOCATION = "locationHandler"
	AGE      = "ageHandler"
)

func NewConversation(entry, cancel string) handlers.Conversation {
	return handlers.NewConversation(
		entry,
		telebot.HandlerFunc(startHandler), // 入口
		map[string][]telebot.IHandler{
			NAME:     {telebot.HandlerFunc(nameHandler)},
			AGE:      {telebot.HandlerFunc(ageHandler)},
			LOCATION: {telebot.HandlerFunc(locationHandler)},
		}, // states状态
		&handlers.ConversationOpts{
			ExitName:     cancel,
			ExitHandler:  telebot.HandlerFunc(cancelHandler),
			AllowReEntry: true,
		}, // options
	)
}

func startHandler(ctx telebot.IContext) error {
	err := ctx.Reply(fmt.Sprintf("Hello, I'm @%s.\nWhat is your name?.", ctx.Bot().Me.Username), &telebot.SendOptions{})
	if err != nil {
		return fmt.Errorf("failed to send start message: %w", err)
	}
	// 设置当前用户下一个入口
	return handlers.NextConversationState(handlers.Entry, NAME)
}

// cancelHandler cancels the conversation.
func cancelHandler(ctx telebot.IContext) error {
	err := ctx.Reply("Oh, goodbye!", &telebot.SendOptions{})
	if err != nil {
		return fmt.Errorf("failed to send cancelHandler message: %w", err)
	}
	//return handlers.EndConversation()
	return nil
}

// nameHandler gets the user's nameHandler
func nameHandler(ctx telebot.IContext) error {
	inputName := ctx.Message().Text
	err := ctx.Reply(fmt.Sprintf("Nice to meet you, %s!\n\nAnd how old are you?", html.EscapeString(inputName)), &telebot.SendOptions{})
	if err != nil {
		return fmt.Errorf("failed to send name message: %w", err)
	}
	return handlers.NextConversationState(NAME, AGE)
}

// ageHandler gets the user's ageHandler
func ageHandler(ctx telebot.IContext) error {
	inputAge := ctx.Message().Text
	ageNumber, err := strconv.ParseInt(inputAge, 10, 64)
	if err != nil {
		// If the number is not valid, try again!
		ctx.Reply(fmt.Sprintf("This doesn't seem to be a number. Could you repeat?"), &telebot.SendOptions{})
		// We try the ageHandler handler again
		return handlers.NextConversationState(NAME, AGE)
	}

	err = ctx.Reply(fmt.Sprintf("age %d\n What's your location?", ageNumber), &telebot.SendOptions{})
	if err != nil {
		return fmt.Errorf("failed to send ageHandler message: %w", err)
	}
	return handlers.NextConversationState(AGE, LOCATION)
}

func locationHandler(ctx telebot.IContext) error {
	inputLocation := ctx.Message().Text
	err := ctx.Reply(fmt.Sprintf("Full name: name\nAge: {age}\nLocation: %s", html.EscapeString(inputLocation)), &telebot.SendOptions{})
	if err != nil {
		return fmt.Errorf("failed to send name message: %w", err)
	}

	// 完成
	return handlers.EndConversation()
}
