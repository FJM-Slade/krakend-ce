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

// HeaderLogs appends logging related headers to request
func IntegrationHub(c *gin.Context) {

	txId := getTxID(c.Request)

	fillPayload(txId, c)
	ginBodyLogMiddleware(txId, c)

	/*
		if err := NewActiveMQ("localhost:61623").Send("/queue/proxyEventsQueue", "{\"create_time\":1619028534209,\"event_level\":400,\"owner\":\"DFT_OWNER\",\"version\":\"DFT_VERSION\",\"type\":\"DFT_TYPE\",\"business_id\":\"DFT_BO\",\"uow_id\":\"198582f8-6414-46c7-bccb-c67910a554d9\"
		,\"route_id\":\"/eGISGIVProxy/services/ProxyConsultaGIV\",\"step_id\":\"TXBEGIN\",\"event_type\":4,\"tx_userid\":\"some_user\",
		\"tx_ttl\":0,
		\"msg_data\":\"<soapenv:Envelope xmlns:soapenv=\\\"http://schemas.xmlsoap.org/soap/envelope/\\\" xmlns:prox=\\\"http://proxy.web.egis.i2s.com\\\"><soapenv:Header/><soapenv:Body> <prox:getValorMaxResgatavelApolice> <xml_in><![CDATA[<?xml version=\\\"1.0\\\" encoding=\\\"UTF-8\\\" standalone=\\\"yes\\\"?><Consulta><Apolice><Sistema>GIVITG</Sistema><Modalidade>5111</Modalidade><NumeroApolice>111185</NumeroApolice><TipoSinistro>1</TipoSinistro><Data>20181204</Data></Apolice></Consulta>]]></xml_in></prox:getValorMaxResgatavelApolice></soapenv:Body></soapenv:Envelope>\"
		,\"user_msg\":null,\"message_id\":null,\"request_url\":null,\"url_template\":null,
		\"local_ip\":null,\"remote_ip\":\"10.0.75.1\",\"queryparam\":\"\",\"http_method\":\"POST\",\"infos\":null,\"confidential\":1,\"appname\":\"PROXYAPP\"}"); err != nil {
			fmt.Println("AMQ ERROR:", err)
		}
	*/
	c.Next()

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

	//if err := NewActiveMQ("localhost:61623").Send("/queue/proxyEventsQueue", string(out)); err != nil {
	//	fmt.Println("AMQ ERROR:", err)
	//}

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

	//if err := NewActiveMQ("localhost:61623").Send("/queue/proxyEventsQueue", string(out)); err != nil {
	//	fmt.Println("AMQ ERROR:", err)
	//}

}

func getTxID(r *http.Request) string {
	if len(r.Header.Get("txId")) > 0 {
		return r.Header.Get("txId")
	}
	return uuid.Must(uuid.NewV1(), nil).String()
}
