package delivery

import (
	"strconv"

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
	if newText != nil {
		if _, err := bot.Send(newText); err != nil {
			return err
		}
	}
	if newKeyboard != nil {
		if _, err := bot.Send(newKeyboard); err != nil {
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

		msgConf := tgbotapi.NewMessage(chatID, "Добро пожаловать.\n\nЧто бы проверить наличие необходимого вам льготного лекарства в аптеках Санкт-Петербурга, просто нажмите на кнопку и введите его навание.\n\nВ случае, если необходимого вам лекарства сейчас нигде нет, вы можете подписаться на него и мы сообщим вам, как только оно появится.\n\nДля получения информационной справки используте команду /help Приятного использования!")
		msgConf.ReplyMarkup = keyboards.HomeKeyboard

		msg = msgConf
		newKeyboard = nil
		newText = nil

		if err := SendMessage(msg, newKeyboard, newText, bot); err != nil {
			return err
		}

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

		msgConf := tgbotapi.NewMessage(chatID, "Что бы вы хотели?")
		msgConf.ReplyMarkup = keyboards.HomeKeyboard

		msg = msgConf
		newKeyboard = nil
		newText = nil

		if err := SendMessage(msg, newKeyboard, newText, bot); err != nil {
			return err
		}

		return nil
	}

	msgConf := tgbotapi.NewMessage(chatID, "Что бы вы хотели?")
	msgConf.ReplyMarkup = keyboards.HomeWithSubKeyboard

	msg = msgConf
	newKeyboard = nil
	newText = nil

	if err := SendMessage(msg, newKeyboard, newText, bot); err != nil {
		return err
	}

	return nil
}

func DefaultCommand(message *tgbotapi.Message, bot tgbotapi.BotAPI) error {
	chatID := message.Chat.ID
	msgConf := tgbotapi.NewMessage(chatID, "К сожалению, я так не умею")

	msg = msgConf
	newKeyboard = nil
	newText = nil

	if err := SendMessage(msg, newKeyboard, newText, bot); err != nil {
		return err
	}

	return nil
}

func Help(message *tgbotapi.Message, bot tgbotapi.BotAPI) error {
	chatID := message.Chat.ID

	msg = tgbotapi.NewMessage(chatID, "---Необольшая справка---\n\nС помощью данного бота вы можете проверить наличие льготных лекарств в аптеках Санкт-Петербурга, а так же подписаться на необходимые вам лекарства, и получать уведомления, как только они появятся в аптеках.\n\nЧто бы подписаться на какое-либо лекарство, вам необходимо нажать на кнопку <Проверить лекарство> и ввести его название, после чего в появившемся сообщении вы увидите всю информацию о нем, а так же кнопку <Подписаться>, если вы, конечно, уже не подписаны на него\n\nПосле того, как вы подписались на ваше первое лекарство, в главном меню появится кнопка <Подписки>, нажав на которую вы увидите все ваши подписки, узнать наличие, а так же отменить подписку.")
	newKeyboard = nil
	newText = nil

	if err := SendMessage(msg, newKeyboard, newText, bot); err != nil {
		return err
	}

	return nil
}

func DefaultMsg(message *tgbotapi.Message, bot tgbotapi.BotAPI, h *Handler) error {

	tguserID := message.From.ID
	chatID := message.Chat.ID

	// Cмотрим состояние его подписок
	isSubscribe, err := h.services.Users.IsHasSubsriptions(tguserID)
	if err != nil {
		return err
	}

	// Если у пользователя нет подписок, то выдаем ему клавиатуру
	// без кнопки просмотра подписок
	if isSubscribe == false {

		msgConf := tgbotapi.NewMessage(chatID, "Простите, я могу обрабатывать только названия лекарств, если вы хотите хотите найти лекарство, то нажмите на кнопку **Проверить лекарство**")
		msgConf.ReplyMarkup = keyboards.HomeKeyboard

		msg = msgConf
		newKeyboard = nil
		newText = nil

		if err := SendMessage(msg, newKeyboard, newText, bot); err != nil {
			return err
		}

		return nil
	}

	msgConf := tgbotapi.NewMessage(chatID, "Простите, я могу обрабатывать только названия лекарств, если вы хотите хотите найти лекарство, то нажмите на кнопку **Проверить лекарство**")
	msgConf.ReplyMarkup = keyboards.HomeWithSubKeyboard

	msg = msgConf
	newKeyboard = nil
	newText = nil

	if err := SendMessage(msg, newKeyboard, newText, bot); err != nil {
		return err
	}
	return nil
}

func SearchMed(callbackQuery *tgbotapi.CallbackQuery, bot tgbotapi.BotAPI, h *Handler) error {
	tguserID := callbackQuery.From.ID
	chatID := callbackQuery.Message.Chat.ID
	msgID := callbackQuery.Message.MessageID
	state := "SearchMed"

	err := h.services.Users.ChangeState(tguserID, state)
	if err != nil {
		return err
	}

	msg = nil
	newKeyboard = tgbotapi.NewEditMessageReplyMarkup(chatID, msgID, keyboards.MedSearchKeyboard)
	newText = tgbotapi.NewEditMessageText(chatID, msgID, "Введите название лекарства:")

	if err := SendMessage(msg, newKeyboard, newText, bot); err != nil {
		return err
	}

	return nil
}

func SearchMedAct(message *tgbotapi.Message, bot tgbotapi.BotAPI, h *Handler) error {
	medTitle := message.Text
	tguserID := message.From.ID
	chatID := message.Chat.ID

	// Проверяем наличие данного лекарства в базе данных льготных лекарств
	isExist, err := h.services.Medicaments.IsExist(medTitle)
	if err != nil {
		return err
	}

	if isExist == false {
		msgConf := tgbotapi.NewMessage(chatID, "Простите, но кажется вы неправильно написали название, либо это лекарство не льготное. Попробуйте еще раз:")
		msgConf.ReplyMarkup = keyboards.MedSearchKeyboard

		msg = msgConf
		newKeyboard = nil
		newText = nil

		if err := SendMessage(msg, newKeyboard, newText, bot); err != nil {
			return err
		}

		return nil
	}

	trueName, err := h.services.Medicaments.GetTrueName(medTitle)
	if err != nil {
		return err
	}

	// Отправляем запрос
	medResp, err := h.services.Medicaments.ReqMedInfo(trueName)
	if err != nil {
		return err
	}

	// Находим id нашго лекарства
	medicamentID, err := h.services.Medicaments.GetID(trueName)
	if err != nil {
		return err
	}

	err = h.services.Users.ChangeSelectedMed(medicamentID, tguserID)
	if err != nil {
		return err
	}

	// Проверяем подписан ли пользователь на это лекарство, что бы решить
	// какую клавиатуру и текст необходимо отобразить
	isSubscribe, err := h.services.Users.IsSubToThisMed(tguserID, medicamentID)
	if err != nil {
		return err
	}

	if isSubscribe == true {
		// Проверяем полученный json на наличе информации об ошибке.
		// Так как перед отправкой запроса мы проверяем наличие лекарства в нашей бд,
		// где хранится список лекарств доступных по льготе, то вариант с неправильным написанием
		// или вводом чего-то вообще неподходящего или несуществующего
		// Значит, ошибка всегда будет означать то, что лекарства сейчас нет в доступе
		isErr := h.services.Medicaments.IsErrExistInJSON(medResp)
		if isErr == true {
			msgConf := tgbotapi.NewMessage(chatID, "К сожалению данного лекарства сейчас нет ни в одной аптеке, но так как вы подписаны, мы уведомим вас, как только оно появится в аптеках")
			msgConf.ReplyMarkup = keyboards.ViewMedKeyboard

			msg = msgConf
			newKeyboard = nil
			newText = nil

			if err := SendMessage(msg, newKeyboard, newText, bot); err != nil {
				return err
			}

			return nil
		}

		// Парсим json и компануем текст сообщения
		msgText := h.services.Medicaments.ParseJSON(medResp)

		msgConf := tgbotapi.NewMessage(chatID, msgText)
		msgConf.ReplyMarkup = keyboards.ViewMedKeyboard

		msg = msgConf
		newKeyboard = nil
		newText = nil

		if err := SendMessage(msg, newKeyboard, newText, bot); err != nil {
			return err
		}

		return nil

	}

	// Опять Проверяем полученный json на наличе информации об ошибке
	isErr := h.services.Medicaments.IsErrExistInJSON(medResp)
	if isErr == true {

		msgConf := tgbotapi.NewMessage(chatID, "К сожалению данного лекарства сейчас нет ни в одной аптеке, но если хотите, вы можете подписаться и мы уведомим вас, как только оно появится")
		msgConf.ReplyMarkup = keyboards.ViewMedWithSubKeyboard

		msg = msgConf
		newKeyboard = nil
		newText = nil

		if err := SendMessage(msg, newKeyboard, newText, bot); err != nil {
			return err
		}

		return nil
	}

	// Парсим json и компануем текст сообщения
	msgText := h.services.Medicaments.ParseJSON(medResp)

	msgConf := tgbotapi.NewMessage(chatID, msgText)
	msgConf.ReplyMarkup = keyboards.ViewMedWithSubKeyboard

	msg = msgConf
	newKeyboard = nil
	newText = nil

	if err := SendMessage(msg, newKeyboard, newText, bot); err != nil {
		return err
	}

	return nil
}

// BackToHome - возвращает пользователя на домашнюю страницу
func BackToHome(callbackQuery *tgbotapi.CallbackQuery, bot tgbotapi.BotAPI, h *Handler) error {

	tguserID := callbackQuery.From.ID
	chatID := callbackQuery.Message.Chat.ID
	msgID := callbackQuery.Message.MessageID

	isSubscribe, err := h.services.IsHasSubsriptions(tguserID)
	if err != nil {
		return err
	}

	state, err := h.services.GetState(tguserID)
	if err != nil {
		return err
	}

	// Если пользователь просматривал подписки, нажал на лекарство и потом нажал кнопку <backToHome>,
	// то обрабатываем другим способом и выводим ему все его подписки, вместо домашней страницы
	if state == "ViewSubMed" {

		h.services.Users.ChangeState(tguserID, "Home")

		//msg, newKeyboard, newText, err := ListSubscriptions(callbackQuery)
		//if err != nil {
		//	return nil, nil, nil, err
		//}
		return nil
	}

	h.services.Users.ChangeState(tguserID, "Home")
	// Если у пользователя нет подписок, то выдаем ему клавиатуру
	// без кнопки просмотра подписок
	if isSubscribe == false {
		msg = nil

		newKeyboard = tgbotapi.NewEditMessageReplyMarkup(chatID, msgID, keyboards.HomeKeyboard)

		newText = tgbotapi.NewEditMessageText(chatID, msgID, "Что бы вы хотели?")

		if err := SendMessage(msg, newKeyboard, newText, bot); err != nil {
			return err
		}

		return nil
	}

	msg = nil

	newKeyboard = tgbotapi.NewEditMessageReplyMarkup(chatID, msgID, keyboards.HomeWithSubKeyboard)

	newText = tgbotapi.NewEditMessageText(chatID, msgID, "Что бы вы хотели?")

	if err := SendMessage(msg, newKeyboard, newText, bot); err != nil {
		return err
	}

	return nil
}

// ListSubscriptions - отправляет пользователю сообщение с информацией о всех его подписках
func ListSubscriptions(callbackQuery *tgbotapi.CallbackQuery, bot tgbotapi.BotAPI, h *Handler) error {

	tguserID := callbackQuery.From.ID
	chatID := callbackQuery.Message.Chat.ID
	msgID := callbackQuery.Message.MessageID

	subscriptions, err := h.services.GetSubscriptions(tguserID)
	if err != nil {
		return err
	}

	subKeyboard := keyboards.CreateKeyboarWithUserSubscriptions(subscriptions)

	msg = nil

	newKeyboard = tgbotapi.NewEditMessageReplyMarkup(chatID, msgID, subKeyboard)

	newText = tgbotapi.NewEditMessageText(chatID, msgID, "Ваши подписки:")

	if err := SendMessage(msg, newKeyboard, newText, bot); err != nil {
		return err
	}

	return nil
}

// Subscribe - оформляет подписку на лекарство для данного пользователя
func Subscribe(callbackQuery *tgbotapi.CallbackQuery, bot tgbotapi.BotAPI, h *Handler) error {

	tguserID := callbackQuery.From.ID
	chatID := callbackQuery.Message.Chat.ID
	msgID := callbackQuery.Message.MessageID

	medicamentID, err := h.services.Users.GetSelectedMed(tguserID)
	if err != nil {
		return err
	}

	err = h.services.Users.Subscribe(tguserID, medicamentID)
	if err != nil {
		return err
	}

	h.services.Users.ChangeState(tguserID, "Home")

	msg = nil

	newKeyboard = tgbotapi.NewEditMessageReplyMarkup(chatID, msgID, keyboards.HomeWithSubKeyboard)

	newText = tgbotapi.NewEditMessageText(chatID, msgID, "Поздравляю, подписка на лекарство успешно оформлена. Теперь вы первым узнаете о появлении данного лекарства в аптеках нашего города.\n\nХотите еще что-нибудь?")

	if err := SendMessage(msg, newKeyboard, newText, bot); err != nil {
		return err
	}

	return nil
}

func Unsubscribe(callbackQuery *tgbotapi.CallbackQuery, bot tgbotapi.BotAPI, h *Handler) error {

	tguserID := callbackQuery.From.ID
	chatID := callbackQuery.Message.Chat.ID
	msgID := callbackQuery.Message.MessageID

	medicamentID, err := h.services.Users.GetSelectedMed(tguserID)
	if err != nil {
		return err
	}

	err = h.services.Unsubscribe(tguserID, medicamentID)
	if err != nil {
		return err
	}

	h.services.Users.ChangeState(tguserID, "Home")

	isSub, err := h.services.Users.IsHasSubsriptions(tguserID)
	if err != nil {
		return err
	}

	if isSub == true {
		msg = nil
		newKeyboard = tgbotapi.NewEditMessageReplyMarkup(chatID, msgID, keyboards.HomeWithSubKeyboard)

		newText = tgbotapi.NewEditMessageText(chatID, msgID, "Поздравляю, подписка отменена.\n\nХотите еще что-нибудь?")

		if err := SendMessage(msg, newKeyboard, newText, bot); err != nil {
			return err
		}

		return nil
	}

	msg = nil

	newKeyboard = tgbotapi.NewEditMessageReplyMarkup(chatID, msgID, keyboards.HomeKeyboard)

	newText = tgbotapi.NewEditMessageText(chatID, msgID, "Поздравляю, подписка отменена.\n\nХотите еще что-нибудь?")

	if err := SendMessage(msg, newKeyboard, newText, bot); err != nil {
		return err
	}

	return nil
}

// InterceptMedicament - перехватывает id лекарств
func InterceptMedicament(callbackQuery *tgbotapi.CallbackQuery, bot tgbotapi.BotAPI, h *Handler) error {

	tguserID := callbackQuery.From.ID
	chatID := callbackQuery.Message.Chat.ID
	msgID := callbackQuery.Message.MessageID

	h.services.Users.ChangeState(tguserID, "ViewSubMed")

	medicamentID, err := strconv.Atoi(callbackQuery.Data)
	if err != nil {
		return err
	}

	title, err := h.services.Medicaments.GetTitle(medicamentID)
	if err != nil {
		return err
	}

	err = h.services.Users.ChangeSelectedMed(medicamentID, tguserID)
	if err != nil {
		return err
	}

	// Отправляем запрос
	medResp, err := h.services.Medicaments.ReqMedInfo(title)
	if err != nil {
		return err
	}

	isErr := h.services.Medicaments.IsErrExistInJSON(medResp)
	if isErr == true {

		msg = nil

		newKeyboard = tgbotapi.NewEditMessageReplyMarkup(chatID, msgID, keyboards.ViewMedKeyboard)

		newText = tgbotapi.NewEditMessageText(chatID, msgID, "К сожалению данного лекарства сейчас нет ни в одной аптеке, но так как вы подписаны, мы уведомим вас, как только оно появится в аптеках")

		if err := SendMessage(msg, newKeyboard, newText, bot); err != nil {
			return err
		}

		return nil
	}

	// Парсим json и компануем текст сообщения
	msgText := h.services.Medicaments.ParseJSON(medResp)

	msg = nil

	newKeyboard = tgbotapi.NewEditMessageReplyMarkup(chatID, msgID, keyboards.ViewMedKeyboard)

	newText = tgbotapi.NewEditMessageText(chatID, msgID, msgText)

	if err := SendMessage(msg, newKeyboard, newText, bot); err != nil {
		return err
	}

	return nil
}
