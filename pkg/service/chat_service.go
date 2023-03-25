package service

import (
	"blog-go/pkg/model"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func GenerateGPTResponse(c *gin.Context,receivedMessage *model.TextMessage)  {
	apiURL := "https://api.openai.com/v1/chat/completions"
	messages := make([]model.Message,0)
	messages = append(messages,model.Message{
		Role:    "user",
		Content: receivedMessage.Content,
	})
	data := &model.OpenAIRequest{
		Model:    "gpt-3.5-turbo",
		Messages: messages,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", model.OpenAIAPIClient))

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var result model.OpenAIResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return
	}
	reply := strings.TrimSpace(result.Choices[0].Message.Content)
	response := model.TextMessage{
		ToUserName:   receivedMessage.ToUserName,
		FromUserName: receivedMessage.FromUserName,
		CreateTime:   time.Now().Unix(),
		MsgType:      receivedMessage.MsgType,
		Content:       reply,
	}
	c.XML(http.StatusOK,response)
	return
}