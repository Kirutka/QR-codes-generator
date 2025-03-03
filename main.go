package main

import (
	"image/color"
	"log"

	"github.com/skip2/go-qrcode"
)

func main() {
	data := "http://192.168.31.50:8080"

	qr, err := qrcode.New(data, qrcode.Medium)
	if err != nil {
		log.Fatal(err)
	}

	qr.BackgroundColor = color.RGBA{255, 255, 0, 255} // Жёлтый фон
	qr.ForegroundColor = color.RGBA{0, 0, 255, 255}   // Синий QR-код

	err = qr.WriteFile(256, "qrcode.png") // задаём размер кода
	if err != nil {
		log.Fatal(err)
	}
}
