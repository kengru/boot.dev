module github.com/kengru/boot.dev/wserver

go 1.22.0

require (
	github.com/go-chi/chi/v5 v5.0.11
	golang.org/x/crypto v0.19.0
	internal/database v1.0.0
)

require (
	github.com/golang-jwt/jwt/v5 v5.2.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
)

replace internal/database => ./internal/database
