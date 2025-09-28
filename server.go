package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"
)

func main() {
	// ---- статика: index.html, mp3, json ----
	http.Handle("/", http.FileServer(http.Dir(".")))

	// ---- простые API для data/done ----
	http.HandleFunc("/api/data", handleJSONFile("data.json"))
	http.HandleFunc("/api/done", handleJSONFile("done.json"))

	addr := ":8080"
	url := "http://localhost" + addr

	go func() {
		time.Sleep(300 * time.Millisecond)
		openBrowser(url)
	}()

	log.Printf("Serving %s (CTRL+C to stop)\n", url)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func handleJSONFile(filename string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// CORS на случай, если откроешь с другого origin
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		switch r.Method {
		case http.MethodGet:
			// если файла нет — отдадим пустую структуру
			if _, err := os.Stat(filename); os.IsNotExist(err) {
				def := "[]"
				if filename == "done.json" {
					def = "{}"
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(def))
				return
			}
			http.ServeFile(w, r, filename)
		case http.MethodPut, http.MethodPost:
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if !json.Valid(body) {
				http.Error(w, "invalid JSON", http.StatusBadRequest)
				return
			}
			if err := os.WriteFile(filename, body, 0644); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusNoContent)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url) // macOS
	case "windows":
		// запускает через shell ассоциированный браузер
		cmd = exec.Command("cmd", "/c", "start", "", url)
	default:
		// Linux/BSD
		cmd = exec.Command("xdg-open", url)
	}
	if err := cmd.Start(); err != nil {
		fmt.Println("Open browser manually:", url, "err:", err)
	}
}
