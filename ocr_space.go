package ocr_space

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/go-querystring/query"
)

func InitAPI(apiKey string, options ApiOptions) OCRSpaceAPI {
	if options.Url == "" {
		options.Url = ocrDefaultUrl
	}

	if options.HTTPClient == nil {
		options.HTTPClient = http.DefaultClient
	}

	return OCRSpaceAPI{
		apiKey:  apiKey,
		options: options,
	}
}

func (a *OCRSpaceAPI) ParseFromUrl(fileUrl string, params Params) (*OCRText, error) {
	values, err := query.Values(params)
	if err != nil {
		return nil, err
	}

	values.Add("url", fileUrl)

	res, err := a.postRequest(values)
	if err != nil {
		return nil, err
	}

	return unmarshalResponse(res)
}

func (a *OCRSpaceAPI) ParseFromBase64(baseString string, params Params) (*OCRText, error) {
	values, err := query.Values(params)
	if err != nil {
		return nil, err
	}

	values.Add("base64Image", baseString)

	res, err := a.postRequest(values)
	if err != nil {
		return nil, err
	}

	return unmarshalResponse(res)
}

func (a *OCRSpaceAPI) ParseFromLocal(file File, params Params) (*OCRText, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", file.Name)
	if err != nil {
		return nil, err
	}

	_, err = part.Write(file.Content)
	if err != nil {
		return nil, err
	}

	values, err := query.Values(params)
	if err != nil {
		return nil, err
	}
	for key := range values {
		err = writer.WriteField(key, values.Get(key))
		if err != nil {
			return nil, err
		}
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
