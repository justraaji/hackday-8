package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ollama/ollama/api"
)

func RunNotSpamReinforcement(ctx context.Context, client *api.Client) {
	notSpamMessages, err := GenerateNotSpamMessages()
	if err != nil {
		log.Fatal(err)
	}

	for _, message := range notSpamMessages {
		fmt.Println(message.Content)

		req := &api.ChatRequest{
			Model:  "testModel",
			Stream: &stream,
			Messages: []api.Message{
				message,
			},
		}

		respFunc := func(resp api.ChatResponse) error {
			fmt.Println(resp.Message.Content)
			output := ExtractJSONStringFromResponse(resp)
			fmt.Println(output)

			if IsSpam(resp) {
				CorrectNotSpam(ctx, client, message.Content)
			}

			return nil
		}

		err = client.Chat(ctx, req, respFunc)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func GenerateNotSpamMessages() ([]api.Message, error) {
	jsonData, err := LoadNotSpamFromCSV("./reinforcement/data/not_spam.csv")
	if err != nil {
		return nil, err
	}

	var messages []api.Message
	for _, jsonStr := range jsonData {
		message := api.Message{
			Role:    "user",
			Content: jsonStr,
		}
		messages = append(messages, message)
	}

	return messages, err
}

func LoadNotSpamFromCSV(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV: %v", err)
	}

	var jsonData []string
	for i, record := range records {
		if i == 0 {
			continue
		}

		data := CompanyData{
			Username:    strings.ToLower(strings.ReplaceAll(record[0][0:1]+record[1], " ", "")),
			FirstName:   record[0],
			LastName:    strings.ReplaceAll(record[1], "\\u0026", "&"),
			CompanyName: strings.ReplaceAll(record[2], "\\u0026", "&"),
			Address:     strings.ReplaceAll(record[3], "\\u0026", "&"),
			City:        strings.ReplaceAll(record[4], "\\u0026", "&"),
			State:       strings.ReplaceAll(record[6], "\\u0026", "&"),
			ZipCode:     strings.ReplaceAll(record[7], "\\u0026", "&"),
			PhoneNumber: strings.ReplaceAll(record[8], "\\u0026", "&"),
			Email:       strings.ReplaceAll(record[10], "\\u0026", "&"),
		}

		jsonBytes, err := json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("error converting to JSON: %v", err)
		}

		jsonData = append(jsonData, string(jsonBytes))
	}

	return jsonData, nil
}

func CorrectNotSpam(ctx context.Context, client *api.Client, input string) error {
	fmt.Println("Applying reinforcement for not spam")
	req := &api.ChatRequest{
		Model:  "testModel",
		Stream: &stream,
		Messages: []api.Message{
			{
				Role:    "user",
				Content: "The following was marked as spam, but is not spam: " + input,
			},
		},
	}

	respFunc := func(resp api.ChatResponse) error {
		fmt.Println(resp.Message.Content)

		return nil
	}

	err := client.Chat(ctx, req, respFunc)
	if err != nil {
		log.Fatal(err)
	}

	return err
}
