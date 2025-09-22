package main

import (
	"bytes"
	"fmt"
	"image/color"
	"image/png"
	"io"
	"log"
	"net/http"
	"text/template"

	"github.com/skip2/go-qrcode"
)

// PageData структура для передачи данных в шаблон
type PageData struct {
	Error string
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
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

func generateHandler(w http.ResponseWriter, r *http.Request) {
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

	// Генерируем QR-код в памяти
	qr, err := qrcode.New(data, qrcode.Medium)
	if err != nil {
		log.Printf("Ошибка создания QR-кода: %v", err)
		http.Error(w, "Не удалось создать QR-код. Попробуйте другой ввод.", http.StatusBadRequest)
		return
	}

	// Устанавливаем цвета (по желанию, можно сделать настраиваемыми)
	qr.BackgroundColor = color.RGBA{255, 255, 255, 255} // Белый фон
	qr.ForegroundColor = color.RGBA{51, 51, 51, 255}     // Темно-серый QR-код (#333333)

	// Создаем буфер для хранения изображения PNG
	var buf bytes.Buffer
	img := qr.Image(300) // Размер 300x300 пикселей
	err = png.Encode(&buf, img)
	if err != nil {
		log.Printf("Ошибка кодирования изображения: %v", err)
		http.Error(w, "Ошибка при создании изображения.", http.StatusInternalServerError)
		return
	}

	// Устанавливаем заголовок Content-Type для изображения
	w.Header().Set("Content-Type", "image/png")
	// Отправляем изображение напрямую в тело ответа
	_, err = io.Copy(w, &buf)
	if err != nil {
		log.Printf("Ошибка отправки изображения: %v", err)
		// Примечание: если заголовки уже отправлены, http.Error не сработает корректно
	}
}

func main() {
	// Обрабатываем статические файлы (CSS)
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Обрабатываем маршруты
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/generate", generateHandler)

	fmt.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
