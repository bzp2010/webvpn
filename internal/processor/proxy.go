package processor

import (
	"net/http"
	"net/url"

	"github.com/bzp2010/webvpn/internal/model"
)

type ProxyProcessor struct {
	Policy *model.Policy
	RequestURL *url.URL
}

func (p *ProxyProcessor) RequestProcess(req *http.Request) bool {
	policyURL, err := url.Parse(p.Policy.To)

	if err != nil {
		panic(err)
	}

	// request host rewrite
	if p.Policy.HostRewrite != "" {
		req.Host = p.Policy.HostRewrite
	} else {
		req.Host = policyURL.Host
	}

	req.URL.Scheme = policyURL.Scheme
	req.URL.Host = policyURL.Host
	req.URL.Path = p.RequestURL.Path
	req.URL.RawQuery = p.RequestURL.RawQuery

	// close gzip
	req.Header.Set("Accept-Encoding", "")

	return true
}

func (p *ProxyProcessor) ResponseProcess(resp *http.Response) bool {
	return true
}