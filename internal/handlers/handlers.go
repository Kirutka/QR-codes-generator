// internal/handlers/handlers.go
package handlers

import (
	"qr/internal/qr" // Импортируем пакет для генерации QR
	"html/template"
	"log"
	"net/http"
)

// PageData структура для передачи данных в шаблон
type PageData struct {
	Error string
}

// HomeHandler обрабатывает запрос к главной странице
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем, что запрашивается корневой путь
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Парсим шаблон из файла
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Printf("Ошибка загрузки шаблона: %v", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	// Отображаем главную страницу с формой
	data := PageData{} // Нет ошибки при первом открытии
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Ошибка выполнения шаблона: %v", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}
}

// GenerateHandler обрабатывает запрос на генерацию QR-кода
func GenerateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// Если метод не POST, перенаправляем на главную
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Получаем данные из формы
	data := r.FormValue("data")
	if data == "" {
		// Если данные пустые, перенаправляем на главную с ошибкой
		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			log.Printf("Ошибка загрузки шаблона: %v", err)
			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, PageData{Error: "Пожалуйста, введите текст или URL."})
		if err != nil {
			log.Printf("Ошибка выполнения шаблона: %v", err)
			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		}
		return
	}

	// Вызываем функцию генерации QR-кода из пакета qr
	img, err := qr.GenerateQR(data)
	if err != nil {
		log.Printf("Ошибка генерации QR-кода: %v", err)
		http.Error(w, "Не удалось создать QR-код. Попробуйте другой ввод.", http.StatusBadRequest)
		return
	}

	// Отправляем изображение клиенту
	qr.SendQRImage(w, img)
}
