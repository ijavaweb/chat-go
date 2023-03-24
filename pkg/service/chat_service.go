package service

import (
	"blog-go/pkg/model"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func GenerateGPTResponse(prompt string) string {
	fmt.Print(prompt)
	apiURL := "https://api.openai.com/v1/chat/completions"
	messages := make([]model.Message,0)
	messages = append(messages,model.Message{
		Role:    "user",
		Content: prompt,
	})
	data := &model.OpenAIRequest{
		Model:    "gpt-3.5-turbo",
		Messages: messages,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "Error generating response"
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "Error generating response"
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", model.OpenAIAPIClient))

	resp, err := client.Do(req)
	if err != nil {
		return "Error generating response"
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "Error generating response"
	}

	fmt.Println(string(body))
	var result model.OpenAIResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "Error generating response"
	}
	return strings.TrimSpace(result.Choices[0].Message.Content)
}