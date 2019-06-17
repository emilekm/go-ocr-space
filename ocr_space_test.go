package ocr_space

import (
	"testing"
)

func TestJustText(t *testing.T) {
	var ocr OCRText
	want := "helloworld"
	ocr.IsErroredOnProcessing = true
	ocr.ErrorMessage = []string{"hello", "world"}
	got := ocr.JustText()
	if got != want {
		t.Errorf("JustText() = %q, want %q", got, want)
	}
}
