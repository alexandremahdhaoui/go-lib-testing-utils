package tUtils

import (
	"crypto/tls"
	"net/http"
	"testing"
)

type Tester interface {
	T() *testing.T
}

type Identifier interface {
	Id() string
}

type HttpConfig interface {
	HttpClient() *http.Client
	Tester
	TlsConfig() *tls.Config
}

type TlsTester interface {
	Tester
	TlsConfig() *tls.Config
}

type HttpTester interface {
	HttpClient() *http.Client
	Tester
	TlsConfig() *tls.Config
}

type HttpConfigBuilder interface {
	Build() HttpConfig
	SetCookies(string, []*http.Cookie) HttpConfigBuilder
	SetHttpClient(*http.Client) HttpConfigBuilder
	SetProxy(string) HttpConfigBuilder
	SetT(*testing.T) HttpConfigBuilder
	SetTlsConfig(*tls.Config) HttpConfigBuilder
}
