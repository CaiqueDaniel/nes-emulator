package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math/rand"
	"nes-emu/src/emulator"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
)

func main() {
	driver.Main(func(s screen.Screen) {
		width, height := 640, 480
		size := image.Point{X: width, Y: height}

		// 1. Cria a janela
		w, err := s.NewWindow(&screen.NewWindowOptions{
			Title:  "Exemplo de Buffer de Pixels - Shiny",
			Width:  width,
			Height: height,
		})
		if err != nil {
			log.Fatalf("Erro ao criar janela: %v", err)
		}
		defer w.Release()

		// 2. Aloca um buffer de pixels na memória RAM
		buf, err := s.NewBuffer(size)
		if err != nil {
			log.Fatalf("Erro ao alocar buffer: %v", err)
		}
		defer buf.Release()

		emulator := emulator.NewEmulatorModule(&w, &buf)

		// 2. Goroutine para injetar eventos de renderização na fila de eventos da janela
		go func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("Recovered in goroutine:", r)
				}
			}()

			err := emulator.GetStartGameController().StartGame("./test/resources/Zelda.NES")

			if err != nil {
				panic(err)
			}
		}()

		// 4. Laço principal de eventos (Event Loop manual)
		for {
			event := w.NextEvent()

			switch e := event.(type) {

			case lifecycle.Event:
				// Fecha a aplicação quando a janela é encerrada
				if e.To == lifecycle.StageDead {
					return
				}

			case paint.Event:
				// Transfere o buffer para a janela e exibe na tela
				w.Upload(image.Point{0, 0}, buf, buf.Bounds())
				w.Publish()
			}
		}
	})
}

func fillPixelBuffer(img *image.RGBA) {
	bounds := img.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Cria um gradiente simples baseado nas coordenadas (X, Y)
			c := color.RGBA{
				R: uint8(rand.Int() % 256),
				G: uint8(y % 256),
				B: 200,
				A: 255,
			}
			img.SetRGBA(x, y, c)
		}
	}
}
