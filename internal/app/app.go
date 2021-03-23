package app

import (
	"log"
	"net/http"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/lib/pq"
	"github.com/nedostupno/seidhr/internal/config"
	"github.com/nedostupno/seidhr/internal/delivery"
	"github.com/nedostupno/seidhr/internal/repository"
	"github.com/nedostupno/seidhr/internal/service"
)

// checkTime - следит за временем и в момент, когда время становится равным 11 часам пишет в канал
// Который будет считан функцией CyclicMedSearch, после чего она будет запущена.
func checkTime(c chan bool) {
	for {
		hour := time.Now().Hour()

		if hour == 11 {
			c <- true
			time.Sleep(23 * time.Hour)
		}
		time.Sleep(20 * time.Minute)
	}
}

// Run - данная функция стартует наш сервер, инициализирует конфиги и базу данных
func Run() {

	// Инициализируем конфиг
	cfg, err := config.Init()
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	// Инициализируем Базу данных
	db, err := repository.NewPostgreDB(cfg)
	if err != nil {
		log.Fatalf("%+v", err)
	}

	// Инициализируем бота
	if bot, err := tgbotapi.NewBotAPI(cfg.Bot.APIToken); err != nil {
		log.Fatalf("%+v", err)
	} else {

		bot.Debug = true

		// Поулчаем инфу о состоянии нашего вебхука
		// Выводим в консоль последнюю возникшую ошибку
		info, err := bot.GetWebhookInfo()
		if err != nil {
			log.Fatalf("%+v", err)
		}
		if info.LastErrorDate != 0 {
			log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
		}

		//Начинаем слушать на 8443 порту
		updates := bot.ListenForWebhook("/" + bot.Token)
		go http.ListenAndServeTLS(":8443", cfg.SSL.Fullchain, cfg.SSL.Privkey, nil)

		// Получаем обновления от телеграма,
		// в зависимости от типа полученного сообщения используем разные обработчики
		//for update := range updates {
		//	//
		//}

		// Хэндлер будет обрабатывать запрос от пользователя, после чего
		// в необходимых случаях вызывать метод сервиса,
		// который в свою очеред будет вызывать методы пакета repository
		//
		repos := repository.NewRepository(db)
		services := service.NewService(repos)
		handlers := delivery.NewHandler(services, bot)

		// Проверяем заполнена ли наша база значениями полученными из файла drugs.txt
		// если нет, то заполняем.
		isExist, err := services.Medicaments.IsMedListExist()
		if err != nil {
			log.Fatalf("%+v", err)
		}

		if isExist == false {

			meds, err := services.Medicaments.ReadFileWithMeds()
			if err != nil {
				log.Fatalf("%+v", err)
			}

			err = services.Medicaments.InitMedList(meds)
			if err != nil {
				log.Fatalf("%+v", err)
			}
		}

		// Создаем канал необходимый для работы функций отвечающих за ежедневную проверку
		// наличия лекарств в аптеке.

		c := make(chan bool)

		go checkTime(c)
		// Запускаем в горутине функцию, которая читает из канала
		// и и после чтения производит проверку наличия лекрств,
		// на которые оформлены подписки, и если лекрство появилось,
		// то уведомляем пользователей с подпиской об этом
		go delivery.CyclicMedSearch(bot, handlers, c)

		// Вызываем метод объекта Handler, с помощью которого мы обрабатываем
		// сообщения, поступающие в канал updates
		handlers.HandleUpdate(updates)
	}

}
