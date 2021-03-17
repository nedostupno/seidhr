package delivery

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Обявим переменные для конфигов сообщений в самом начале,
// что бы не приходилось повторять их объявление в каждой функции
var msg, newKeyboard, newText tgbotapi.Chattable

// SendMessage - принимаем на вход интерфейсы описывающие наше сообщение,
//а так же объект бота и выполняем отправку
func SendMessage(msg tgbotapi.Chattable, newKeyboard tgbotapi.Chattable, newText tgbotapi.Chattable, bot tgbotapi.BotAPI) error {

	if msg != nil {
		if _, err := bot.Send(msg); err != nil {
			return err
		}
	}
	if newKeyboard != nil {
		if _, err := bot.Send(newKeyboard); err != nil {
			return err
		}
	}
	if newText != nil {
		if _, err := bot.Send(newText); err != nil {
			return err
		}
	}
	return nil
}

func Start(message *tgbotapi.Message, bot tgbotapi.BotAPI) error {

	msgConf := tgbotapi.NewMessage(message.Chat.ID, "Добрый день, мой ярл, рад видеть вас живым!")

	msg = msgConf
	newKeyboard = nil
	newText = nil

	if err := SendMessage(msg, newKeyboard, newText, bot); err != nil {
		return err
	}

	return nil
}

func Default(message *tgbotapi.Message, bot tgbotapi.BotAPI) error {

	msgConf := tgbotapi.NewMessage(message.Chat.ID, "К сожалению, я так не умею")

	msg = msgConf
	newKeyboard = nil
	newText = nil

	if err := SendMessage(msg, newKeyboard, newText, bot); err != nil {
		return err
	}

	return nil
}
