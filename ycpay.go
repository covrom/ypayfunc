package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
)

type RequestBody struct {
	HttpMethod string `json:"httpMethod"`
	Body       []byte `json:"body"`
}

type ResponseBody struct {
	StatusCode int         `json:"statusCode"`
	Body       interface{} `json:"body"`
}

// Точка входа для Yandex Cloud Functions
// Переменные окружения:
// MSRV - сервер и порт SSL майлера, например, smtp.yandex.ru:465
// MLGN - логин мыла
// MPSW - пароль мыла
// MLCC - копия уведомления на твое мыло
// PUSHPSW - yandex payment секрет
// После загрузки кода функции в Яндекс.Облако нужно сделать ее публичной и указать ссылку на нее здесь:
// https://yoomoney.ru/transfer/myservices/http-notification
func YaPay(ctx context.Context, request []byte) (*ResponseBody, error) {
	requestBody := &RequestBody{}
	if err := json.Unmarshal(request, &requestBody); err != nil {
		return nil, fmt.Errorf("an error has occurred when parsing request: %v", err)
	}
	YandexMoneyIncomingPush(string(requestBody.Body), os.Getenv("PUSHPSW"))
	return &ResponseBody{StatusCode: 200, Body: ""}, nil
}
