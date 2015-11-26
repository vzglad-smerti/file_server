package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type ConfigJson struct {
	Port        string `json:"port"`
	Dir         string `json:"dir"`
	SecondCache string `json:"second_cache"`
}

func maxAgeHandler(seconds int64, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", fmt.Sprintf("max-age=%d, public, must-revalidate, proxy-revalidate", seconds))
		h.ServeHTTP(w, r)
	})
}

func main() {
	jsonSrc, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}
	var config ConfigJson
	json.Unmarshal(jsonSrc, &config)

	i, err := strconv.ParseInt(config.SecondCache, 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", maxAgeHandler(i, http.FileServer(http.Dir(config.Dir))))
	log.Fatal(http.ListenAndServe(":"+config.Port, nil))
}
