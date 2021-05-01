package middleware

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

const (
	headerTenant = "X-Tenant"
)

type message struct {
	create_time  int32
	event_level  int32
	owner        string
	version      string
	type2        string
	business_id  string
	uow_id       string
	route_id     string
	step_id      string
	event_type   string
	tx_userid    string
	tx_ttl       int32
	msg_data     string
	user_msg     string
	message_id   string
	request_url  string
	url_template string
	local_ip     string
	remote_ip    string
	queryparam   string
	http_method  string
	infos        string
	confidential int8
	appname      string
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
		if err := NewActiveMQ("localhost:61623").Send("/queue/proxyEventsQueue", "{\"create_time\":1619028534209,\"event_level\":400,\"owner\":\"DFT_OWNER\",\"version\":\"DFT_VERSION\",\"type\":\"DFT_TYPE\",\"business_id\":\"DFT_BO\",\"uow_id\":\"198582f8-6414-46c7-bccb-c67910a554d9\",\"route_id\":\"/eGISGIVProxy/services/ProxyConsultaGIV\",\"step_id\":\"TXBEGIN\",\"event_type\":4,\"tx_userid\":\"some_user\",\"tx_ttl\":0,\"msg_data\":\"<soapenv:Envelope xmlns:soapenv=\\\"http://schemas.xmlsoap.org/soap/envelope/\\\" xmlns:prox=\\\"http://proxy.web.egis.i2s.com\\\"><soapenv:Header/><soapenv:Body> <prox:getValorMaxResgatavelApolice> <xml_in><![CDATA[<?xml version=\\\"1.0\\\" encoding=\\\"UTF-8\\\" standalone=\\\"yes\\\"?><Consulta><Apolice><Sistema>GIVITG</Sistema><Modalidade>5111</Modalidade><NumeroApolice>111185</NumeroApolice><TipoSinistro>1</TipoSinistro><Data>20181204</Data></Apolice></Consulta>]]></xml_in></prox:getValorMaxResgatavelApolice></soapenv:Body></soapenv:Envelope>\",\"user_msg\":null,\"message_id\":null,\"request_url\":null,\"url_template\":null,\"local_ip\":null,\"remote_ip\":\"10.0.75.1\",\"queryparam\":\"\",\"http_method\":\"POST\",\"infos\":null,\"confidential\":1,\"appname\":\"PROXYAPP\"}"); err != nil {
			fmt.Println("AMQ ERROR:", err)
		}
	*/
	c.Next()

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

}



//get queue name  and queue address from file configuration

//fill string for version
func fillPayload(txId string, c *gin.Context) {

	fmt.Println("Request header X-tenant: " + c.GetHeader("X-tenant"))
	fmt.Println("Request method: " + c.Request.Method)
	fmt.Println("Request proto minor: " + strconv.FormatInt((int64(c.Request.ProtoMinor)), 10))
	rsBodyBytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return
	}
	fmt.Println("Request rsBodyBytes: " + string(rsBodyBytes))

	c.Request.Header.Set("txId", txId)
	fmt.Println("Request txId: " + c.GetHeader("txId") + "\n")

}


func getTxID(r *http.Request) string {
	if len(r.Header.Get("txId")) > 0 {
		return r.Header.Get("txId")
	}
	return uuid.Must(uuid.NewV1(), nil).String()
}

