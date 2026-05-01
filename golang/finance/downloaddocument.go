package finance

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/domonda/go-types/uu"
	fs "github.com/ungerik/go-fs"
)

// DownloadOption modifies the PDF download request URL.
type DownloadOption func(u *url.URL)

// WithAuditTrail appends an audit trail to the downloaded PDF.
// Valid values: "append", "prepend", "configured", "only".
func WithAuditTrail(trailType string) DownloadOption {
	return func(u *url.URL) {
		q := u.Query()
		q.Add("auditTrail", trailType)
		u.RawQuery = q.Encode()
	}
}

// WithAuditTrailLang sets the language for the audit trail (e.g. "en", "de").
func WithAuditTrailLang(lang string) DownloadOption {
	return func(u *url.URL) {
		q := u.Query()
		q.Add("auditTrailLang", lang)
		u.RawQuery = q.Encode()
	}
}

// WithEmbedXML embeds the UN/CEFACT format XML into the downloaded PDF.
func WithEmbedXML() DownloadOption {
	return func(u *url.URL) {
		q := u.Query()
		q.Add("embedXML", "1")
		u.RawQuery = q.Encode()
	}
}

// DownloadDocumentPDF downloads the PDF file for a document identified by docID.
// The returned fs.MemFile is named "<docID>.pdf" and contains the PDF bytes.
//
// Use DownloadOption functions to modify the request:
//   - WithAuditTrail("append"|"prepend"|"configured"|"only")
//   - WithAuditTrailLang("en"|"de"|...)
//   - WithEmbedXML()
//
// API endpoint: https://idwell.ai/api/public/document/{docID}.pdf
func DownloadDocumentPDF(ctx context.Context, apiKey string, docID uu.ID, opts ...DownloadOption) (fs.MemFile, error) {
	pdfURL, err := url.Parse(fmt.Sprintf("%s/document/%s.pdf", baseURLFromCtx(ctx), docID))
	if err != nil {
		return fs.MemFile{}, err
	}

	for _, opt := range opts {
		opt(pdfURL)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, pdfURL.String(), nil)
	if err != nil {
		return fs.MemFile{}, err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := httpClientFromCtx(ctx).Do(req) //nolint:gosec // intentional HTTP call to API URL from context
	if err != nil {
		return fs.MemFile{}, err
	}

	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fs.MemFile{}, fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fs.MemFile{}, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, body)
	}

	return fs.NewMemFile(fmt.Sprintf("%s.pdf", docID), body), nil
}
