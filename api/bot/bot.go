package bot

import telebot "github.com/jxo-me/gfbot"

var (
	insBotRouter = Routers{
		List: map[string]telebot.HandlerFunc{
			telebot.OnText:       Event.OnText,
			telebot.OnCallback:   Event.OnCallback,
			telebot.OnUserJoined: Event.OnUserJoined,
		},
		Btns: map[*telebot.Btn]telebot.HandlerFunc{
			//&bot.BtnBetting:       bot.Event.OnBtnBetting,
		},
	}
)

type Routers struct {
	List map[string]telebot.HandlerFunc
	Btns map[*telebot.Btn]telebot.HandlerFunc
}

func Router() *Routers {
	return &insBotRouter
}
