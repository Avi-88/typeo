package config

import (
	"encoding/json"
	"log"
	"os"
	"time"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/pborman/uuid"
)

type ClientMessage struct {
	Type      string     `json:"type"`
	SessionId string `json:"session_id,omitempty"`
	Data      string     `json:"data"`
}

type ServerMessage struct {
	Type      string    `json:"type"`
	SessionId string `json:"session_id"`
	Data      string    `json:"data"`
	WPM       int       `json:"wpm"`
	Accuracy  float64   `json:"accuracy"`
}

type UserSession struct {
	SessionId    string `json:"session_id"`
	TypedText    string    `json:"typed_text"`
	CorrectChars int       `json:"correct_chars"`
	Prompt       string    `json:"prompt"`
	WPM          int       `json:"wpm"`
	Accuracy     float64   `json:"accuracy"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

var userSessions map[string]UserSession = make(map[string]UserSession)

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

func LoadSocketServer(c *gin.Context) {
	upgrader := websocket.Upgrader{}
	w, r := c.Writer, c.Request
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	defer socket.Close()

	HandleMessages(socket)
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
			if len(messageText.SessionId) == 0 {
				log.Println("Missing sessionId in progress message")
				return
			}
			HandleProgress(socket, messageText.SessionId, messageText.Data)
		}
	}
}

func HandleInit(socket *websocket.Conn) {
	sessionId := uuid.New()
	message := ServerMessage{
		Type:      "Init",
		SessionId: sessionId,
		Data:      demoText,
		WPM:       0,
		Accuracy:  0,
	}
	userSession := UserSession{
		SessionId: sessionId,
		TypedText: "",
		Prompt:    demoText,
		WPM:       0,
		Accuracy:  0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	userSessions[sessionId] = userSession

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

func HandleProgress(socket *websocket.Conn, sessionId string, data string) {
	currentSession, ok := userSessions[sessionId]
	if !ok {
		log.Fatal("Session not found")
		return
	}
	prevLength := len(currentSession.TypedText)
	currentSession.TypedText = data
	currentSession.UpdatedAt = time.Now()
	charsTyped := len(currentSession.TypedText)
	elapsedTime := time.Since(currentSession.CreatedAt).Minutes()
	if elapsedTime > 0 { // Avoid division by zero
		currentSession.WPM = int((float64(charsTyped) / 5) / elapsedTime)
	}
	if prevLength < charsTyped { // User typed a new character
		lastIndex := charsTyped - 1
		if lastIndex < len(currentSession.Prompt) && data[lastIndex] == currentSession.Prompt[lastIndex] {
			currentSession.CorrectChars++
		}
	} else if prevLength > charsTyped { // User deleted a character
		// If deleted character was correct, decrement count
		if prevLength-1 < len(currentSession.Prompt) && currentSession.TypedText[prevLength-1] == currentSession.Prompt[prevLength-1] {
			currentSession.CorrectChars--
		}
	}
	if charsTyped > 0 {
		currentSession.Accuracy = (float64(currentSession.CorrectChars) / float64(charsTyped)) * 100
	} else {
		currentSession.Accuracy = 0
	}
	userSessions[sessionId] = currentSession
	message := ServerMessage{
		Type:      "Progress",
		SessionId: sessionId,
		Data:      currentSession.TypedText,
		WPM:       currentSession.WPM,
		Accuracy:  currentSession.Accuracy,
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
