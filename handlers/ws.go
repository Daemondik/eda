package handlers

import (
	"eda/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/streadway/amqp"
	"log"
	"net/http"
	"sync"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	conn *websocket.Conn
	id   string
}

var clients = make(map[string]*Client)
var mutex = &sync.Mutex{}

func Chat(c *gin.Context) {
	id := c.Param("user_id")
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Print("Ошибка при подключении WebSocket:", err)
		return
	}

	mutex.Lock()
	if len(clients) >= 2 {
		conn.Close()
		mutex.Unlock()
		log.Println("Превышено количество пользователей")
		return
	}

	client := &Client{conn: conn, id: id}
	clients[id] = client
	mutex.Unlock()

	channel, err := models.RabbitMQQueue(models.ChatQueue)
	if err != nil {
		log.Println(err)
		return
	}
	//defer channel.Close()

	// Чтение сообщений от пользователя
	go func() {
		defer func() {
			conn.Close()
			mutex.Lock()
			delete(clients, id)
			mutex.Unlock()
			channel.Close()
		}()
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("Ошибка чтения:", err)
				break
			}

			// Отправка сообщения в RabbitMQ
			err = channel.Publish(
				"",               // exchange
				models.ChatQueue, // routing key
				false,            // mandatory
				false,            // immediate
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        msg,
				})
			if err != nil {
				log.Println("Ошибка отправки сообщения:", err)
				break
			}
		}
		mutex.Lock()
		delete(clients, id)
		mutex.Unlock()
	}()

	// Потребление сообщений из RabbitMQ и отправка пользователям
	go func() {
		defer channel.Close()
		msgs, err := channel.Consume(
			models.ChatQueue, // queue
			"",               // consumer
			true,             // auto-ack
			false,            // exclusive
			false,            // no-local
			false,            // no-wait
			nil,              // args
		)
		if err != nil {
			log.Println(err)
		}

		for d := range msgs {
			mutex.Lock()
			for _, client := range clients {
				err := client.conn.WriteMessage(websocket.TextMessage, d.Body)
				if err != nil {
					log.Println("Ошибка отправки сообщения клиенту:", err)
					client.conn.Close()
					delete(clients, client.id)
				}
			}
			mutex.Unlock()
		}
	}()
}
