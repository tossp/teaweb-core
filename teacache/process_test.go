package teacache

import (
	"io/ioutil"
	"net/http"
	"testing"
)

func TestPurgeCache(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8882/webhook", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Tea-Cache-Purge", "1")
	req.Header.Set("Tea-Key", "z8O4MuXixbKH6aiVyZigYTxxovRblR3u")
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(data))
}
