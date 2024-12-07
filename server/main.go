package main

import (
	"encoding/json"
	"feh-map-editor/decoder"
	"feh-map-editor/updater"
	"net/http"
)

func corsMiddleware(next http.Handler) http.Handler {
	var corsMiddleware = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Add("Access-Control-Allow-Origin", "*")
		writer.Header().Add("Access-Control-Allow-Methods", "GET")
		next.ServeHTTP(writer, request)
	})

	return corsMiddleware
}

func main() {
	mux := http.NewServeMux()
	go updater.Update()
	var getMapHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		var filename = req.Form.Get("filename")
		var mapData = decoder.Decode(filename)
		var bytes, _ = json.Marshal(mapData)
		w.Write(bytes)
	})
	mux.Handle("GET /map", corsMiddleware(getMapHandler))

	http.ListenAndServe(":3535", mux)
}
