package ytp

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// Hop-by-hop headers. These are removed when sent to the backend.
// http://www.w3.org/Protocols/rfc2616/rfc2616-sec13.html
var hopHeaders = []string{
	"Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te", // canonicalized version of "TE"
	"Trailers",
	"Transfer-Encoding",
	"Upgrade",
}

func dropHopHeaders(header http.Header) {
	for _, h := range hopHeaders {
		header.Del(h)
	}
}

func copyHeader(src, dst http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func setXForwardHeader(header http.Header, host string) {
	if hosts, ok := header["X-Forwarded-For"]; ok {
		host = strings.Join(hosts, ", ") + ", " + host
	}
	header.Set("X-Forwarded-For", host)
}

type Proxy struct {
	host string
	authToken string
	client *http.Client
}

func New(host, authToken string) (*Proxy, error) {
	return &Proxy{
		host: host,
		authToken: authToken,

		client: &http.Client{
			Transport: &http.Transport{
				MaxIdleConns: 10,
				IdleConnTimeout: 30 * time.Second,
			},
		},
	}, nil
}

func (p *Proxy) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	req.URL.Scheme = "https"
	req.URL.Host = p.host
	req.Host = ""
	req.RequestURI = ""
	dropHopHeaders(req.Header)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", p.authToken))
	log.Println(req.URL.String())

	if resp, err := p.client.Do(req); err != nil {
		http.Error(wr, "Server Error", http.StatusInternalServerError)
		log.Fatal("Proxy.ServeHTTP:", err)
	} else {
		defer resp.Body.Close()
		dropHopHeaders(resp.Header)
		copyHeader(resp.Header, wr.Header())
		wr.WriteHeader(resp.StatusCode)
		io.Copy(wr, resp.Body)
	}
}
