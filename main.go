package main

import (
	"bytes"
	"github.com/go-chi/chi"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
)

func main() {
	r := chi.NewRouter()
	//r.Use(middleware.Logger)
	r.Get("/{service}/*", ServeHTTP)
	r.Post("/{service}/*", ServeHTTP)
	r.Options("/{service}/*", ServeHTTP)
	http.ListenAndServe(":8085", r)
}

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "/favicon.ico" {
		io.WriteString(w, "Request path Error")
		return
	}

	service := chi.URLParam(r, "service")
	pathWithQuery := strings.Replace(r.RequestURI, "/"+service, "",1)

	rawURL := ""
	switch service {
	case "test1":
		rawURL = "http://10.10.10.10" + pathWithQuery
		break
	case "test2":
		rawURL = "http://10.10.10.10" + pathWithQuery
		break
	case "debug":
		rawURL = "http://10.10.10.10" + pathWithQuery
		break
	case "baidu":
		rawURL = "http://www.baidu.com" + pathWithQuery
		break
	}

	remote, err := url.Parse(rawURL)
	if err != nil {
		panic(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.Director = func(req *http.Request) {
		// set the request host header
		req.Host = remote.Host
		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host
		req.URL.Path = remote.Path
		req.URL.RawQuery = remote.RawQuery

		// close gzip
		req.Header.Set("Accept-Encoding", "")
	}
	proxy.ModifyResponse = func(response *http.Response) error {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return  err
		}

		err = response.Body.Close()
		if err != nil {
			return err
		}

		response.Body = ioutil.NopCloser(bytes.NewReader(body))
		response.ContentLength = int64(len(body))
		response.Header.Set("Content-Length", strconv.Itoa(len(body)))

		return nil
	}

	r.URL.Path = ""

	proxy.ServeHTTP(w, r)
}
