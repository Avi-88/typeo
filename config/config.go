package config

import (
	"log"
	"os"
	"github.com/google/uuid"
    "encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"time"
)


type ClientMessage struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

type ServerMessage struct {
	Type string `json:"type"`
	Data string `json:"data"`
	WPM int `json:"wpm"`
	Accuracy float64 `json:"accuracy"`
}

type UserSession struct {
	SessionId uuid.UUID `json:"session_id"`
	TypedText string `json:"typed_text"`
	Prompt string `json:"prompt"`
	WPM int `json:"wpm"`
	Accuracy float64 `json:"accuracy"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var

// LoadEnv loads environment variables from .env file
func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

// GetEnv gets the value of an environment variable
func GetEnv(key string) string {
	return os.Getenv(key)
}

var demoText string = "An application can also send and receive messages using the io.WriteCloser and io.Reader interfaces. To send a message, call the connection NextWriter method to get an io.WriteCloser, write the message to the writer and close the writer when done. To receive a message, call the connection NextReader method to get an io.Reader and read until io.EOF is returned. This snippet shows how to echo messages using the NextWriter and NextReader methods:"


func LoadSocketServer (c *gin.Context) {
	upgrader := websocket.Upgrader{}
	w, r := c.Writer, c.Request
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	defer socket.Close()
	
	HandleMessages(socket);
}

func HandleMessages(socket *websocket.Conn) {
	for {
		_, message, err := socket.ReadMessage()
		if err != nil {
			log.Fatal(err)
			return
		}
		var messageText ClientMessage
		if err := json.Unmarshal(message, &messageText); err != nil {
			log.Fatal("There was an error unmarshalling the message: ", err)
			return
		}
		switch messageText.Type {
		case "Init":
			HandleInit(socket)
		case "Progress":
			HandleProgress(socket, messageText.Data)
		}
	}
}

func HandleInit(socket *websocket.Conn) {
	message := ServerMessage{
		Type: "Init",
		Data: demoText,
		WPM: 0,
		Accuracy: 0,
	}
	sessionId := uuid.New()
	userSession := UserSession{
		SessionId: sessionId,
		TypedText: "",
		Prompt: demoText,
		WPM: 0,
		Accuracy: 0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Fatal("There was an error marshalling the message: ", err)
		return
	}
	if err := socket.WriteMessage(websocket.TextMessage, []byte(messageBytes)); err != nil {
		log.Fatal(err)
		return
	}
}

func HandleProgress (socket *websocket.Conn, data string) {

}