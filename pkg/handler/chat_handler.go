package handler

import (
	"blog-go/pkg/model"
	"blog-go/pkg/service"
	"crypto/sha1"
	"encoding/hex"
	"encoding/xml"
	"github.com/gin-gonic/gin"
	"github.com/go-xmlpath/xmlpath"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"time"
)

const (
	wechatToken     = "67_0sfguUnDaeQxKVC-ja-JHl64ses0ZF_6tyvf9Pd-kq0_FSF_O0Zo2xmNZY8C-WnzydWVjnDD4j9eOu_Ep6X0PKjFxdVa8hlKBrHgyTnio8hiKZq287k4B_tBx-cQJGdABAEXU"
)
func  VerifyData(c *gin.Context) {
	req := c.Request
	signature := req.URL.Query().Get("signature")
	timestamp := req.URL.Query().Get("timestamp")
	nonce := req.URL.Query().Get("nonce")
	echostr := req.URL.Query().Get("echostr")

	if checkSignature(wechatToken, signature, timestamp, nonce) {
		c.JSON(http.StatusOK,echostr)
		return
	} else {
		c.JSON(http.StatusOK,"invalid signature")
	}
}
func MessageHandler (c *gin.Context) {
	req := c.Request
	body, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()

	if err != nil {
		c.JSON(http.StatusBadRequest,err.Error())
		return
	}

	content, toUserName := parseXMLMessage(string(body))

	reply := service.GenerateGPTResponse(content)
	response := model.TextMessage{
		ToUserName:   toUserName,
		FromUserName: "gh_a835fe2e54c7",
		CreateTime:   time.Now().Unix(),
		MsgType:      "text",
		Content:       reply,
	}

	responseXML, err := xml.MarshalIndent(response, "", "  ")
	if err != nil {
		c.JSON(http.StatusOK,err.Error())
		return
	}
	c.XML(http.StatusOK,responseXML)
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