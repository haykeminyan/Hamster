package middleware

import (
	"net/http"
	"strings"

	"github.com/matscus/Hamster/Package/JWTToken/jwttoken"
)

//AdminsMiddleware - the admins http middleware func,for check auth and set default headers for admins panel
func AdminsMiddleware(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Max-Age", "600")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000")
		w.Header().Set("X-Frame-Options", "SAMEORIGIN")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		header := r.Header.Get("Authorization")
		splitToken := strings.Split(header, "Bearer ")
		if len(splitToken) != 2 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("{\"Message\":\"Bearer token not in proper format\"}"))
			return
		}
		header = strings.TrimSpace(splitToken[1])
		if header != "" {
			check := jwttoken.IsAdmin(header)
			if check {
				f(w, r)
				return
			}
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("{\"Message\":\"you are not administrator, contact god for access.\"}"))
			return
		}
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("{\"Message\":\"you are not administrator, contact god for access.\"}"))
	}
}
