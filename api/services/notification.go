package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SendData struct {
	ID          primitive.ObjectID `bson:"id" json:"id"`
	Description string             `bson:"description" json:"description"`
}

func NotificationSend(id primitive.ObjectID, str string, session string) error {
	data := SendData{
		ID:          id,
		Description: str,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Ошибка при маршалинге JSON:", err)
		return err
	}

	req, err := http.NewRequest("POST", "http://192.168.0.18:8084/api/telegram/notification", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Ошибка при создании запроса:", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	if session != "" {
		cookie := &http.Cookie{
			Name:  "session",
			Value: session,
		}
		req.AddCookie(cookie)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка при выполнении запроса:", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		fmt.Println("Ошибка ответа сервера:", resp.Status)
		return fmt.Errorf("server returned status: %d", resp.StatusCode)
	}

	return nil
}
