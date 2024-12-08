package api

import "net/http"

func PageCreate(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/create" {
		http.NotFound(w, r)
		return
	}

	// Указываем, что возвращаем HTML
	w.Header().Set("Content-Type", "text/html")
	http.ServeFile(w, r, "../frontend/create.html") // Путь к вашему HTML-файлу
}
