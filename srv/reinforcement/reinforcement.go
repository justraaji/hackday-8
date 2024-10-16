package main

import (
	"strings"

	"github.com/ollama/ollama/api"
)

type CompanyData struct {
	Username    string `json:"username"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	CompanyName string `json:"companyName"`
	Address     string `json:"address"`
	City        string `json:"city"`
	State       string `json:"state"`
	ZipCode     string `json:"zip_code"`
}

var stream = false

func ExtractJSONStringFromResponse(response api.ChatResponse) string {
	return response.Message.Content[strings.Index(response.Message.Content, "```json")+len("```json") : strings.LastIndex(response.Message.Content, "```")]
}

func IsSpam(response api.ChatResponse) bool {
	return strings.Contains(response.Message.Content, "\"isSpam\": true")
}

func IsNotSpam(response api.ChatResponse) bool {
	return strings.Contains(response.Message.Content, "\"isSpam\": false")
}
