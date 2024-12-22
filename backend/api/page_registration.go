package api

import "net/http"

func PageRegistration(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/registration" {
		http.NotFound(w, r)
		return
	}

	// Указываем, что возвращаем HTML
	w.Header().Set("Content-Type", "text/html")
	http.ServeFile(w, r, "../frontend/register.html")
}
