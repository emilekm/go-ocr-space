package ocr_space

import "net/http"

const (
	ocrDefaultUrl = "https://api.ocr.space/parse/image"
)

type OCRText struct {
	ParsedResults []struct {
		TextOverlay struct {
			Lines []struct {
				Words []struct {
					WordText string  `json:"WordText"`
					Left     float64 `json:"Left"`
					Top      float64 `json:"Top"`
					Height   float64 `json:"Height"`
					Width    float64 `json:"Width"`
				} `json:"Words"`

				MaxHeight float64 `json:"MaxHeight"`
				MinTop    float64 `json:"MinTop"`
			} `json:"Lines"`

			HasOverlay bool   `json:"HasOverlay"`
			Message    string `json:"Message"`
		} `json:"TextOverlay"`

		TextOrientation   string `json:"TextOrientation"`
		FileParseExitCode int    `json:"FileParseExitCode"`
		ParsedText        string `json:"ParsedText"`
		ErrorMessage      string `json:"ErrorMessage"`
		ErrorDetails      string `json:"ErrorDetails"`
	} `json:"ParsedResults"`

	OCRExitCode                  int      `json:"OCRExitCode"`
	IsErroredOnProcessing        bool     `json:"IsErroredOnProcessing"`
	ErrorMessage                 []string `json:"ErrorMessage"`
	ErrorDetails                 string   `json:"ErrorDetails"`
	ProcessingTimeInMilliseconds string   `json:"ProcessingTimeInMilliseconds"`
	SearchablePDFURL             string   `json:"SearchablePDFURL"`
}

type OCRSpaceAPI struct {
	apiKey  string
	options ApiOptions
}

type ApiOptions struct {
	Url        string
	HTTPClient *http.Client
}

type Lang string

const (
	LangArabic             Lang = "ara"
	LangBulgarian          Lang = "bul"
	LangChineseSimplified  Lang = "chs"
	LangChineseTraditional Lang = "cht"
	LangCroatian           Lang = "hrv"
	LangCzech              Lang = "cze"
	LangDanish             Lang = "dan"
	LangDutch              Lang = "dut"
	LangEnglish            Lang = "eng"
	LangFinnish            Lang = "fin"
	LangFrench             Lang = "fre"
	LangGerman             Lang = "ger"
	LangGreek              Lang = "gre"
	LangHungarian          Lang = "hun"
	LangKorean             Lang = "kor"
	LangItalian            Lang = "ita"
	LangJapanese           Lang = "jpn"
	LangPolish             Lang = "pol"
	LangPortuguese         Lang = "por"
	LangRussian            Lang = "rus"
	LangSlovenian          Lang = "slv"
	LangSpanish            Lang = "spa"
	LangSwedish            Lang = "swe"
	LangTurkish            Lang = "tur"
)

type Filetype string

const (
	FiletypePDF Filetype = "PDF"
	FiletypeGIF Filetype = "GIF"
	FiletypePNG Filetype = "PNG"
	FiletypeJPG Filetype = "JPG"
	FiletypeTIF Filetype = "TIF"
	FiletypeBMP Filetype = "BMP"
)

type OCREngineVer int

const (
	OCREngineV1 OCREngineVer = 1
	OCREngineV2 OCREngineVer = 2
)

type Params struct {
	// Language used for OCR.
	// Default = eng
	Language Lang `url:"language,omitempty"`

	// If true, returns the coordinates of the bounding boxes for each word.
	// If false, the OCR'ed text is returned only as a text block (this makes the JSON reponse smaller).
	// Default = False
	IsOverlayRequired bool `url:"isOverlayRequired"`

	// Overwrites the automatic file type detection based on content-type.
	// Supported image file formats are png, jpg (jpeg), gif, tif (tiff) and bmp.
	// For document ocr, the api supports the Adobe PDF format. Multi-page TIFF files are supported.
	Filetype Filetype `url:"filetype,omitempty"`

	// If set to true, the api autorotates the image correctly and sets the TextOrientation parameter in the JSON response.
	// If the image is not rotated, then TextOrientation=0, otherwise it is the degree of the rotation, e. g. "270".
	// Default = False
	DetectOrientation bool `url:"detectOrientation"`

	// If true, API generates a searchable PDF.
	// This parameter automatically sets isOverlayRequired = true.
	// Default = False
	IsCreateSearchablePDF bool `url:"isCreateSearchablePdf"`

	// If true, the text layer is hidden (not visible)
	// Default = False
	IsSearchablePDFHideTextLayer bool `url:"isSearchablePdfHideTextLayer"`

	// If set to true, the api does some internal upscaling.
	// This can improve the OCR result significantly, especially for low-resolution PDF scans.
	// Default = False
	Scale bool `url:"scale"`

	// If set to true, the OCR logic makes sure that the parsed text result is always returned line by line.
	// This switch is recommended for table OCR, receipt OCR, invoice processing and all other type of input documents that have a table like structure.
	// Default = False
	IsTable bool `url:"isTable"`

	// OCR engine version: 1 or 2
	// Default = 1
	OCREngine *OCREngineVer `url:"OCREngine,omitempty"`
}
