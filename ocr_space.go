package ocr_space

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

func InitConfig(apiKey string, language string) Config {
	var config Config
	config.ApiKey = apiKey
	config.Language = language
	config.Url = "https://api.ocr.space/parse/image"
	return config
}

func (c Config) ParseFromUrl(fileUrl string) (OCRText, error) {
	var results OCRText
	var resp, err = http.PostForm(c.Url,
		url.Values{
			"url":                          {fileUrl},
			"language":                     {c.Language},
			"apikey":                       {c.ApiKey},
			"isOverlayRequired":            {"true"},
			"isSearchablePdfHideTextLayer": {"true"},
			"scale": {"true"},
		},
	)
	if err != nil {
		return results, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return results, err
	}

	err = json.Unmarshal(body, &results)
	if err != nil {
		return results, err
	}

	return results, nil
}

func (c Config) ParseFromBase64(baseString string) (OCRText, error) {
	var results OCRText
	resp, err := http.PostForm("https://api.ocr.space/parse/image",
		url.Values{
			"base64Image":                  {baseString},
			"language":                     {c.Language},
			"apikey":                       {c.ApiKey},
			"isOverlayRequired":            {"true"},
			"isSearchablePdfHideTextLayer": {"true"},
			"scale": {"true"},
		},
	)
	if err != nil {
		return results, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return results, err
	}

	err = json.Unmarshal(body, &results)
	if err != nil {
		return results, err
	}

	return results, nil
}

func (c Config) ParseFromLocal(localPath string) (OCRText, error) {
	var results OCRText
	params := map[string]string{
		"language":                     c.Language,
		"apikey":                       c.ApiKey,
		"isOverlayRequired":            "true",
		"isSearchablePdfHideTextLayer": "true",
		"scale": "true",
	}

	file, err := os.Open(localPath)
	if err != nil {
		return results, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filepath.Base(localPath))
	if err != nil {
		return results, err
	}
	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return results, err
	}

	req, err := http.NewRequest("POST", c.Url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	response, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	} else {
		body := &bytes.Buffer{}
		_, err := body.ReadFrom(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		response.Body.Close()
		err = json.Unmarshal(body.Bytes(), &results)
		if err != nil {
			return results, err
		}
	}

	return results, nil
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
