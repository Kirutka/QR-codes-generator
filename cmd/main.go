// main.go
package main

import (
	"fmt"
	"log"
	"net/http"

	"qr/internal/handlers" // Импортируем наш внутренний пакет
)

func main() {
	// Обрабатываем статические файлы (CSS)
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Обрабатываем маршруты с помощью обработчиков из internal/handlers
	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/generate", handlers.GenerateHandler)

	fmt.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
