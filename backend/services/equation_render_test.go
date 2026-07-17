package services

import (
	"bytes"
	"image/png"
	"strings"
	"testing"
)

func TestRenderEquationImageSVG(t *testing.T) {
	data, ct, err := RenderEquationImage("E = mc^{2}", "svg", "medium", "#000000")
	if err != nil {
		t.Fatalf("svg render: %v", err)
	}
	if ct != "image/svg+xml" {
		t.Errorf("content type = %q", ct)
	}
	if !bytes.Contains(data, []byte("<svg")) || !bytes.Contains(data, []byte("E = mc^{2}")) {
		t.Errorf("svg output missing expected markup: %s", data)
	}
}

func TestRenderEquationImagePDF(t *testing.T) {
	data, ct, err := RenderEquationImage("a^2+b^2=c^2", "pdf", "medium", "")
	if err != nil {
		t.Fatalf("pdf render: %v", err)
	}
	if ct != "application/pdf" {
		t.Errorf("content type = %q", ct)
	}
	if !bytes.HasPrefix(data, []byte("%PDF")) {
		lim := 8
		if len(data) < lim {
			lim = len(data)
		}
		t.Errorf("pdf output not a valid PDF header: %s", data[:lim])
	}
}

func TestRenderEquationImagePNG(t *testing.T) {
	data, ct, err := RenderEquationImage("x+1", "png", "small", "")
	if err != nil {
		t.Fatalf("png render: %v", err)
	}
	if ct != "image/png" {
		t.Errorf("content type = %q", ct)
	}
	if _, err := png.Decode(bytes.NewReader(data)); err != nil {
		t.Errorf("png output is not a decodable PNG: %v", err)
	}
}

func TestRenderEquationImageUnknownFormat(t *testing.T) {
	if _, _, err := RenderEquationImage("x", "gif", "medium", ""); err == nil {
		t.Errorf("expected error for unsupported format")
	}
}

func TestStrokesIntersect(t *testing.T) {
	hrs := &HandwritingRecognitionService{}
	s1 := HandwritingStroke{X: []float64{0, 10}, Y: []float64{0, 10}}
	s2 := HandwritingStroke{X: []float64{5, 15}, Y: []float64{5, 15}}
	if !hrs.strokesIntersect(s1, s2) {
		t.Errorf("overlapping strokes should intersect")
	}
	s3 := HandwritingStroke{X: []float64{100, 110}, Y: []float64{100, 110}}
	if hrs.strokesIntersect(s1, s3) {
		t.Errorf("disjoint strokes should not intersect")
	}
}

func TestBuildComplexExpressionReflectsStrokes(t *testing.T) {
	hrs := &HandwritingRecognitionService{}
	out := hrs.buildComplexExpression([]HandwritingStroke{
		{X: []float64{0, 10, 20}, Y: []float64{0, 10, 0}},
	})
	if !strings.Contains(out, "strokes:1") {
		t.Errorf("buildComplexExpression output does not reflect input: %q", out)
	}
	if out == "x^{2} + y^{2} = z^{2}" {
		t.Errorf("buildComplexExpression still returns hardcoded constant")
	}
}
