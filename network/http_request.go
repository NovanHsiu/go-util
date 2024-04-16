package network

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type HttpClient struct {
	TimeoutSeconds             int
	InsecureSkipVerify         bool
	EnabledSingledResuedClient bool
	client                     *http.Client
}

func NewHttpClient(timeoutSeconds int, insecureSkipVerify, enabledSingledResuedClient bool) *HttpClient {
	httpClient := HttpClient{
		TimeoutSeconds:             timeoutSeconds,
		InsecureSkipVerify:         insecureSkipVerify,
		EnabledSingledResuedClient: enabledSingledResuedClient,
	}
	httpClient.client = &http.Client{
		Timeout: time.Duration(time.Duration(timeoutSeconds) * time.Second),
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: insecureSkipVerify},
		},
	}
	return &httpClient
}

func NewDefaultHttpClient() *HttpClient {
	return NewHttpClient(30, true, true)
}

func (c *HttpClient) getClient() *http.Client {
	if c.EnabledSingledResuedClient {
		return c.client
	}
	client := &http.Client{
		Timeout: time.Duration(time.Duration(c.TimeoutSeconds) * time.Second),
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: c.InsecureSkipVerify},
		},
	}
	return client
}

func (c *HttpClient) closeIdleConnections(client *http.Client) {
	if !c.EnabledSingledResuedClient {
		client.CloseIdleConnections()
	}
}

func (c *HttpClient) Close() {
	c.client.CloseIdleConnections()
}

func extractBody(res *http.Response) (int, map[string]interface{}, error) {
	body, _ := ioutil.ReadAll(res.Body)
	if len(body) >= 3 {
		// remove bom
		if body[0] == 239 && body[1] == 187 && body[2] == 191 {
			body = body[3:]
		}
	}
	jsonMap := make(map[string]interface{})
	err := json.Unmarshal(body, &jsonMap)
	if err != nil {
		jsonMap["data"] = body
		return res.StatusCode, jsonMap, err
	}
	return res.StatusCode, jsonMap, nil
}

func setHeader(req *http.Request, header map[string]string) {
	for key := range header {
		req.Header.Set(key, header[key])
	}
}

func (c *HttpClient) SendBodyRequest(method, url, jsonStr string, header map[string]string) (int, map[string]interface{}, error) {
	pt := time.Now()
	jsonBytes := []byte(jsonStr)
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		log.Println("SendBodyRequest http client new request error:", err)
		return 500, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connection", "application/json")
	setHeader(req, header)
	client := c.getClient()
	res, err := client.Do(req)
	defer c.closeIdleConnections(client)
	if err != nil {
		log.Println("SendBodyRequest http client do error:", err)
		return 500, nil, err
	}
	defer res.Body.Close()
	showRespTimeLog(url, pt)
	return extractBody(res)
}

type SendFile struct {
	ParamName string
	Paths     []string
}

func (c *HttpClient) SendFormDataWithFilesRequest(method, url string, params map[string]string, sendFiles []SendFile, header map[string]string) (int, map[string]interface{}, error) {
	pt := time.Now()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for _, sendFile := range sendFiles {
		for _, path := range sendFile.Paths {
			file, err := os.Open(path)
			if err != nil {
				log.Println("SendFormDataWithFilesRequest error: ", err)
				return 500, nil, err
			}
			defer file.Close()
			part, err := writer.CreateFormFile(sendFile.ParamName, filepath.Base(path))
			if err != nil {
				log.Println("SendFormDataWithFilesRequest create file error: ", err)
				return 500, nil, err
			}
			_, err = io.Copy(part, file)
			if err != nil {
				log.Println("SendFormDataWithFilesRequest file copy error: ", err)
				return 500, nil, err
			}
		}
	}
	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err := writer.Close()
	if err != nil {
		log.Println("SendFormDataWithFilesRequest writer error: ", err)
		return 500, nil, err
	}
	req, err := http.NewRequest(method, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	setHeader(req, header)
	if err != nil {
		log.Println("SendFormDataWithFilesRequest http new request error: ", err)
		return 500, nil, err
	}
	client := c.getClient()
	res, err := client.Do(req)
	defer c.closeIdleConnections(client)
	if err != nil {
		log.Println("SendFormDataWithFilesRequest http client do error: ", err)
		return 500, nil, err
	}
	showRespTimeLog(url, pt)
	return extractBody(res)
}

func (c *HttpClient) SendSoapRequest(method, url string, payload []byte, header map[string]string) (int, []byte, error) {
	pt := time.Now()

	// prepare the request
	req, err := http.NewRequest(method, url, bytes.NewReader(payload))
	if err != nil {
		log.Println("SendSoapRequest error creating request object", err)
		return 500, nil, err
	}

	// set the content type header, as well as the oter required headers
	req.Header.Set("Content-type", "text/xml;charset=utf-8")
	setHeader(req, header)
	// prepare the client request
	client := c.getClient()

	// dispatch the request
	res, err := client.Do(req)
	defer c.closeIdleConnections(client)
	if err != nil {
		log.Println("SendSoapRequest error http client do", err)
		return 500, nil, err
	}
	showRespTimeLog(url, pt)
	// read and parse the response body
	bodyBytes, err := ioutil.ReadAll(res.Body)
	return res.StatusCode, bodyBytes, err
}

func (c *HttpClient) SendFormDataRequest(method, url string, params map[string]string, header map[string]string) (int, map[string]interface{}, error) {
	pt := time.Now()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err := writer.Close()
	if err != nil {
		log.Println("SendFormDataRequest writer error: ", err)
		return 500, nil, err
	}
	req, err := http.NewRequest(method, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	setHeader(req, header)
	if err != nil {
		log.Println("SendFormDataRequest http new request error: ", err)
		return 500, nil, err
	}
	client := c.getClient()
	res, err := client.Do(req)
	defer c.closeIdleConnections(client)
	if err != nil {
		log.Println("SendFormDataRequest http client do error: ", err)
		return 500, nil, err
	}
	showRespTimeLog(url, pt)
	return extractBody(res)
}

func (c *HttpClient) SendQueryRequest(method, url string, params map[string]string, header map[string]string) (int, map[string]interface{}, error) {
	pt := time.Now()
	query := ""
	for key, val := range params {
		if query == "" {
			query += "?"
		} else {
			query += "&"
		}
		query += key + "=" + val
	}
	client := c.getClient()
	req, err := http.NewRequest(method, url+query, nil)
	if err != nil {
		log.Println("SendQueryRequest http client new request error:", err)
		return 500, nil, err
	}
	setHeader(req, header)
	res, err := client.Do(req)
	defer c.closeIdleConnections(client)
	if err != nil {
		log.Println("SendQueryRequest http client do error:", err)
		return 500, nil, err
	}
	defer res.Body.Close()
	showRespTimeLog(url, pt)
	return extractBody(res)
}

func (c *HttpClient) PostBodyRequest(url string, jsonStr string, header map[string]string) (int, map[string]interface{}, error) {
	return c.SendBodyRequest("POST", url, jsonStr, header)
}

func (c *HttpClient) PostFormDataWithFilesRequest(url string, params map[string]string, sendFiles []SendFile, header map[string]string) (int, map[string]interface{}, error) {
	return c.SendFormDataWithFilesRequest("POST", url, params, sendFiles, header)
}

func (c *HttpClient) PostSoapRequest(url string, payload []byte, header map[string]string) (int, []byte, error) {
	return c.SendSoapRequest("POST", url, payload, header)
}

func (c *HttpClient) PostFormDataRequest(url string, params map[string]string, header map[string]string) (int, map[string]interface{}, error) {
	return c.SendFormDataRequest("POST", url, params, header)
}

func (c *HttpClient) GetQueryRequest(url string, params map[string]string, header map[string]string) (int, map[string]interface{}, error) {
	return c.SendQueryRequest("GET", url, params, header)
}

func (c *HttpClient) CheckInternetConnected() bool {
	return c.CheckHttpServiceConnected("http://clients3.google.com/generate_204")
}

func (c *HttpClient) CheckHttpServiceConnected(httpUrl string) bool {
	client := c.getClient()
	req, err := http.NewRequest("GET", httpUrl, nil)
	if err != nil {
		return false
	}
	res, err := client.Do(req)
	defer c.closeIdleConnections(client)
	if err != nil {
		return false
	}
	defer res.Body.Close()
	return true
}

func (c *HttpClient) DownloadFile(url string, filepath string, header map[string]string) error {
	// download the file and check this url is ok
	client := c.getClient()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	setHeader(req, header)
	resp, err := client.Do(req)
	defer c.closeIdleConnections(client)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("status code is %d", resp.StatusCode)
	}
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()
	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func showRespTimeLog(logname string, ptime time.Time) {
	diff := float32(time.Now().UnixNano()-ptime.UnixNano()) / 1000000000
	if diff > 2 {
		log.Println(logname, "response time:", fmt.Sprintf("%.2f s", diff))
	}
}
