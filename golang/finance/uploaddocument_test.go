package finance

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/domonda/go-types/uu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	fs "github.com/ungerik/go-fs"
)

var testCatID = uu.ID{0xb1, 0xee, 0xbc, 0x99, 0x9c, 0x0b, 0x4e, 0xf8, 0xbb, 0x6d, 0x6b, 0xb9, 0xbd, 0x38, 0x0a, 0x22}

func TestUploadDocument_Success(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/upload", r.URL.Path)
		assert.Contains(t, r.Header.Get("Authorization"), "Bearer ")
		assert.Contains(t, r.Header.Get("Content-Type"), "multipart/form-data")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(testDocID.String()))
	}))
	defer srv.Close()

	ctx := WithBaseURL(t.Context(), srv.URL)
	got, err := UploadDocument(ctx, "test-key", testCatID, fs.NewMemFile("invoice.pdf", []byte("%PDF")), nil)
	require.NoError(t, err)
	assert.Equal(t, testDocID, got)
}

func TestUploadDocument_WithInvoiceFile(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.NoError(t, r.ParseMultipartForm(10<<20))
		assert.NotNil(t, r.MultipartForm.File["invoice"])
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(testDocID.String()))
	}))
	defer srv.Close()

	ctx := WithBaseURL(t.Context(), srv.URL)
	_, err := UploadDocument(ctx, "key", testCatID,
		fs.NewMemFile("invoice.pdf", []byte("%PDF")),
		fs.NewMemFile("invoice.json", []byte(`{"invoiceNumber":"INV-001"}`)),
	)
	require.NoError(t, err)
}

func TestUploadDocument_WithTags(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.NoError(t, r.ParseMultipartForm(10<<20))
		assert.Equal(t, []string{"invoice", "2024"}, r.MultipartForm.Value["tag"])
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(testDocID.String()))
	}))
	defer srv.Close()

	ctx := WithBaseURL(t.Context(), srv.URL)
	_, err := UploadDocument(ctx, "key", testCatID,
		fs.NewMemFile("invoice.pdf", []byte("%PDF")),
		nil,
		"invoice", "2024",
	)
	require.NoError(t, err)
}

func TestUploadDocument_ErrorStatus(t *testing.T) {
	t.Parallel()

	for _, status := range []int{400, 409, 500} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			t.Parallel()

			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(status)
			}))
			defer srv.Close()

			ctx := WithBaseURL(t.Context(), srv.URL)
			_, err := UploadDocument(ctx, "key", testCatID, fs.NewMemFile("f.pdf", []byte("%PDF")), nil)
			require.Error(t, err)
		})
	}
}

func TestUploadDocument_NetworkError(t *testing.T) {
	t.Parallel()

	ctx := WithBaseURL(t.Context(), "http://127.0.0.1:1")
	_, err := UploadDocument(ctx, "key", testCatID, fs.NewMemFile("f.pdf", []byte("%PDF")), nil)
	require.Error(t, err)
}
