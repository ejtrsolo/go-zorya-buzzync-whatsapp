package tools

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

type DataGlobal struct {
	Data          map[string]interface{} `json:"data"`
	Message       string                 `json:"message"`
	Success       bool                   `json:"success"`
	Authorization bool                   `json:"-"`
}

// ExecuteConsultWithApikey Function that allows you to execute a query to a url with the apikey already configured and return the response data
func ExecuteConsultWithApikey(queryUrl, apikey, methodType string, sendBodyMap map[string]interface{}, isJsons bool) (map[string]interface{}, error) {
	sendBodyByte, err := json.Marshal(sendBodyMap)
	if err != nil {
		return nil, err
	}
	body := bytes.NewBuffer(sendBodyByte)

	headers := map[string]string{
		"x-api-key": apikey,
	}
	response := ConsultClient(queryUrl, methodType, headers, body, isJsons)

	bodyStr, ok := response.Data["body"].(string)
	if !ok {
		return nil, errors.New("ERROR: Body is not a string")
	}
	bodyBytes := []byte(bodyStr)

	var bodyMap map[string]interface{}
	err = json.Unmarshal(bodyBytes, &bodyMap)
	if err != nil {
		return nil, err
	}

	if !response.Success {
		return nil, errors.New(bodyMap["message"].(string))
	}

	return bodyMap["data"].(map[string]interface{}), nil
}

// ConsultClient Function that allows you to execute a query to a url with the apikey already configured and return the response data
func ConsultClient(url, method string, headers map[string]string, body io.Reader, isJsons ...bool) DataGlobal {
	isJson := false
	if len(isJsons) > 0 {
		isJson = isJsons[0]
	}
	log.Printf("ConsultClient Url: %v", url)
	data := DataGlobal{Data: make(map[string]interface{}), Authorization: true}
	if method == "GET" {
		body = nil
	}

	var req *http.Request
	var resp *http.Response
	var bodySend io.Reader
	var err error
	client := &http.Client{}
	if isJson {
		buf := new(bytes.Buffer)
		buf.ReadFrom(body)
		bodySend = bytes.NewBuffer(buf.Bytes())
	} else {
		bodySend = body
	}

	// Create the HTTP request with the prepared method, URL and body
	req, err = http.NewRequest(method, url, bodySend)
	if err != nil {
		fmt.Println("Error al crear la solicitud:", err)
		data.Message = fmt.Sprintf("ERROR: in url %s", url)
		return data
	}

	// Set request headers based on whether it is JSON or not
	if isJson {
		req.Header.Set("Content-Type", "application/json")
		for i, h := range headers {
			req.Header.Set(i, h)
		}
	} else {
		for i, h := range headers {
			req.Header.Add(i, h)
		}
	}

	// Submit the request
	resp, err = client.Do(req)
	if err != nil {
		log.Println("ERROR: Consult Client", err)
		data.Message = "ERROR: Consult Client"
		return data
	}

	defer resp.Body.Close()
	log.Println("ConsultClient Status Response: ", resp.Status)

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	dataResponseStr := buf.String()

	data.Success = true
	data.Data["body"] = dataResponseStr
	data.Data["status"] = resp.Status
	data.Data["status_code"] = resp.StatusCode

	return data
}
