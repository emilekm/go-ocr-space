package ocr_space

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func InitAPI(apiKey string, language string, options ApiOptions) OCRSpaceAPI {
	if options.Url == "" {
		options.Url = ocrDefaultUrl
	}

	if options.HTTPClient == nil {
		options.HTTPClient = http.DefaultClient
	}

	return OCRSpaceAPI{
		apiKey:   apiKey,
		language: language,
		options:  options,
	}
}

func (a *OCRSpaceAPI) ParseFromUrl(fileUrl string) (*OCRText, error) {
	var resp, err = a.postRequest(
		url.Values{
			"url":                          {fileUrl},
			"language":                     {a.language},
			"isOverlayRequired":            {"true"},
			"isSearchablePdfHideTextLayer": {"true"},
			"scale":                        {"true"},
		},
	)
	if err != nil {
		return nil, err
	}

	return unmarshalResponse(resp)
}

func (a *OCRSpaceAPI) ParseFromBase64(baseString string) (*OCRText, error) {
	resp, err := a.postRequest(
		url.Values{
			"base64Image":                  {baseString},
			"language":                     {a.language},
			"isOverlayRequired":            {"true"},
			"isSearchablePdfHideTextLayer": {"true"},
			"scale":                        {"true"},
		},
	)
	if err != nil {
		return nil, err
	}

	return unmarshalResponse(resp)
}

func (a *OCRSpaceAPI) ParseFromLocal(localPath string) (*OCRText, error) {
	params := map[string]string{
		"language":                     a.language,
		"isOverlayRequired":            "true",
		"isSearchablePdfHideTextLayer": "true",
		"scale":                        "true",
	}

	file, err := os.Open(localPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filepath.Base(localPath))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := a.preparePostRequest(body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := a.sendRequest(req)
	if err != nil {
		return nil, err
	}

	return unmarshalResponse(res)
}

func (a *OCRSpaceAPI) preparePostRequest(body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest("POST", ocrDefaultUrl, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("apikey", a.apiKey)

	return req, nil
}

func (a *OCRSpaceAPI) sendRequest(req *http.Request) (*http.Response, error) {
	return a.options.HTTPClient.Do(req)
}

func (a *OCRSpaceAPI) postRequest(values url.Values) (*http.Response, error) {
	req, err := a.preparePostRequest(strings.NewReader(values.Encode()))
	if err != nil {
		return nil, err
	}

	return a.sendRequest(req)
}

func (ocr OCRText) JustText() string {
	text := ""
	if ocr.IsErroredOnProcessing {
		for _, page := range ocr.ErrorMessage {
			text += page
		}
	} else {
		for _, page := range ocr.ParsedResults {
			text += page.ParsedText
		}
	}
	return text
}

func unmarshalResponse(res *http.Response) (*OCRText, error) {
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var result OCRText
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
