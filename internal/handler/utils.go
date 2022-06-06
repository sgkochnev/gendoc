package handler

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

const userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.63 Safari/537.36"

func XMLFile(url string) ([]byte, error) {

	timeout := time.Second * 3
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cl := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)

	req = req.WithContext(ctx)
	resp, err := cl.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s, status code %d", req.URL.String(), resp.StatusCode)
	}

	buf := &bytes.Buffer{}
	io.Copy(buf, resp.Body)
	return buf.Bytes(), nil
}
