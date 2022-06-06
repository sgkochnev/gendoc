package service

import (
	"io"
	"os"
)

type GendocService struct {
}

func NewGengocService() *GendocService {
	return &GendocService{}
}

func (svc *GendocService) File(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return io.ReadAll(f)
}

func (svc *GendocService) Gendoc(xmlDate []byte, recordID int) (*URLS, error) {
	newURLWord, err := createNewDoc(xmlDate, recordID)
	// newURLPDF, err := createNewPDF(newURLWord)
	resp := &URLS{
		URLWord: newURLWord,
		// URLPDF: newURLPDF,
	}
	return resp, err
}
