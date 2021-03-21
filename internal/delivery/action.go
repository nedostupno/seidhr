package delivery

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/nedostupno/seidhr/keyboards"
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

func Start(message *tgbotapi.Message, bot tgbotapi.BotAPI, h *Handler) error {

	tguserID := message.From.ID
	chatID := message.Chat.ID
	isExist, err := h.services.Users.Check(tguserID)
	if err != nil {
		return err
	}

	// Проверяем первый ли раз пользователь обращается к боту
	if isExist == false {
		err := h.services.Users.Create(tguserID, chatID)
		if err != nil {
			return err
		}

		msgConf := tgbotapi.NewMessage(message.Chat.ID, "Добро пожаловать.\n\nЧто бы проверить наличие необходимого вам льготного лекарства в аптеках Санкт-Петербурга, просто нажмите на кнопку и введите его навание.\n\nВ случае, если необходимого вам лекарства сейчас нигде нет, вы можете подписаться на него и мы сообщим вам, как только оно появится.\n\nДля получения информационной справки используте команду /help Приятного использования!")
		msgConf.ReplyMarkup = keyboards.HomeKeyboard

		msg = msgConf
		newKeyboard = nil
		newText = nil

		SendMessage(msg, newKeyboard, newText, bot)
		return nil
	}

	// Если пользователь уже взаимодействовал с ботом,
	// смотрим состояние его подписок
	isSubscribe, err := h.services.Users.IsHasSubsriptions(tguserID)
	if err != nil {
		return err
	}

	// Если у пользователя нет подписок, то выдаем ему клавиатуру
	// без кнопки просмотра подписок
	if isSubscribe == false {

		msgConf := tgbotapi.NewMessage(message.Chat.ID, "Что бы вы хотели?")
		msgConf.ReplyMarkup = keyboards.HomeKeyboard

		msg = msgConf
		newKeyboard = nil
		newText = nil

		SendMessage(msg, newKeyboard, newText, bot)
		return nil
	}

	msgConf := tgbotapi.NewMessage(message.Chat.ID, "Что бы вы хотели?")
	msgConf.ReplyMarkup = keyboards.HomeWithSubKeyboard

	msg = msgConf
	newKeyboard = nil
	newText = nil

	SendMessage(msg, newKeyboard, newText, bot)

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
