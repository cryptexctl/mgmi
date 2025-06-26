package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	maxbot "github.com/max-messenger/max-bot-api-client-go"
	"github.com/max-messenger/max-bot-api-client-go/schemes"
)

func main() {
	token := os.Getenv("MAX_BOT_TOKEN")
	if token == "" {
		log.Fatal("MAX_BOT_TOKEN не установлен")
	}

	api := maxbot.New(token)

	info, err := api.Bots.GetBot()
	if err != nil {
		log.Fatal("Не удалось подключиться к боту:", err)
	}
	log.Printf("Бот подключен: %s", info.Name)

	ch := make(chan interface{})

	http.HandleFunc("/webhook", api.GetHandler(ch))

	go func() {
		for {
			upd := <-ch
			log.Printf("Получено обновление: %#v", upd)

			switch upd := upd.(type) {
			case *schemes.MessageCreatedUpdate:
				if upd.Message.Body.Text == "/start" {
					chatID := upd.Message.Recipient.ChatId
					userID := upd.Message.Sender.UserId

					responseText := fmt.Sprintf("Chat ID: %d\nUser ID: %d", chatID, userID)

					_, err := api.Messages.Send(
						maxbot.NewMessage().
							SetChat(chatID).
							SetText(responseText),
					)

					if err != nil {
						log.Printf("Ошибка отправки сообщения: %v", err)
					} else {
						log.Printf("Отправлено в чат %d для пользователя %d", chatID, userID)
					}
				}
			default:
				log.Printf("Неизвестный тип обновления: %#v", upd)
			}
		}
	}()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Бот запущен на порту %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
