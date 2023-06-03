package bot

import telebot "github.com/jxo-me/gfbot"

var (
	Event = hEvent{}
)

type (
	hEvent struct{}
)

func (h *hEvent) OnText(c telebot.Context) error {
	return c.Send("OnText")
}

func (h *hEvent) OnCallback(c telebot.Context) error {
	return c.Send("OnCallback")
}

func (h *hEvent) OnUserJoined(c telebot.Context) error {
	return c.Send("OnUserJoined")
}
