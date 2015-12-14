package httpreq

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type HTTPService interface {
	HTTPRequest(string, string, io.Reader, map[string]string, bool) (*HTTPResponse, error)
}

type HTTPResponse struct {
	Body      string
	BodyBytes []byte
	Resp      *http.Response
}

func getRespBody(resp *http.Response) (string, []byte, error) {
	defer resp.Body.Close()
	bb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", []byte{}, err
	}
	return string(bb), bb, nil
}

// HTTPRequest executes a given HTTP API request, returning response body
func HTTPRequest(url string, method string, body io.Reader, headers map[string]string, failOnError bool) (*HTTPResponse, error) {
	hresp := &HTTPResponse{}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return hresp, err
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	hc := http.Client{}
	resp, err := hc.Do(req)
	if err != nil {
		return hresp, err
	}
	bs, bb, err := getRespBody(resp)
	if err != nil {
		return hresp, err
	}
	hresp.Body = bs
	hresp.BodyBytes = bb
	hresp.Resp = resp
	if resp.StatusCode > 399 && failOnError {
		return hresp, fmt.Errorf("Server response indicates failure: %v %v", resp.StatusCode, bs)
	}
	return hresp, nil
}
