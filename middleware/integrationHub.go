package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Message struct {
	Create_time  int32
	Event_level  int
	Owner        string
	Version      string
	Type         string
	Business_id  string
	Uow_id       string
	Route_id     string
	Step_id      string
	Event_type   string
	Tx_userid    string
	Tx_ttl       int32
	Msg_data     string
	User_msg     string
	Message_id   string
	Request_url  string
	Url_template string
	Local_ip     string
	Remote_ip    string
	Queryparam   string
	Http_method  string
	Infos        string
	Confidential int8
	Appname      string
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

var i int

// HeaderLogs appends logging related headers to request
func IntegrationHub(c *gin.Context) {

	fmt.Println("\n>>>>vez :", i)

	txId := getTxID(c.Request)

	fillPayload(txId, c)

	//c.Next()

}

//get queue name  and queue address from file configuration

//fill string for version
func fillPayload(txId string, c *gin.Context) {

	var msg Message

	msg.Create_time = timestamppb.Now().GetNanos()
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
	/*msgBodyBytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		//return
	}
	msg.Msg_data = string(msgBodyBytes) */
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

	fmt.Println("\nRequest method: " + c.Request.Method)
	fmt.Println("Request proto minor: " + strconv.FormatInt((int64(c.Request.ProtoMinor)), 10))
	//fmt.Println("Request rsBodyBytes: " + string(msgBodyBytes))

	c.Request.Header.Set("txId", txId)

	fmt.Println("Request txId: " + c.GetHeader("txId") + "\n")

	go func() {
		if err := NewActiveMQ("localhost:61623").Send("/queue/proxyEventsQueue", string(out)); err != nil {
			fmt.Println("AMQ ERROR:", err)
		}

	}()

	ginBodyLogMiddleware(txId, c)

	//time.Sleep(5 * time.Second)

}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func ginBodyLogMiddleware(txId string, c *gin.Context) {
	blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
	c.Writer = blw
	c.Next()
	statusCode := c.Writer.Status()
	c.Writer.Header().Set("txId", txId)
	fmt.Println("Response status code: " + strconv.FormatInt(int64(statusCode), 10))
	fmt.Println("Response txId " + c.Writer.Header().Get("txId"))

	fmt.Println("Response body: " + blw.body.String())

	var msgOut Message

	msgOut.Create_time = timestamppb.Now().GetNanos()
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
	msgOut.Msg_data = blw.body.String()
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
		if err := NewActiveMQ("localhost:61623").Send("/queue/proxyEventsQueue", string(out)); err != nil {
			fmt.Println("AMQ ERROR:", err)
		}

	}()

}

func getTxID(r *http.Request) string {
	if len(r.Header.Get("txId")) > 0 {
		return r.Header.Get("txId")
	}
	return uuid.Must(uuid.NewV1(), nil).String()
}
