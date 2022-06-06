package service

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Data struct {
	Result            int    `josn:"result"`
	Resultdescription string `json:"resultdescription"`
	Resultdata        string `json:"resultdata"`
}

const urlAddr = "https://sycret.ru/service/apigendoc/apigendoc"

const userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.63 Safari/537.36"

func PatientData(textField []byte, recordID int) ([]byte, error) {

	timeout := time.Second * 1
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cl := http.Client{}

	b := fmt.Sprintf("{\"text\":\"%s\",\"recordid\":%d}", textField, recordID)
	reqBoby := strings.NewReader(b)

	req, err := http.NewRequest("POST", urlAddr, reqBoby)
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

	jsonResponse := &Data{}

	jsonDec := json.NewDecoder(resp.Body)
	err = jsonDec.Decode(jsonResponse)

	if jsonResponse.Resultdescription != "OK" {
		return nil, err
	}

	return []byte(jsonResponse.Resultdata), nil
}

func createNewDoc(xmlData []byte, recordID int) (string, error) {

	rd := bytes.NewReader(xmlData)

	f, err := createFile()
	if err != nil {
		return "", err
	}
	defer f.Close()

	err = replaceText(f, rd, recordID)
	if err != nil {
		f.Close()
		os.Remove(f.Name())
		return "", err
	}

	return f.Name(), nil
}

func createFile() (*os.File, error) {
	const uploadDir = "upload/word/"

	dir := uploadDir + time.Now().Format("2006-01")

	err := os.MkdirAll(dir, os.ModeDir)
	if !errors.Is(err, os.ErrExist) && err != nil {
		return nil, err
	}

	fileName := time.Now().Format("2006-01-02 15-04-05") + ".doc"
	filepath := fmt.Sprintf("%s/%s", dir, fileName)

	return os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY, os.ModeDevice)
}

func replaceText(w io.Writer, r io.Reader, recordID int) error {

	const tagText = "<ns1:text field="
	const tagTextEnd = "</ns1:text"

	bufRd := bufio.NewReader(r)
	var field []byte

	for line, err := bufRd.ReadBytes('>'); err != io.EOF; line, err = bufRd.ReadBytes('>') {

		if err != nil {
			return err
		}

		if bytes.Contains(line, []byte(tagText)) {
			newLine := bytes.Trim(line, " ")
			startIndex := bytes.Index(newLine, []byte("\""))
			finishIndex := bytes.LastIndex(newLine, []byte("\""))
			field = newLine[startIndex+1 : finishIndex]

		} else if bytes.Contains(line, []byte(tagTextEnd)) {
			field = []byte("")
		}

		if len(field) != 0 {
			line = tryReplace(line, field, recordID)
		}

		w.Write(line)

	}

	return nil
}

func tryReplace(line, field []byte, recordID int) []byte {
	const tagWTEnd = "</w:t>"

	if bytes.Contains(line, []byte(tagWTEnd)) {
		text, err := PatientData(field, recordID)
		if err != nil {
			log.Println(err)
			text = []byte("")
		}
		str := fmt.Sprintf("%s</w:t>", text)
		return []byte(str)
	}
	return line
}
