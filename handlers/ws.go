package handlers

import (
	"eda/logger"
	"eda/models"
	"eda/utils/security"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"log"
	"net/http"
	"strconv"
	"sync"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return r.Host == "localhost:8080"
	},
}

type Client struct {
	conn *websocket.Conn
	id   string
}

var clients = make(map[string]*Client)
var mutex = &sync.Mutex{}

func Chat(c *gin.Context) {
	id := c.Param("recipient_id")
	recipientId, err := strconv.Atoi(id)
	if err != nil {
		log.Println(err)
		return
	}

	senderId, err := security.GetUserIdByJWTOrOauth(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизирован"})
		return
	}

	sender, err := models.GetUserByID(senderId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	recipient, err := models.GetUserByID(uint(recipientId))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при подключении WebSocket"})
		return
	}

	mutex.Lock()
	if len(clients) >= 2 {
		conn.Close()
		log.Println("Превышено количество пользователей")
		return
	}
	mutex.Unlock()

	mutex.Lock()
	defer mutex.Unlock()

	if _, exists := clients[id]; exists {
		conn.Close()
		log.Println("Клиент уже существует")
		return
	}

	client := &Client{conn: conn, id: strconv.Itoa(int(senderId))}
	clients[id] = client

	channel, err := models.RabbitMQQueue(models.ChatQueue)
	if err != nil {
		log.Println(err)
		return
	}

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

			message := models.Message{
				Text:      string(msg),
				Sender:    sender,
				Recipient: recipient,
			}

			encodedMessage, err := encodeMessage(message)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при подключении WebSocket"})
				return
			}

			err = channel.Publish(
				"",
				models.ChatQueue,
				false,
				false,
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        encodedMessage,
				})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка отправки сообщения"})
				break
			}
		}
		mutex.Lock()
		delete(clients, id)
		mutex.Unlock()
	}()

	go func() {
		defer channel.Close()
		msgs, err := channel.Consume(
			models.ChatQueue,
			"",
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
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

func encodeMessage(message models.Message) ([]byte, error) {
	encoded, err := json.Marshal(message)
	if err != nil {
		logger.Log.Error("Ошибка при кодировании сообщения: ", zap.Error(err))
		return nil, err
	}
	return encoded, nil
}
