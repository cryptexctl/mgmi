package main

import (
	"fmt"
	"log"
	"os"

	maxbot "github.com/max-messenger/max-bot-api-client-go"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Использование: go run utils/webhook_setup.go <webhook_url>")
		fmt.Println("Пример: go run utils/webhook_setup.go https://mybot.ngrok.io/webhook")
		os.Exit(1)
	}

	token := os.Getenv("MAX_BOT_TOKEN")
	if token == "" {
		log.Fatal("MAX_BOT_TOKEN не установлен")
	}

	webhookURL := os.Args[1]
	api := maxbot.New(token)

	subs, err := api.Subscriptions.GetSubscriptions()
	if err != nil {
		log.Fatal("Ошибка получения подписок:", err)
	}

	for _, s := range subs.Subscriptions {
		_, err := api.Subscriptions.Unsubscribe(s.Url)
		if err != nil {
			log.Printf("Ошибка отписки от %s: %v", s.Url, err)
		} else {
			log.Printf("Отписка от %s выполнена", s.Url)
		}
	}

	resp, err := api.Subscriptions.Subscribe(webhookURL, []string{})
	if err != nil {
		log.Fatal("Ошибка установки webhook:", err)
	}

	fmt.Printf("Webhook установлен: %s\n", webhookURL)
	fmt.Printf("Ответ: %#v\n", resp)
}
