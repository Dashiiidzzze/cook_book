package api

import "net/http"

// рендеринг страницы профиля
func PageProfile(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/profile" {
		http.NotFound(w, r)
		return
	}

	// Указываем, что возвращаем HTML
	w.Header().Set("Content-Type", "text/html")
	http.ServeFile(w, r, "../frontend/profile.html")
}
