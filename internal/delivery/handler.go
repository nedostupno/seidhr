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
	if message.IsCommand() {
		cmd := message.Command()
		switch cmd {
		case "start":
			if err := Start(message, *h.bot); err != nil {
				return err
			}
		default:
			if err := Default(message, *h.bot); err != nil {
				return err
			}
		}
	}
	return nil
}

// HandleCallbackQuery - обрабатываем сообщения,
// которые вызваны нажатием пользователя на кнопки бота
func (h *Handler) HandleCallbackQuery(callbackQuery *tgbotapi.CallbackQuery) error {
	//
	return nil
}
