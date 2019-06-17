package ocr_space

import "testing"

func (ocr OCRText) TestJustText(t *testing.T) {
	want := "hello world"
	ocr.IsErroredOnProcessing = true
	ocr.ErrorMessage = []string{"hello", "world"}
	if got := ocr.JustText(); got != want {
		t.Errorf("Hello() = %q, want %q", got, want)
	}
}
