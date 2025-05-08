package ws

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/IBM/sarama"
	"github.com/gorilla/websocket"

	srvApi "githib.com/zamatay/otus/arch/lesson-1/internal/api"
	dialogs2 "githib.com/zamatay/otus/arch/lesson-1/internal/api/ws/dialogs"
	domainmes "githib.com/zamatay/otus/arch/lesson-1/internal/api/ws/domain"
	"githib.com/zamatay/otus/arch/lesson-1/internal/api/ws/posts"
	"githib.com/zamatay/otus/arch/lesson-1/internal/config"
	"githib.com/zamatay/otus/arch/lesson-1/internal/domain"
	"githib.com/zamatay/otus/arch/lesson-1/internal/kafka"
	"githib.com/zamatay/otus/arch/lesson-1/internal/repository"
)

type service interface {
}

type WS struct {
	srv               service
	activeConnections map[int]*websocket.Conn
	dialogs           *dialogs2.Dialogs
	posts             *posts.Posts
	consumer          *kafka.Consumer
	messageChanel     chan *sarama.ConsumerMessage
}

func (w WS) handleWebSocket(writer http.ResponseWriter, request *http.Request) {
	up := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := up.Upgrade(writer, request, nil)
	if err != nil {
		log.Println("Ошибка апгрейда:", err)
		return
	}
	log.Println("Клиент подключен")
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Ошибка чтения:", err)
			break
		}

		if err := w.Handle(conn, messageType, message); err != nil {
			log.Println("Ошибка записи:", err)
			break
		}
	}
	for key, connection := range w.activeConnections {
		if connection == conn {
			delete(w.activeConnections, key)
		}
	}
	log.Println("Клиент отключен")
}

func NewWS(ctx context.Context, service *repository.Repo, s *srvApi.Service, config *config.Config) (*WS, error) {
	ws := new(WS)
	ws.srv = service
	ws.dialogs = dialogs2.NewDialogs(service)
	ws.posts = posts.NewPosts(service)
	ws.activeConnections = make(map[int]*websocket.Conn)

	ws.messageChanel = kafka.GetChannel()
	var err error
	if ws.consumer, err = kafka.NewConsumer(&config.Kafka, ws.messageChanel); err != nil {
		return nil, err
	}

	s.GetRoute().HandleFunc("/ws", ws.handleWebSocket)

	go ws.consumer.Process(ctx)

	go ws.processConsumer()

	return ws, nil
}

func (w WS) processConsumer() {
	for msg := range w.messageChanel {
		for _, header := range msg.Headers {
			if string(header.Key) != "message-type" {
				continue
			}

			if string(header.Value) == "posts/create" {
				tmp := struct {
					Id        string    `json:"id"`
					UserId    int       `json:"user_id"`
					Text      string    `json:"text"`
					CreatedAt time.Time `json:"created_at"`
				}{}
				err := json.Unmarshal(msg.Value, &tmp)
				if err != nil {
					return
				}
				if w.activeConnections[tmp.UserId] != nil {
					post := domain.Post{
						ID:        tmp.Id,
						UserID:    tmp.UserId,
						Text:      tmp.Text,
						CreatedAt: tmp.CreatedAt,
					}
					value := domainmes.PostMessage{OkMessage: domainmes.OkMessage{Result: true, Action: posts.PostsPostAction}, Post: post}
					marshal, err := json.Marshal(value)
					if err != nil {
						return
					}
					if err := w.activeConnections[tmp.UserId].WriteMessage(websocket.TextMessage, marshal); err != nil {
						slog.Error("Ошибка при отправке сообщения пользователю", err)
						delete(w.activeConnections, tmp.UserId)
					}
				}
			}

			slog.Info("Поступило сообщение", msg)

		}
	}
}
