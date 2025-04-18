package main

import (
	"html/template"
	"net/http"
	"sync"
	"time"
)

type Check struct {
	URL       string
	Status    string
	Error     string
	CheckedAt string
}

type PageData struct {
	History []Check
}

var (
	checkHistory []Check
	mu           sync.Mutex
)

func main() {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			url := r.FormValue("url")
			check := Check{
				URL:       url,
				CheckedAt: time.Now().Format("2006-01-02 15:04:05"),
			}

			resp, err := http.Get(url)
			if err != nil {
				check.Error = err.Error()
			} else {
				check.Status = resp.Status
				resp.Body.Close()
			}

			mu.Lock()
			checkHistory = append([]Check{check}, checkHistory...) // newest on top
			if len(checkHistory) > 10 {
				checkHistory = checkHistory[:10] // limit to 10 entries
			}
			mu.Unlock()
		}

		mu.Lock()
		data := PageData{History: checkHistory}
		mu.Unlock()

		tmpl.Execute(w, data)
	})

	http.ListenAndServe(":8080", nil)
}
