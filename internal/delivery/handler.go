package delivery

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/nedostupno/seidhr/internal/service"
)

// Handler - структура первичного обработчика, которая
// включает в себя структуру сервиса
type Handler struct {
	services *service.Service
	bot      *tgbotapi.BotAPI
}

// Создаем новый хендлер, принимая в него указатель на сервис,
// что бы в будущем иметь возможность обращаться к его методам
//
func NewHandler(services *service.Service, bot *tgbotapi.BotAPI) *Handler {
	return &Handler{services: services, bot: bot}
}

// HandleUpdate - перехватываем сообщения из канала и обрабатываем их и
// в зависимости от содержимого вызываем методы пакета service,
// делегируя дальнейшую обработку
func (h *Handler) HandleUpdate(updates tgbotapi.UpdatesChannel) error {
	for update := range updates {
		if update.Message != nil {

			err := h.HandleMessage(update.Message)
			if err != nil {
				return err
			}

		} else if update.CallbackQuery != nil {

			err := h.HandleCallbackQuery(update.CallbackQuery)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

// HandleMessage - обрабатывем объекты обычных сообщений
func (h *Handler) HandleMessage(message *tgbotapi.Message) error {

	// Если нам пришла команда, то сверяем ее с доступными командами,
	// если ни одного совпадения нет, то отвечаем Default
	// -------------------------------------------------------------
	// Если это не команда, то значит нам пришло обычное сообщение
	// Данный бот в качестве обычных сообщений принимает только названия лекарств,
	// Поэтому выполняем проверку лекарства
	if message.IsCommand() {
		cmd := message.Command()
		switch cmd {
		case "start":
			if err := Start(message, *h.bot, h); err != nil {
				return err
			}
		case "help":
			if err := Help(message, *h.bot); err != nil {
				return err
			}
		default:
			if err := DefaultCommand(message, *h.bot); err != nil {
				return err
			}
		}
	} else {
		tguserID := message.From.ID

		state, err := h.services.Users.GetState(tguserID)
		if err != nil {
			return err
		}

		switch state {
		case "SearchMed":
			if err := SearchMedAct(message, *h.bot, h); err != nil {
				return err
			}

			return nil
		default:
			if err := DefaultMsg(message, *h.bot, h); err != nil {
				return err
			}
		}
	}
	return nil
}

// HandleCallbackQuery - обрабатываем сообщения,
// которые вызваны нажатием пользователя на кнопки бота
func (h *Handler) HandleCallbackQuery(callbackQuery *tgbotapi.CallbackQuery) error {
	//Перехватываем нажатие на кнопку <Проверить лекарство> и
	// с помощью SearchMed() меняем state пользователя на "SearchMed"
	// и предлагаем ввести название лекарства

	switch callbackQuery.Data {
	case "searchMed":
		if err := SearchMed(callbackQuery, *h.bot, h); err != nil {
			return err
		}
	case "backToHome":
		if err := BackToHome(callbackQuery, *h.bot, h); err != nil {
			return err
		}
	case "subscribe":
		if err := Subscribe(callbackQuery, *h.bot, h); err != nil {
			return err
		}
	case "unsubscribe":
		if err := Unsubscribe(callbackQuery, *h.bot, h); err != nil {
			return err
		}
	case "lsSub":
		if err := ListSubscriptions(callbackQuery, *h.bot, h); err != nil {
			return err
		}
	}
	return nil
}
