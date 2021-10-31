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
)

func InitAPI(apiKey string, language string, options ApiOptions) OCRSpaceAPI {
	if options.Url == "" {
		options.Url = ocrDefaultUrl
	}

	return OCRSpaceAPI{
		apiKey:   apiKey,
		language: language,
		options:  options,
	}
}

func (c OCRSpaceAPI) ParseFromUrl(fileUrl string) (*OCRText, error) {
	var resp, err = http.PostForm(c.options.Url,
		url.Values{
			"url":                          {fileUrl},
			"language":                     {c.language},
			"apikey":                       {c.apiKey},
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

func (c OCRSpaceAPI) ParseFromBase64(baseString string) (*OCRText, error) {
	resp, err := http.PostForm(c.options.Url,
		url.Values{
			"base64Image":                  {baseString},
			"language":                     {c.language},
			"apikey":                       {c.apiKey},
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

func (c OCRSpaceAPI) ParseFromLocal(localPath string) (*OCRText, error) {
	params := map[string]string{
		"language":                     c.language,
		"apikey":                       c.apiKey,
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

	req, err := http.NewRequest("POST", c.options.Url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return unmarshalResponse(response)
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
