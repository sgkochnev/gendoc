package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type docForPatient struct {
	URLTemplate string `json:"URLTemplate"`
	RecordID    int    `json:"RecordID"`
}

func (h *handler) gendoc(w http.ResponseWriter, r *http.Request) {

	doc := &docForPatient{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(doc)
	if err != nil {
		doc.URLTemplate = r.URL.Query().Get("URLTemplate")
		doc.RecordID, err = strconv.Atoi(r.URL.Query().Get("RecordID"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			// log.Println("cannot convert string to inte")
			h.log.Err.Println("cannot convert string to inte")
			return
		}
	}
	if doc.URLTemplate == "" || doc.RecordID <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		h.log.Err.Println("request body missing")
		return
	}

	fileXML, err := XMLFile(doc.URLTemplate)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Err.Println(err)
		return
	}

	urls, err := h.svc.Gendoc(fileXML, doc.RecordID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Err.Println(err)
		return
	}
	urls.URLWord = r.Host + "/" + urls.URLWord

	j, err := json.Marshal(urls)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Err.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)

}

func (h *handler) upload(w http.ResponseWriter, r *http.Request) {
	file, err := h.svc.File(r.URL.Path[1:])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}
	if strings.HasSuffix(r.URL.Path, ".doc") {
		w.Header().Set("Content-Type", "application/msword")
	} else if strings.HasSuffix(r.URL.Path, ".pdf") {
		w.Header().Set("Content-Type", "application/pdf")
	}
	w.Write(file)
}
