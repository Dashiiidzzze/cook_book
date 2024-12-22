package api

import "net/http"

func PageLogin(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/auth" {
		http.NotFound(w, r)
		return
	}

	// Указываем, что возвращаем HTML
	w.Header().Set("Content-Type", "text/html")
	http.ServeFile(w, r, "../frontend/auth.html")
}
