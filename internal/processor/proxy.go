/*
 * Copyright (C) 2021
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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