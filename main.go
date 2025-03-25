package main

import (
	"log"
	"time"

	"github.com/gdamore/tcell/v2"
)

const (
	screenWidth  = 40
	screenHeight = 20
	paddleHeight = 4
)

var (
	ballX, ballY   = screenWidth / 2, screenHeight / 2
	ballDX, ballDY = 1, 1
	paddleLeftY    = screenHeight/2 - paddleHeight/2
	paddleRightY   = screenHeight/2 - paddleHeight/2
	scoreLeft      = 0
	scoreRight     = 0
)

func drawScreen(s tcell.Screen) {
	s.Clear()

	// Desenha a bola
	s.SetContent(ballX, ballY, '●', nil, tcell.StyleDefault.Foreground(tcell.ColorRed))

	// Desenha as raquetes
	for i := 0; i < paddleHeight; i++ {
		s.SetContent(1, paddleLeftY+i, '█', nil, tcell.StyleDefault.Foreground(tcell.ColorGreen))
		s.SetContent(screenWidth-2, paddleRightY+i, '█', nil, tcell.StyleDefault.Foreground(tcell.ColorBlue))
	}

	// Atualiza a tela
	s.Show()
}

func moveBall() {
	ballX += ballDX
	ballY += ballDY

	// Rebote nas bordas superior e inferior
	if ballY <= 0 || ballY >= screenHeight-1 {
		ballDY *= -1
	}

	// Rebote na raquete esquerda
	if ballX == 2 && ballY >= paddleLeftY && ballY < paddleLeftY+paddleHeight {
		ballDX *= -1
	}

	// Rebote na raquete direita
	if ballX == screenWidth-3 && ballY >= paddleRightY && ballY < paddleRightY+paddleHeight {
		ballDX *= -1
	}

	// Pontuação
	if ballX <= 0 {
		scoreRight++
		resetBall()
	} else if ballX >= screenWidth-1 {
		scoreLeft++
		resetBall()
	}
}

func resetBall() {
	ballX, ballY = screenWidth/2, screenHeight/2
	ballDX, ballDY = 1, 2
}

func main() {
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatal(err)
	}
	if err := screen.Init(); err != nil {
		log.Fatal(err)
	}
	defer screen.Fini()

	quit := make(chan struct{})

	// Captura entrada do usuário
	go func() {
		for {
			ev := screen.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyEscape {
					close(quit)
					return
				} else if ev.Key() == tcell.KeyUp {
					if paddleRightY > 0 {
						paddleRightY--
					}
				} else if ev.Key() == tcell.KeyDown {
					if paddleRightY < screenHeight-paddleHeight {
						paddleRightY++
					}
				} else if ev.Rune() == 'w' {
					if paddleLeftY > 0 {
						paddleLeftY--
					}
				} else if ev.Rune() == 's' {
					if paddleLeftY < screenHeight-paddleHeight {
						paddleLeftY++
					}
				}
			}
		}
	}()

	// Loop principal do jogo
	for {
		select {
		case <-quit:
			return
		default:
			moveBall()
			drawScreen(screen)
			time.Sleep(50 * time.Millisecond)
		}
	}
}
