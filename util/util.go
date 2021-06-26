package util

import (
	"bufio"
	"encoding/json"
	"net/http"
	"os"
)

func RespondWithInfo(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"info": message})
}
func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"error": message})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func ReadAllLines(filePath string) []string {
	file, _ := os.Open(filePath)
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	file.Close()
	return lines
}
