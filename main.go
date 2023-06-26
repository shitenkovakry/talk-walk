package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

var (
	address       = ":8080"
	listOfClients = make(map[*websocket.Conn]bool)
	//channelForSendingMessage = make(chan []byte)
)

func serveHome(writer http.ResponseWriter, request *http.Request) {
	http.ServeFile(writer, request, "index.html")
}

func sendMessagesToClients(connection *websocket.Conn) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		messageForClients := []byte("hello. i am here")

		err := connection.WriteMessage(websocket.TextMessage, messageForClients)
		if err != nil {
			log.Print(errors.Wrapf(err, "can not write message"))
			break
		}
	}
}

func connectWithWebSocket(writer http.ResponseWriter, request *http.Request) {
	// настройка веб-сокет соединения
	upgrader := websocket.Upgrader{} // обработка http-соединения до веб-сокета

	connection, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		log.Print(errors.Wrapf(err, "failed to set up connection with websocket"))

		return
	}
	defer connection.Close()

	// добавление клиента в список
	listOfClients[connection] = true

	go sendMessagesToClients(connection)

	// ожидание закрытия соединения
	<-make(chan struct{})

	// удаление клиента после закрытия соединения
	delete(listOfClients, connection)
}

func main() {
	// создаем обработчик для двух путей. '/' будет обрабатываться функцией serveHome, которая возвращает файл index.html
	http.HandleFunc("/", serveHome)

	// '/ws' будет обрабатываться функцией connectWithWebSocket, которая устанавливает веб-сокет соединение
	http.HandleFunc("/ws", connectWithWebSocket)

	log.Println("server listen and serve at:")
	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Print(errors.Wrapf(err, "can not connect with server"))
	}
}
