package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

var (
	address            = ":8080"
	listOfClients      = make(map[*websocket.Conn]bool)
	stopWordFromClient = "goodbye"
)

func serveHome(writer http.ResponseWriter, request *http.Request) {
	http.ServeFile(writer, request, "index.html")
}

func sendMessagesToClients() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		messageForClients := []byte("hello. i am here")

		listOfClientsToBeRemoved := make([]*websocket.Conn, 0) //список для передачи клиентов, чтобы в последующем их удалить

		for connectionForClient := range listOfClients {
			err := connectionForClient.WriteMessage(websocket.TextMessage, messageForClients)
			if err != nil {
				log.Print(errors.Wrapf(err, "can not send message"))

				listOfClientsToBeRemoved = append(listOfClientsToBeRemoved, connectionForClient)
			}
		}

		for _, connectionForRemoveClient := range listOfClientsToBeRemoved {
			delete(listOfClients, connectionForRemoveClient) //удаляем клиента
		}
	}
}

func handleMessageFromClient(connection *websocket.Conn, messageInByte []byte) bool {
	messageInString := string(messageInByte)

	if messageInString == stopWordFromClient {
		log.Print("client sent stop-word", stopWordFromClient, ". close connection")

		return false // чтобы указать, что нужно закрыть соединение
	}

	responseMessage := []byte("thank you for message")
	err := connection.WriteMessage(websocket.TextMessage, responseMessage)
	if err != nil {
		log.Print(errors.Wrapf(err, "can not send response message"))
	}

	return true // продолжаем обработку сообщений
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

	for {
		typeOfMessage, messageInBytes, err := connection.ReadMessage()
		if err != nil {
			log.Print(errors.Wrapf(err, "can not read message from client:", typeOfMessage))

			break
		}

		if !handleMessageFromClient(connection, messageInBytes) {
			break // при false выходим из цикла
		}
	}

	// удаление клиента после закрытия соединения
	delete(listOfClients, connection)
}

func main() {
	// создаем обработчик для двух путей. '/' будет обрабатываться функцией serveHome, которая возвращает файл index.html
	http.HandleFunc("/", serveHome)

	// '/ws' будет обрабатываться функцией connectWithWebSocket, которая устанавливает веб-сокет соединение
	http.HandleFunc("/ws", connectWithWebSocket)

	go sendMessagesToClients()

	log.Println("server listen and serve at:", address)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Print(errors.Wrapf(err, "can not connect with server"))
	}
}
