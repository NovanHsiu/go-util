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

var RequestTimeOutSecond = 30

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

func PostBodyRequest(url string, jsonStr string) (int, map[string]interface{}, error) {
	pt := time.Now()
	jsonBytes := []byte(jsonStr)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		log.Println("PostBodyRequest http client new request error:", err)
		return 500, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connection", "application/json")

	client := &http.Client{
		Timeout: time.Duration(time.Duration(RequestTimeOutSecond) * time.Second),
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	res, err := client.Do(req)
	defer client.CloseIdleConnections()
	if err != nil {
		log.Println("PostBodyRequest http client do error:", err)
		return 500, nil, err
	}
	defer res.Body.Close()
	showRespTimeLog(url, pt)
	return extractBody(res)
}

type PostFile struct {
	ParamName string
	Paths     []string
}

func PostFormDataWithFilesRequest(url string, params map[string]string, postFiles []PostFile) (int, map[string]interface{}, error) {
	pt := time.Now()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for _, postFile := range postFiles {
		for _, path := range postFile.Paths {
			file, err := os.Open(path)
			if err != nil {
				log.Println("PostFormDataWithFilesRequest error: ", err)
				return 500, nil, err
			}
			defer file.Close()
			part, err := writer.CreateFormFile(postFile.ParamName, filepath.Base(path))
			if err != nil {
				log.Println("PostFormDataWithFilesRequest create file error: ", err)
				return 500, nil, err
			}
			_, err = io.Copy(part, file)
			if err != nil {
				log.Println("PostFormDataWithFilesRequest file copy error: ", err)
				return 500, nil, err
			}
		}
	}
	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err := writer.Close()
	if err != nil {
		log.Println("PostFormDataWithFilesRequest writer error: ", err)
		return 500, nil, err
	}
	req, err := http.NewRequest("POST", url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if err != nil {
		log.Println("PostFormDataWithFilesRequest http new request error: ", err)
		return 500, nil, err
	}
	client := &http.Client{
		Timeout: time.Duration(time.Duration(RequestTimeOutSecond) * time.Second),
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	res, err := client.Do(req)
	defer client.CloseIdleConnections()
	if err != nil {
		log.Println("PostFormDataWithFilesRequest http client do error: ", err)
		return 500, nil, err
	}
	showRespTimeLog(url, pt)
	return extractBody(res)
}

func PostSoapRequest(url string, payload []byte) (int, []byte, error) {
	pt := time.Now()
	httpMethod := "POST"

	// prepare the request
	req, err := http.NewRequest(httpMethod, url, bytes.NewReader(payload))
	if err != nil {
		log.Println("PostSoapRequest error creating request object", err)
		return 500, nil, err
	}

	// set the content type header, as well as the oter required headers
	req.Header.Set("Content-type", "text/xml;charset=utf-8")

	// prepare the client request
	client := &http.Client{
		Timeout: time.Duration(time.Duration(RequestTimeOutSecond) * time.Second),
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	// dispatch the request
	res, err := client.Do(req)
	defer client.CloseIdleConnections()
	if err != nil {
		log.Println("PostSoapRequest error http client do", err)
		return 500, nil, err
	}
	showRespTimeLog(url, pt)
	// read and parse the response body
	bodyBytes, err := ioutil.ReadAll(res.Body)
	return res.StatusCode, bodyBytes, err
}

func PostFormDataRequest(url string, params map[string]string) (int, map[string]interface{}, error) {
	pt := time.Now()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err := writer.Close()
	if err != nil {
		log.Println("PostFormDataRequest writer error: ", err)
		return 500, nil, err
	}
	req, err := http.NewRequest("POST", url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if err != nil {
		log.Println("PostFormDataRequest http new request error: ", err)
		return 500, nil, err
	}
	client := &http.Client{
		Timeout: time.Duration(time.Duration(RequestTimeOutSecond) * time.Second),
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	res, err := client.Do(req)
	defer client.CloseIdleConnections()
	if err != nil {
		log.Println("PostFormDataRequest http client do error: ", err)
		return 500, nil, err
	}
	showRespTimeLog(url, pt)
	return extractBody(res)
}

func GetQueryRequest(url string, params map[string]string) (int, map[string]interface{}, error) {
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
	client := &http.Client{
		Timeout: time.Duration(time.Duration(RequestTimeOutSecond) * time.Second),
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	req, err := http.NewRequest("GET", url+query, nil)
	if err != nil {
		log.Println("GetQueryRequest http client new request error:", err)
		return 500, nil, err
	}
	//req.Header.Add("Authorization", "")
	res, err := client.Do(req)
	defer client.CloseIdleConnections()
	if err != nil {
		log.Println("GetQueryRequest http client do error:", err)
		return 500, nil, err
	}
	defer res.Body.Close()
	showRespTimeLog(url, pt)
	return extractBody(res)
}

func CheckInternetConnected() bool {
	return CheckHttpServiceConnected("http://clients3.google.com/generate_204", RequestTimeOutSecond)
}

func CheckHttpServiceConnected(httpUrl string, timeOutSeconds int) bool {
	client := &http.Client{
		Timeout: time.Duration(time.Duration(timeOutSeconds) * time.Second),
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	req, err := http.NewRequest("GET", httpUrl, nil)
	if err != nil {
		return false
	}
	res, err := client.Do(req)
	defer client.CloseIdleConnections()
	if err != nil {
		return false
	}
	defer res.Body.Close()
	return true
}

func DownloadFile(url string, filepath string) error {
	// download the file and check this url is ok
	client := &http.Client{
		Timeout: time.Duration(time.Duration(RequestTimeOutSecond) * time.Second),
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	//req.Header.Add("Authorization", "")
	resp, err := client.Do(req)
	defer client.CloseIdleConnections()
	if err != nil {
		return err
	}
	defer resp.Body.Close()
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
