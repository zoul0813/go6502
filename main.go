package main

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/zoul0813/go6502/pkg/CPU"
	"github.com/zoul0813/go6502/pkg/Memory"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
)

const (
	CLOCK_SPEED         = time.Nanosecond * 1000 // 1Mhz
	ROM_HEAD            = 0x8000
	ZP_HEAD             = 0x000
	STACK_HEAD          = 0x100
	SCREEN_HEAD  uint16 = 0x400
	scale               = 4
	cols                = 40
	rows                = 25
	fontSize            = 8
	padding             = 32
	frameRate           = 30
	pixelWidth          = (fontSize * cols)
	pixelHeight         = (fontSize * rows)
	screenWidth         = pixelWidth * scale
	screenHeight        = pixelHeight * scale
)

var (
	normalFont  font.Face
	cpu         *CPU.CPU
	rom         *Memory.Memory
	screenColor = color.RGBA{75, 220, 125, 20}
)

type Game struct {
	// runes   []rune
	// text    string
	counter int
}

func (g *Game) Update() error {
	// Keyboard input
	// If the enter key is pressed, add a line break.
	// if repeatingKeyPressed(ebiten.KeyEnter) || repeatingKeyPressed(ebiten.KeyNumpadEnter) {
	// 	text += "\n"
	// }

	// If the backspace key is pressed, remove one character.
	// if repeatingKeyPressed(ebiten.KeyBackspace) {
	// 	if len(text) >= 1 {
	// 		text = text[:len(g.text)-1]
	// 	}
	// }

	g.counter++
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()
	// Blink the cursor.

	t := ""
	// if g.counter%60 < 30 {
	// 	t += "_"
	// }

	for i := 0; i < 256; i++ {
		if i%cols == 0 && len(t) > 0 {
			t += "\n"
		}
		c, _ := rom.Get(SCREEN_HEAD + uint16(i))
		if c >= 32 && c <= 126 {
			t += string(c)
		} else {
			t += " "
		}
	}

	bound := text.BoundString(normalFont, "W")

	x := 0
	y := 0 + bound.Dy()

	text.Draw(screen, t, normalFont, x, y, screenColor)

	b, _ := rom.Get(0x00)
	debug := fmt.Sprintf("%02x", b)

	// fmt.Printf("0x00: %02x\n", b)
	// fmt.Printf("%v", t)
	text.Draw(screen, debug, normalFont, pixelWidth-bound.Dx(), pixelHeight-bound.Dy(), screenColor)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth / scale, screenHeight / scale
}

func main() {
	fmt.Printf("Go 6502... \n")

	// ram := Memory.New(0x8000, 0x0000, false)
	rom = Memory.New(0xffff, 0x0000, false)
	f, err := os.ReadFile("rom/rom.bin")
	if err != nil {
		fmt.Printf("Can't read file rom/rom.bin\n\n")
		panic(err)
	}

	rom.Load(f)

	// rom.Dump(0x0000, 0xff) // Zero Page
	rom.Dump(0xfff0, 0x0f) // Reset Vectors
	rom.Dump(0x8000, 0xff) // Start of Program?

	// should just loop infinitely now ...

	cpu = CPU.New(
		0xfffc,     // PC
		0xFF,       // SP
		0x00,       // A
		0xf0,       // X
		0xFE,       // Y
		0b00110000, // Status
		false,      // Single Step
		false,      // DebugMode
	)

	word, _ := rom.GetWord(cpu.PC)
	cpu.PC = word

	cpu.Debug()

	// for {
	// 	if cpu.SingleStep {
	// 		DebugConsole(cpu, rom)
	// 	}

	// 	halted, _ := cpu.Step(rom)

	// 	if halted {
	// 		DebugConsole(cpu, rom)
	// 	}

	// 	fmt.Print("\n") // always end the instructions debug lines
	// 	cpu.Debug()
	// }

	g := &Game{
		// text:    "GO6502\nv0.0.0\n\n% ",
		counter: 0,
	}

	fontFile, err := os.ReadFile("assets/fonts/C64_Pro_Mono-STYLE.ttf")
	if err != nil {
		log.Fatal(err)
	}

	fontFace, err := sfnt.Parse(fontFile)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72

	normalFont, err = opentype.NewFace(fontFace, &opentype.FaceOptions{
		Size:    8,
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		log.Fatal(err)
	}

	ebiten.SetWindowSize(screenWidth+padding, screenHeight+padding)
	ebiten.SetWindowTitle("TypeWriter (Ebitengine Demo)")
	ebiten.SetTPS(frameRate)

	doQuit := false
	cpuClock := time.NewTicker(CLOCK_SPEED)
	defer cpuClock.Stop()

	go func() {
		for {
			if doQuit {
				return
			}
			<-cpuClock.C
			halted, _ := cpu.Step(rom)
			if halted {
				fmt.Printf("Halted: %v", cpu)
			}
		}
	}()

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
	doQuit = true
}
