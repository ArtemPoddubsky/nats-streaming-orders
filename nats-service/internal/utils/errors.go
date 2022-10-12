package utils

import (
	"main/internal/log"
	"net/http"
	"os"
)

// PageNotFound writes custom "404 Not Found" error page to ResponseWriter.
func PageNotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)

	page, err := os.ReadFile("./materials/website/404.html")
	if err != nil {
		log.Logger.Errorln(err)
		http.Error(w, "Not Found (404) ", http.StatusNotFound)

		return
	}

	_, err = w.Write(page)
	if err != nil {
		log.Logger.Errorln(err)
		http.Error(w, "Not Found (404) ", http.StatusNotFound)
	}
}

// PageInternalError writes custom "500 Internal Server Error" error page to ResponseWriter.
func PageInternalError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)

	page, err := os.ReadFile("./materials/website/500.html")
	if err != nil {
		log.Logger.Errorln(err)
		http.Error(w, "Temporary Error (500) ", http.StatusInternalServerError)

		return
	}

	_, err = w.Write(page)
	if err != nil {
		log.Logger.Errorln(err)
		http.Error(w, "Temporary Error (500) ", http.StatusInternalServerError)
	}
}
