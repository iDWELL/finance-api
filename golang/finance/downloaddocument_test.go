package finance

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/domonda/go-types/uu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testDocID = uu.ID{0xa0, 0xee, 0xbc, 0x99, 0x9c, 0x0b, 0x4e, 0xf8, 0xbb, 0x6d, 0x6b, 0xb9, 0xbd, 0x38, 0x0a, 0x11}

func TestWithAuditTrail(t *testing.T) {
	t.Parallel()

	u, _ := url.Parse("http://example.com/doc.pdf")
	WithAuditTrail("append")(u)
	assert.Equal(t, "append", u.Query().Get("auditTrail"))
}

func TestWithAuditTrailLang(t *testing.T) {
	t.Parallel()

	u, _ := url.Parse("http://example.com/doc.pdf")
	WithAuditTrailLang("en")(u)
	assert.Equal(t, "en", u.Query().Get("auditTrailLang"))
}

func TestWithEmbedXML(t *testing.T) {
	t.Parallel()

	u, _ := url.Parse("http://example.com/doc.pdf")
	WithEmbedXML()(u)
	assert.Equal(t, "1", u.Query().Get("embedXML"))
}

func TestDownloadDocumentPDF_Success(t *testing.T) {
	t.Parallel()

	pdfContent := []byte("%PDF-1.4 fake content")

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/document/"+testDocID.String()+".pdf", r.URL.Path)
		assert.Contains(t, r.Header.Get("Authorization"), "Bearer ")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(pdfContent)
	}))
	defer srv.Close()

	ctx := WithBaseURL(t.Context(), srv.URL)
	file, err := DownloadDocumentPDF(ctx, "test-api-key", testDocID)
	require.NoError(t, err)
	assert.Equal(t, testDocID.String()+".pdf", file.Name())
	assert.Equal(t, pdfContent, file.FileData)
}

func TestDownloadDocumentPDF_WithOptions(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		assert.Equal(t, "append", q.Get("auditTrail"))
		assert.Equal(t, "en", q.Get("auditTrailLang"))
		assert.Equal(t, "1", q.Get("embedXML"))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("%PDF"))
	}))
	defer srv.Close()

	ctx := WithBaseURL(t.Context(), srv.URL)
	_, err := DownloadDocumentPDF(ctx, "key", testDocID,
		WithAuditTrail("append"),
		WithAuditTrailLang("en"),
		WithEmbedXML(),
	)
	require.NoError(t, err)
}

func TestDownloadDocumentPDF_ServerError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("internal error"))
	}))
	defer srv.Close()

	ctx := WithBaseURL(t.Context(), srv.URL)
	_, err := DownloadDocumentPDF(ctx, "key", testDocID)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "500")
}

func TestDownloadDocumentPDF_NetworkError(t *testing.T) {
	t.Parallel()

	ctx := WithBaseURL(t.Context(), "http://127.0.0.1:1")
	_, err := DownloadDocumentPDF(ctx, "key", testDocID)
	require.Error(t, err)
}
