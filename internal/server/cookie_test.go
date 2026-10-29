package server

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"gotest.tools/v3/assert"
)

func TestResendAuthCookie(t *testing.T) {
	baseDomain := "example.com"
	req := &http.Request{
		Host:   fmt.Sprintf("dev.%s", baseDomain),
		Header: make(http.Header),
		TLS:    &tls.ConnectionState{}, // if this isnt nil the request is considered secure
	}
	maxAge := int(time.Until(time.Now().Add(5 * time.Second)).Seconds())
	req.AddCookie(&http.Cookie{
		Name:     cookieSignupName,
		Value:    "aaa",
		MaxAge:   maxAge,
		Path:     cookiePath,
		Domain:   baseDomain,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
		HttpOnly: true,
	})

	resp := httptest.NewRecorder()
	bearer := exchangeSignupCookieForSession(
		req,
		resp,
		Options{SessionDuration: 1 * time.Minute, BaseDomain: baseDomain})
	assert.Equal(t, "aaa", bearer)
	assert.Equal(t, len(resp.Header()["Set-Cookie"]), 2)
	cookies := resp.Header()["Set-Cookie"]

	matched, err := regexp.MatchString(
		"auth=aaa; Path=/; Domain=dev.example.com; Max-Age=\\d\\d; HttpOnly; Secure; SameSite=Strict",
		cookies[0])
	assert.NilError(t, err)
	assert.Assert(t, matched)

	assert.Equal(t, "signup=; Path=/; Domain=example.com; Max-Age=0; HttpOnly; Secure; SameSite=Strict", cookies[1])
}

func TestCookieSecureFlagIsSetForTLS(t *testing.T) {
	baseDomain := "example.com"
	req := &http.Request{
		Host:   fmt.Sprintf("dev.%s", baseDomain),
		Header: make(http.Header),
		TLS:    &tls.ConnectionState{}, // if this isnt nil the request is considered secure
	}

	resp := httptest.NewRecorder()
	setCookie(
		req,
		resp,
		cookieConfig{
			Name:  "tls-cookie-test",
			Value: "1234",
		})
	assert.Equal(t, len(resp.Header()["Set-Cookie"]), 1)
	cookies := resp.Header()["Set-Cookie"]

	assert.Equal(t, "tls-cookie-test=1234; Path=/; Max-Age=0; HttpOnly; Secure; SameSite=Strict", cookies[0])
}

func TestSessionCookiesSecureBehindTLSProxy(t *testing.T) {
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com/", nil)
	resp := httptest.NewRecorder()
	setCookie(req, resp, cookieConfig{Name: cookieAuthorizationName, Value: "secret", Expires: time.Now().Add(time.Minute)})
	deleteCookie(req, resp, cookieSignupName, "example.com")
	for _, cookie := range resp.Result().Cookies() {
		assert.Assert(t, cookie.Secure)
		assert.Assert(t, cookie.HttpOnly)
		assert.Equal(t, cookie.SameSite, http.SameSiteStrictMode)
	}
}
