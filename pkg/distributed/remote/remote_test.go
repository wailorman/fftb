package remote_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
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
