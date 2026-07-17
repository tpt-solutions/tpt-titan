package services

import (
	"bytes"
	"encoding/base64"
	"io"
	"mime"
	"mime/multipart"
	"net/mail"
	"net/textproto"
	"testing"
)

func bodyHeader(contentType, disp, filename string) textproto.MIMEHeader {
	h := textproto.MIMEHeader{}
	h.Set("Content-Type", contentType)
	if disp != "" {
		if filename != "" {
			h.Set("Content-Disposition", disp+`; filename="`+filename+`"`)
		} else {
			h.Set("Content-Disposition", disp)
		}
	}
	return h
}

func buildMultipartMessage(t *testing.T) *mail.Message {
	t.Helper()
	var buf bytes.Buffer
	boundary := "BOUNDARY42"
	buf.WriteString("From: sender@example.com\r\n")
	buf.WriteString("To: receiver@example.com\r\n")
	buf.WriteString("Subject: Test\r\n")
	buf.WriteString("MIME-Version: 1.0\r\n")
	buf.WriteString("Content-Type: multipart/mixed; boundary=" + boundary + "\r\n")
	buf.WriteString("\r\n")

	w := multipart.NewWriter(&buf)
	if err := w.SetBoundary(boundary); err != nil {
		t.Fatalf("set boundary: %v", err)
	}

	tw, _ := w.CreatePart(bodyHeader("text/plain", "inline", ""))
	tw.Write([]byte("Hello, see attached files."))

	aw1, _ := w.CreatePart(bodyHeader("application/pdf", "attachment", "report.pdf"))
	aw1.Write([]byte("%PDF-1.4 fake pdf contents for testing"))

	aw2, _ := w.CreatePart(bodyHeader("image/png", "attachment", "logo.png"))
	aw2.Write([]byte("PNGDATA-bytes-here"))

	if err := w.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}

	msg, err := mail.ReadMessage(&buf)
	if err != nil {
		t.Fatalf("read message: %v", err)
	}
	return msg
}

type parsedPart struct {
	filename string
	ct       string
	inline   bool
	data     string
}

// extractAttachmentParts mirrors the MIME-walking logic used by
// ProcessEmailAttachments but returns the parsed parts directly so the parsing
// can be verified without a database.
func extractAttachmentParts(t *testing.T, email *mail.Message) []parsedPart {
	t.Helper()
	var out []parsedPart

	mediaType, params, err := mime.ParseMediaType(email.Header.Get("Content-Type"))
	if err != nil || mediaType == "" {
		return out
	}
	mr := multipart.NewReader(email.Body, params["boundary"])
	for {
		part, perr := mr.NextPart()
		if perr != nil {
			break
		}
		ct := part.Header.Get("Content-Type")
		pmt, _, _ := mime.ParseMediaType(ct)
		if isAttachmentHeader(part.Header) {
			data, _ := io.ReadAll(part)
			out = append(out, parsedPart{
				filename: decodeFilename(part.FileName()),
				ct:       pmt,
				inline:   isInlineHeader(part.Header),
				data:     string(data),
			})
		}
	}
	return out
}

func TestProcessEmailAttachmentsRealParsing(t *testing.T) {
	msg := buildMultipartMessage(t)
	parts := extractAttachmentParts(t, msg)
	if len(parts) != 2 {
		t.Fatalf("expected 2 attachments, got %d", len(parts))
	}

	byName := map[string]parsedPart{}
	for _, p := range parts {
		byName[p.filename] = p
	}
	a, ok := byName["report.pdf"]
	if !ok {
		t.Fatalf("report.pdf attachment missing; got %v", names(parts))
	}
	if a.ct != "application/pdf" {
		t.Errorf("report.pdf content type = %q, want application/pdf", a.ct)
	}
	if !bytes.Contains([]byte(a.data), []byte("fake pdf contents")) {
		t.Errorf("report.pdf data not parsed correctly: %q", a.data)
	}
	if _, ok := byName["logo.png"]; !ok {
		t.Errorf("logo.png attachment missing")
	}
}

func TestDecodeFilenameStripsPathsAndEncodedWords(t *testing.T) {
	if got := decodeFilename(`../evil/report.pdf`); got != "report.pdf" {
		t.Errorf("path not stripped: %q", got)
	}
	if got := decodeFilename(`=?UTF-8?B?` + base64.StdEncoding.EncodeToString([]byte("résumé.pdf")) + `?=`); got != "résumé.pdf" {
		t.Errorf("rfc2047 not decoded: %q", got)
	}
	if got := decodeFilename(""); got != "attachment.bin" {
		t.Errorf("empty filename fallback wrong: %q", got)
	}
}

func TestGenerateImageThumbnail(t *testing.T) {
	pngData := mustDecodeBase64("iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAIAAACQd1PeAAAAEElEQVR4nGL6z8AACAAA//8DCQECWLbVUAAAAABJRU5ErkJggg==")
	if _, err := generateImageThumbnail(pngData, 32, 32); err != nil {
		t.Fatalf("generateImageThumbnail: %v", err)
	}
}

func mustDecodeBase64(s string) []byte {
	data := make([]byte, base64.StdEncoding.DecodedLen(len(s)))
	n, err := base64.StdEncoding.Decode(data, []byte(s))
	if err != nil {
		panic(err)
	}
	return data[:n]
}

func names(parts []parsedPart) []string {
	var ns []string
	for _, p := range parts {
		ns = append(ns, p.filename)
	}
	return ns
}
