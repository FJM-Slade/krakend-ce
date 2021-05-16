package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"

	"github.com/devopsfaith/krakend/config"
)

type Message struct {
	Create_time  int64  `json:"create_time"`
	Event_level  int    `json:"event_level"`
	Owner        string `json:"owner"`
	Version      string `json:"version"`
	Type         string `json:"type"`
	Business_id  string `json:"business_id"`
	Uow_id       string `json:"uow_id"`
	Route_id     string `json:"route_id"`
	Step_id      string `json:"step_id"`
	Event_type   string `json:"event_type"`
	Tx_userid    string `json:"tx_userid"`
	Tx_ttl       int32  `json:"tx_ttl"`
	Msg_data     string `json:"msg_data"`
	User_msg     string `json:"user_msg"`
	Message_id   string `json:"message_id"`
	Request_url  string `json:"request_url"`
	Url_template string `json:"url_template"`
	Local_ip     string `json:"local_ip"`
	Remote_ip    string `json:"remote_ip"`
	Queryparam   string `json:"queryparam"`
	Http_method  string `json:"http_method"`
	Infos        string `json:"infos"`
	Confidential int8   `json:"confidential"`
	Appname      string `json:"appname"`
}

var activeMQcfg string

func IntegrationHub(cfg config.ServiceConfig, activeMQaddr string) gin.HandlerFunc {
	return func(c *gin.Context) {

		fmt.Println("\nall extra config -->  ", &cfg.ExtraConfig)
		fmt.Println("\n config part", cfg.ExtraConfig["pt/i2s/utl/integrationhub/gateway"])

		txId := getTxID(c.Request)

		fmt.Println("txid :", txId)

		activeMQcfg = activeMQaddr

		fillRequest(txId, c)

		fillResponse(txId, c)

	}

}

//fill string for version
func fillRequest(txId string, c *gin.Context) {

	var msg Message

	msg.Create_time = time.Now().UnixNano() / int64(time.Millisecond)
	msg.Event_level = 200
	msg.Owner = "DFT_OWNER"
	msg.Version = "DFT_VERSION"
	msg.Type = "DFT_TYPE"
	msg.Business_id = "DFT_BO"
	msg.Uow_id = txId
	msg.Route_id = c.Request.URL.Path
	msg.Step_id = "TXBEGIN"
	msg.Event_type = "4"
	msg.Tx_ttl = 0

	buf := make([]byte, 1024)
	num, _ := c.Request.Body.Read(buf)
	reqBody := string(buf[0:num])
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(reqBody))) // Write body back
	msg.Msg_data = string(buf[0:num])

	msg.User_msg = c.Request.URL.User.Username()
	//msg.Message_id =
	msg.Request_url = c.Request.URL.RawPath
	//msg.url_template = c.Request.
	msg.Local_ip = c.ClientIP()
	msg.Remote_ip = c.Request.RemoteAddr
	msg.Queryparam = c.Request.URL.RawQuery
	msg.Http_method = c.Request.Method
	//msg.infos
	msg.Confidential = 1
	msg.Appname = "XXXXXXXXX"
	out, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	fmt.Println("Request struct: " + string(out))

	c.Request.Header.Set("txId", txId)

	go func() {
		if err := NewActiveMQ(activeMQcfg).Send("/queue/proxyEventsQueue", string(out)); err != nil {
			fmt.Println("AMQ ERROR:", err)
		}

	}()
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func fillResponse(txId string, c *gin.Context) {
	fmt.Println("Request txId: " + c.GetHeader("txId") + "\n")

	var msgOut Message

	blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
	c.Writer = blw

	c.Next()
	c.Writer.Header().Set("txId", txId)

	if strings.Contains((c.Writer.Header().Get("Content-Type")), "xml") {

		body, _ := ioutil.ReadAll(blw.body)

		msgOut.Msg_data = string(body)
	} else {
		fmt.Println("cont2" + (c.Writer.Header().Get("Content-Type")))
		msgOut.Msg_data = blw.body.String()
	}

	msgOut.Create_time = time.Now().UnixNano() / int64(time.Millisecond)
	msgOut.Event_level = c.Writer.Status()
	msgOut.Owner = "DFT_OWNER"
	msgOut.Version = "DFT_VERSION"
	msgOut.Type = "DFT_TYPE"
	msgOut.Business_id = "DFT_BO"
	msgOut.Uow_id = txId
	msgOut.Route_id = c.Request.URL.Path
	msgOut.Step_id = "TXEND"
	msgOut.Event_type = "4"
	msgOut.Tx_ttl = 0

	msgOut.User_msg = c.Request.URL.User.Username()
	//msgOut.Message_id =
	msgOut.Request_url = c.Request.URL.RawPath
	//msgOut.url_template = c.Request.
	msgOut.Local_ip = c.Request.RemoteAddr
	msgOut.Remote_ip = c.Request.RemoteAddr
	msgOut.Queryparam = c.Request.URL.RawQuery
	msgOut.Http_method = c.Request.Method
	//msgOut.infos
	msgOut.Confidential = 1
	msgOut.Appname = "XXXXXXXXX"
	out, err := json.Marshal(msgOut)
	if err != nil {
		panic(err)
	}
	fmt.Println("Response struct: " + string(out) + "\n")

	go func() {
		if err := NewActiveMQ(activeMQcfg).Send("/queue/proxyEventsQueue", string(out)); err != nil {
			fmt.Println("AMQ ERROR:", err)
		}

	}()

}

func getTxID(r *http.Request) string {
	if len(r.Header.Get("txId")) > 0 {
		fmt.Println("Same txid. Will not generate a new one.")
		return r.Header.Get("txId")
	}
	return uuid.Must(uuid.NewV4(), nil).String()
}
