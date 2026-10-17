package httpapi

import "net/http"

func cors(origins []string) func(http.Handler) http.Handler {
	allowed := make(map[string]struct{}, len(origins))
	for _, origin := range origins {
		allowed[origin] = struct{}{}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			_, ok := allowed[origin]
			if origin != "" {
				w.Header().Add("Vary", "Origin")
			}
			if r.Method == http.MethodOptions {
				if origin == "" {
					writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed.")
					return
				}
				if !ok {
					writeError(w, http.StatusForbidden, "forbidden", "Origin is not allowed.")
					return
				}
				allow(w, origin)
				w.WriteHeader(http.StatusNoContent)
				return
			}
			if ok {
				allow(w, origin)
			}
			next.ServeHTTP(w, r)
		})
	}
}

func allow(w http.ResponseWriter, origin string) {
	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type")
	w.Header().Set("Access-Control-Max-Age", "600")
}
