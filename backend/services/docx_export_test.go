package services

import (
	"archive/zip"
	"bytes"
	"testing"
)

func TestDOCXExportWithHeaderFooter(t *testing.T) {
	des := NewDOCXExportService()
	content := &DocumentContent{
		Title:  "Test",
		Header: "Confidential Header",
		Footer: "Page Footer Text",
		Content: []ContentBlock{
			{Type: "paragraph", Text: "Hello world"},
		},
	}

	data, err := des.ExportToDOCX(content)
	if err != nil {
		t.Fatalf("ExportToDOCX failed: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("generated DOCX is not a valid ZIP: %v", err)
	}

	want := map[string]bool{
		"word/document.xml": true,
		"word/header1.xml":  true,
		"word/footer1.xml":  true,
	}
	for _, f := range zr.File {
		delete(want, f.Name)
	}
	if len(want) > 0 {
		t.Errorf("missing expected DOCX parts: %v", want)
	}
}

func TestDOCXExportWithoutHeaderFooter(t *testing.T) {
	des := NewDOCXExportService()
	content := &DocumentContent{
		Title: "Test",
		Content: []ContentBlock{
			{Type: "paragraph", Text: "Hello world"},
		},
	}

	data, err := des.ExportToDOCX(content)
	if err != nil {
		t.Fatalf("ExportToDOCX failed: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("generated DOCX is not a valid ZIP: %v", err)
	}
	for _, f := range zr.File {
		if f.Name == "word/header1.xml" || f.Name == "word/footer1.xml" {
			t.Errorf("header/footer part should not be present: %s", f.Name)
		}
	}
}
