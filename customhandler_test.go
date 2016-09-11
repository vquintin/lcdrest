package lcdrest

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFirstPutLeadsToCreation(t *testing.T) {
	req, err := makePUTRequest("abc", "xyz")
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := makeHandler()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
}

func TestSecondPutLeadsToUpdate(t *testing.T) {
	req1, err := makePUTRequest("abc", "xyz")
	if err != nil {
		t.Fatal(err)
	}
	req2, err := makePUTRequest("abc", "rst")
	if err != nil {
		t.Fatal(err)
	}
	handler := makeHandler()
	rr1 := httptest.NewRecorder()
	rr2 := httptest.NewRecorder()

	handler.ServeHTTP(rr1, req1)
	handler.ServeHTTP(rr2, req2)

	assert.Equal(t, http.StatusOK, rr2.Code)
}

func TestPutRequestWithNoMessageLeadsToBadRequest(t *testing.T) {
	req, err := http.NewRequest("PUT", "/toto", strings.NewReader("Not valid!"))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := makeHandler()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestPutRequestWithTooManyMessagesLeadsToBadRequest(t *testing.T) {
	values := url.Values{
		"message": {"xyz", "rst"},
	}
	req, err := http.NewRequest("PUT", "/toto", strings.NewReader(values.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := makeHandler()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetRequestOnExistingKeyLeadsToSucess(t *testing.T) {
	req, err := http.NewRequest("GET", "/abc", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := makeHandlerWithValues()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "xyz", rr.Body.String())
}

func TestGetRequestOnNonExistingKeyLeadsToBadRequest(t *testing.T) {
	req, err := http.NewRequest("GET", "/golang", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := makeHandlerWithValues()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestDeleteRequestOnExistingKeyLeadsToSucess(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/abc", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := makeHandlerWithValues()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestDeleteRequestOnNonExistingKeyLeadsToBadRequest(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/golang", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := makeHandlerWithValues()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func makeHandlerWithValues() http.Handler {
	return makeAndInitializeHandler(map[string]string{
		"abc": "xyz",
	})
}

func makeHandler() http.Handler {
	return makeAndInitializeHandler(map[string]string{})
}

func makeAndInitializeHandler(messages map[string]string) http.Handler {
	var buf bytes.Buffer
	rm := NewRandomMessage(&buf, 15*time.Second)
	for k, v := range messages {
		rm.Put(k, v)
	}
	return NewCustomHandler(rm)
}

func makePUTRequest(key string, message string) (*http.Request, error) {
	values := url.Values{
		"message": {message},
	}
	return http.NewRequest("PUT", "/toto", strings.NewReader(values.Encode()))
}
