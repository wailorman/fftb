package remote_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/wailorman/fftb/pkg/distributed/handlers"
	"github.com/wailorman/fftb/pkg/distributed/models"
)

var authoritySecret = []byte("authority_secret_remote_dealer")
var sessionSecret = []byte("session_secret_remote_dealer")

type authorEq struct {
	author models.IAuthor
}

// Matches _
func (aeq authorEq) Matches(other interface{}) bool {
	otherAuthor := other.(models.IAuthor)

	return aeq.author.IsEqual(otherAuthor)
}

// String _
func (aeq authorEq) String() string {
	return fmt.Sprintf("%#v", aeq.author)
}

func AuthorEq(author models.IAuthor) gomock.Matcher {
	return authorEq{author: author}
}

// hasString _
type hasString struct {
	subStr string
}

func HasString(subStr string) gomock.Matcher {
	return &hasString{subStr: subStr}
}

// String _
func (hs hasString) String() string {
	return fmt.Sprintf("Contains %s", hs.subStr)
}

// Matches _
func (hs hasString) Matches(other interface{}) bool {
	var otherStr string

	if otherStrStringer, ok := other.(fmt.Stringer); ok {
		otherStr = otherStrStringer.String()
	}

	if otherStrStringer, ok := other.(error); ok {
		otherStr = otherStrStringer.Error()
	}

	return strings.Contains(otherStr, hs.subStr)
}

type echoClientWrap struct {
	e *echo.Echo
}

func newEchoClientWrap(e *echo.Echo) *echoClientWrap {
	return &echoClientWrap{
		e: e,
	}
}

// Do _
func (ew *echoClientWrap) Do(req *http.Request) (*http.Response, error) {
	rec := httptest.NewRecorder()
	ew.e.ServeHTTP(rec, req)
	return rec.Result(), nil
}

type storageClaim struct {
	id string
}

// GetID _
func (sc *storageClaim) GetID() string {
	return sc.id
}

// GetName _
func (sc *storageClaim) GetName() string {
	return sc.id
}

// GetURL _
func (sc *storageClaim) GetURL() string {
	return fmt.Sprintf("test://%s", sc.id)
}

// GetSize _
func (sc *storageClaim) GetSize() int {
	return 1
}

// GetType _
func (sc *storageClaim) GetType() string {
	return "test"
}

// WriteFrom _
func (sc *storageClaim) WriteFrom(io.Reader) error {
	return nil
}

// ReadTo _
func (sc *storageClaim) ReadTo(io.Writer) error {
	return nil
}

func createAuthor(t *testing.T) models.IAuthor {
	author := &models.Author{Name: "remote_dealer_test"}

	authorityToken, err := handlers.CreateAuthorityToken(authoritySecret, author.GetName())

	if err != nil {
		t.Fatalf("Failed to build author: %s", err)
		return nil
	}

	author.SetAuthorityKey(authorityToken)

	return author
}

func createStorageClaim(t *testing.T) models.IStorageClaim {
	return &storageClaim{id: uuid.New().String()}
}
