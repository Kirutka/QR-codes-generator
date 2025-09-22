// internal/qr/generator.go
package qr

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"io"
	"net/http"
	"github.com/skip2/go-qrcode"
)

// GenerateQR создает изображение QR-кода из строки данных.
// Возвращает image.Image или ошибку.
func GenerateQR(data string) (image.Image, error) {
	// Генерируем QR-код в памяти
	qrCode, err := qrcode.New(data, qrcode.Medium)
	if err != nil {
		return nil, err
	}

	// Устанавливаем цвета
	qrCode.BackgroundColor = color.RGBA{255, 255, 255, 255} // Белый фон
	qrCode.ForegroundColor = color.RGBA{51, 51, 51, 255}     // Темно-серый QR-код (#333333)

	// Возвращаем изображение
	return qrCode.Image(300), nil // Размер 300x300 пикселей
}

// SendQRImage кодирует изображение в PNG и отправляет его в http.ResponseWriter.
func SendQRImage(w http.ResponseWriter, img image.Image) {
	// Создаем буфер для хранения изображения PNG
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		// Это внутренняя ошибка, которую сложно обработать на этом этапе
		// Лучше логировать её в вызывающей функции
		http.Error(w, "Ошибка при кодировании изображения.", http.StatusInternalServerError)
		return
	}

	// Устанавливаем заголовок Content-Type для изображения
	w.Header().Set("Content-Type", "image/png")
	// Отправляем изображение напрямую в тело ответа
	_, err = io.Copy(w, &buf)
	if err != nil {
		// Ошибка при отправке, скорее всего клиент отменил запрос
		// Логирование может быть полезно, но не критично для пользователя
		// log.Printf("Ошибка отправки изображения: %v", err)
		return
	}
}
