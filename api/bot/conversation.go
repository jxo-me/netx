package bot

import (
	"fmt"
	telebot "github.com/jxo-me/gfbot"
	"github.com/jxo-me/gfbot/handlers"
	"html"
	"strconv"
)

func NewConversation() handlers.Conversation {
	return handlers.NewConversation(
		[]telebot.IHandler{telebot.HandlerFunc(start)}, // 入口
		map[string][]telebot.IHandler{
			NAME: {telebot.HandlerFunc(name)},
			AGE:  {telebot.HandlerFunc(age)},
		}, // states状态
		&handlers.ConversationOpts{
			Exits:        []telebot.IHandler{telebot.HandlerFunc(cancel)},
			AllowReEntry: true,
		}, // options
	)
}

// age gets the user's age
func age(ctx telebot.IContext) error {
	inputAge := ctx.Message().Text
	ageNumber, err := strconv.ParseInt(inputAge, 10, 64)
	if err != nil {
		// If the number is not valid, try again!
		ctx.Reply(fmt.Sprintf("This doesn't seem to be a number. Could you repeat?"), &telebot.SendOptions{})
		// We try the age handler again
		return handlers.NextConversationState(AGE)
	}

	err = ctx.Reply(fmt.Sprintf("Ah, you're %d years old!", ageNumber), &telebot.SendOptions{})
	if err != nil {
		return fmt.Errorf("failed to send age message: %w", err)
	}
	return handlers.EndConversation()
}

// cancel cancels the conversation.
func cancel(ctx telebot.IContext) error {
	err := ctx.Reply("Oh, goodbye!", &telebot.SendOptions{})
	if err != nil {
		return fmt.Errorf("failed to send cancel message: %w", err)
	}
	return handlers.EndConversation()
}

// name gets the user's name
func name(ctx telebot.IContext) error {
	fmt.Println("888888888888888888888888888888888:", ctx.Message())
	inputName := ctx.Message().Text
	err := ctx.Reply(fmt.Sprintf("Nice to meet you, %s!\n\nAnd how old are you?", html.EscapeString(inputName)), &telebot.SendOptions{})
	if err != nil {
		return fmt.Errorf("failed to send name message: %w", err)
	}
	return handlers.NextConversationState(AGE)
}

func start(ctx telebot.IContext) error {
	err := ctx.Reply(fmt.Sprintf("Hello, I'm @%s.\nWhat is your name?.", ctx.Bot().Me.Username), &telebot.SendOptions{})
	if err != nil {
		return fmt.Errorf("failed to send start message: %w", err)
	}
	// 设置当前用户下一个入口
	return handlers.NextConversationState(NAME)
}
