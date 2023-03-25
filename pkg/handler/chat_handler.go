package handler

import (
	"blog-go/pkg/model"
	"blog-go/pkg/service"
	"crypto/sha1"
	"encoding/hex"
	"encoding/xml"
	"github.com/gin-gonic/gin"
	"github.com/go-xmlpath/xmlpath"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
)

const (
	wechatToken     = "67_cLACSYnbUDKi0lCk_PKk8CidfO5RonJeLWWStL7cI5-_VE4e-bATO-oinkJSWUM4b6W1rQWumCIyfJAfH8s8UVcGcYzvjqNBNh_DY5qRxhfrvpFzf6qUgfAXZR4CDWcAHALM"
)
func  VerifyData(c *gin.Context) {
	req := c.Request
	signature := req.URL.Query().Get("signature")
	timestamp := req.URL.Query().Get("timestamp")
	nonce := req.URL.Query().Get("nonce")
	echostr := req.URL.Query().Get("echostr")
	log.Println(req)
	log.Println(signature)
	log.Println(nonce)
	log.Println(timestamp)
	log.Println(echostr)
	c.String(http.StatusOK,echostr)
	return
	//if checkSignature(wechatToken, signature, timestamp, nonce) {
	//	c.JSON(http.StatusOK,echostr)
	//	return
	//} else {
	//	c.JSON(http.StatusOK,"invalid signature")
	//}
}
func MessageHandler (c *gin.Context) {
	timeout := time.After(4000 * time.Millisecond)
	var receivedMessage model.TextMessage
	err := c.ShouldBindXML(&receivedMessage)
	if err != nil {
		log.Println(err.Error())
		c.String(http.StatusBadRequest, "Invalid XML")
		return
	}
	select {
	case <- timeout:
		response := model.TextMessage{
			ToUserName:   receivedMessage.FromUserName,
			FromUserName: receivedMessage.ToUserName,
			CreateTime:   time.Now().Unix(),
			MsgType:      receivedMessage.MsgType,
			Content:      "问题似乎太复杂了，我超时了~",
		}
		msg, err := xml.Marshal(&response)
		if err != nil {
			return
		}
		_, _ = c.Writer.Write(msg)
		return
	default:
		service.GenerateGPTResponse(c,&receivedMessage)
	}
}
func checkSignature(token, signature, timestamp, nonce string) bool {
	values := []string{token, timestamp, nonce}
	sort.Strings(values)

	hash := sha1.New()
	hash.Write([]byte(strings.Join(values, "")))
	generatedSignature := hex.EncodeToString(hash.Sum(nil))

	return signature == generatedSignature
}
func parseXMLMessage(xmlData string) (string, string) {
	root, err := xmlpath.Parse(strings.NewReader(xmlData))
	if err != nil {
		return "", ""
	}

	contentPath := xmlpath.MustCompile("//xml/Content")
	content, _ := contentPath.String(root)

	toUserNamePath := xmlpath.MustCompile("//xml/FromUserName")
	toUserName, _ := toUserNamePath.String(root)

	return content, toUserName
}