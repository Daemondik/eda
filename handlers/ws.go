package handlers

import (
	"eda/models"
	"github.com/gorilla/websocket"
	"github.com/streadway/amqp"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Chat(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	ch, err := models.NewRabbitMQQueue(models.ChatQueue)
	if err != nil {
		log.Println(err)
		return
	}
	defer ch.Close()

	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}
			err = ch.Publish(
				"",
				"chat",
				false,
				false,
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        message,
				})
			if err != nil {
				log.Println("publish:", err)
				break
			}
		}
	}()

	for {
		msgs, err := ch.Consume(
			models.ChatQueue,
			"",
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			log.Println("consume:", err)
			break
		}

		for d := range msgs {
			err = conn.WriteMessage(websocket.TextMessage, d.Body)
			if err != nil {
				log.Println("write:", err)
				break
			}
		}
	}
}
