package appcodec

import (
	"gateway/gw"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func proxy(message gw.Message, connector *gw.Connector) gw.Message {
	zuul := connector.Conf.Zuul
	sm := message.(*StringMessage)
	body := sm.Body
	serverName := body.SvrName
	resource := body.Resource
	action := strings.ToUpper(body.Action)
	content := body.Content
	atta := body.Attachment
	sid := body.SessionId
	var req *http.Request
	switch action {
	case "GET":
		req = newGetHttpRequest(zuul, serverName, resource, content, atta, sid)
	case "POST":
		req = newPostHttpRequest(zuul, serverName, resource, content, atta, sid)
	case "PUT":
		req = newPostHttpRequest(zuul, serverName, resource, content, atta, sid)
	case "DELETE":
		req = newPostHttpRequest(zuul, serverName, resource, content, atta, sid)
	case "OPTION":
	default:
		req = newGetHttpRequest(zuul, serverName, resource, content, atta, sid)
	}
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		gw.Logger.Errorf("Send %s Request Failure:%s,%s", action, serverName, resource)
		return newHttpErrorResponse(sm)
	}
	return processResponse(sm, resp)
}

func processResponse(src *StringMessage, resp *http.Response) *StringMessage {
	sm := NewEmptyMsg().(*StringMessage)
	header := sm.Header
	body := sm.Body
	header.EncryptType = src.Header.EncryptType
	header.MsgType = src.Header.MsgType
	header.ReqType = Response

	bs, err := ioutil.ReadAll(resp.Body)
	body.BType = src.Body.BType
	if err != nil {
		body.Attachment = "500"
	} else {
		body.Content = gw.Byte2str(bs)
		body.Attachment = strconv.FormatInt(int64(resp.StatusCode), 10)
	}
	return sm
}

func newHttpErrorResponse(srcStringMessage *StringMessage) *StringMessage {
	sm := NewEmptyMsg().(*StringMessage)
	header := sm.Header
	body := sm.Body
	header.EncryptType = srcStringMessage.Header.EncryptType
	header.MsgType = srcStringMessage.Header.MsgType
	header.ReqType = Response

	body.BType = srcStringMessage.Body.BType
	body.Attachment = "500"

	return sm
}

func newGetHttpRequest(zuul string, serverName string, resource string, content string, atta string, sid string) *http.Request {

	urlStr := zuul + "/" + serverName + resource
	if content != "" {
		urlStr = urlStr + "?" + content
	}
	gw.Logger.Infof("Get Request URL: %s", urlStr)
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		gw.Logger.Errorf("New GET Request error,url:%s,Error:%s", urlStr, err.Error())
		return nil
	}
	return req
}

func newPostHttpRequest(zuul string, serverName string, resource string, content string, atta string, sid string) *http.Request {
	urlStr := zuul + "/" + serverName + resource
	var req *http.Request
	var err error
	if content == "" {
		req, err = http.NewRequest("POST", urlStr, nil)
	} else {
		req, err = http.NewRequest("POST", urlStr, strings.NewReader(content))
	}

	if err != nil {
		gw.Logger.Errorf("New POST Request error,url:%s,Error:%s", urlStr, err.Error())
		return nil
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "sid")
	return req
}
