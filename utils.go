package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/hagen1778/chproxy/log"
)

func doQuery(query, addr string) error {
	resp, err := http.Post(fmt.Sprintf("%s/?query=%s", addr, url.QueryEscape(query)), "", nil)
	if err != nil {
		return fmt.Errorf("error in clickhouse query %q at %q: %s", query, addr, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		responseBody, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code returned from query %q at %q: %d. Expecting %d. Response body: %q",
			query, addr, resp.StatusCode, http.StatusOK, responseBody)
	}
	return nil
}

func respondWIthErr(rw http.ResponseWriter, err error) {
	log.Errorf("proxy failed: %s", err)
	rw.WriteHeader(http.StatusInternalServerError)
	rw.Write([]byte(err.Error()))
}

func extractUserFromRequest(req *http.Request) string {
	if name, _, ok := req.BasicAuth(); ok {
		return name
	}

	if name := req.Form.Get("user"); name != "" {
		return name
	}

	return "default"
}