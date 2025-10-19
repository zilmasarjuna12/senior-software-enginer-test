package model

import "time"

type Response struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp string      `json:"timestamp"`
}

func NewResponseSuccess(data interface{}) *Response {
	return &Response{
		Success:   true,
		Message:   "success",
		Data:      data,
		Timestamp: time.Now().Format(time.RFC3339),
	}
}

func NewResponseError(message string) *Response {
	return &Response{
		Success:   false,
		Message:   message,
		Timestamp: time.Now().Format(time.RFC3339),
	}
}
