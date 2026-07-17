package services

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
)

func newTestReportResult() *ReportResult {
	return &ReportResult{
		ReportID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		Columns:    []ReportColumn{{Label: "Name"}, {Label: "Count"}},
		Data:       []map[string]interface{}{{"Name": "Alpha", "Count": 3}, {"Name": "Beta", "Count": 7}},
		ExecutedAt: time.Date(2026, 7, 17, 12, 0, 0, 0, time.UTC),
	}
}

func TestExportToJSON_Real(t *testing.T) {
	frs := &FormReportingService{}
	data, err := frs.exportToJSON(newTestReportResult())
	if err != nil {
		t.Fatalf("exportToJSON: %v", err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("exportToJSON produced invalid JSON: %v", err)
	}
}

func TestExportToExcel_RealXLSX(t *testing.T) {
	frs := &FormReportingService{}
	data, err := frs.exportToExcel(newTestReportResult())
	if err != nil {
		t.Fatalf("exportToExcel: %v", err)
	}
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("exportToExcel did not produce a valid XLSX zip: %v", err)
	}
	found := false
	for _, f := range zr.File {
		if f.Name == "xl/worksheets/sheet1.xml" {
			found = true
		}
	}
	if !found {
		t.Errorf("XLSX missing worksheet part")
	}
}

func TestExportToPDF_Real(t *testing.T) {
	frs := &FormReportingService{}
	data, err := frs.exportToPDF(newTestReportResult())
	if err != nil {
		t.Fatalf("exportToPDF: %v", err)
	}
	if !bytes.Contains(data, []byte("%PDF")) {
		t.Errorf("PDF export missing PDF header")
	}
}
