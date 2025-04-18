package main

import (
	"encoding/json"
	"feh-map-editor/decoder"
	"feh-map-editor/encoder"
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

	var postMapHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var mapData encoder.EncodePayload
		err := json.NewDecoder(req.Body).Decode(&mapData)
		if err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		err = encoder.Encode(mapData, "dump.bin")
		if err != nil {
			http.Error(w, "Failed to write file", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("File written successfully"))
	})
	mux.Handle("GET /map", corsMiddleware(getMapHandler))
	mux.Handle("POST /map", corsMiddleware(postMapHandler))

	http.ListenAndServe(":3535", mux)
}
