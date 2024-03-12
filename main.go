package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func stringToSlice(str string) []string {
	slice := make([]string, 0, len(str))
	for _, char := range str {
		slice = append(slice, string(char))
	}
	return slice
}

func getClientIP(ipAddress string) string {
	if strings.HasPrefix(ipAddress, "[") && strings.Contains(ipAddress, "]:") {
		colonIndex := strings.LastIndex(ipAddress, "]:")
		if colonIndex != -1 {
			ipAddress = ipAddress[1:colonIndex]
		}
	} else {
		if colonIndex := strings.LastIndex(ipAddress, ":"); colonIndex != -1 {
			ipAddress = ipAddress[:colonIndex]
		}
	}

	return ipAddress
}

func getHeaders(h http.Header) map[string]interface{} {
	var headers map[string]interface{}
	headers = make(map[string]interface{})
	for key, values := range h {
		for _, value := range values {
			headers[key] = value
		}
	}

	return headers
}

func echoResponse(w http.ResponseWriter, r *http.Request) {
	var body map[string]interface{}
	body = make(map[string]interface{})
	bodyContent, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	body["base64_encoded"] = bodyContent
	body["base64_decoded"] = string(bodyContent)

	var echo map[string]interface{}
	echo = make(map[string]interface{})

	echo["host"] = r.Host
	echo["protocol"] = r.Proto
	echo["path"] = r.URL.Path
	echo["ip"] = getClientIP(r.RemoteAddr)
	echo["method"] = r.Method
	echo["uri"] = r.RequestURI
	echo["headers"] = r.Header
	echo["params"] = r.URL.Query()
	echo["body"] = body

	w.Header().Set("Content-Type", "application/json")
	echoJson, err := json.MarshalIndent(echo, "", "    ")
	if err != nil {
		http.Error(w, "Failed to marshal headers to JSON", http.StatusInternalServerError)
		return
	}

	w.Write(echoJson)
}

func main() {
	http.HandleFunc("/", echoResponse)
	fmt.Println("Server listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
