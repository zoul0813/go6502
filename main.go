package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/zoul0813/go6502/pkg/CPU"
	"github.com/zoul0813/go6502/pkg/Display"
	"github.com/zoul0813/go6502/pkg/IO"
	"github.com/zoul0813/go6502/pkg/Keyboard"
	"github.com/zoul0813/go6502/pkg/Memory"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
)

const (
	// CLOCK_SPEED = time.Nanosecond * 1000 // 1Mhz
	CLOCK_SPEED = time.Nanosecond * 10000 // 1Mhz
	// CLOCK_SPEED         = time.Millisecond * 100 // 1Mhz
	ROM_HEAD            = 0x8000
	ZP_HEAD             = 0x000
	STACK_HEAD          = 0x100
	SCREEN_HEAD  uint16 = 0x400
	scale               = 4
	cols                = 40
	rows                = 25
	fontSize            = 8
	padding             = 32
	frameRate           = 60
	pixelWidth          = (fontSize * cols)
	pixelHeight         = (fontSize * rows)
	screenWidth         = pixelWidth * scale
	screenHeight        = pixelHeight * scale
)

var (
	normalFont  font.Face
	io          *IO.IO
	cpu         *CPU.CPU
	ram         *Memory.Memory
	keyboard    *Keyboard.Keyboard
	display     *Display.Display
	rom         *Memory.Memory
	screenColor = color.RGBA{4, 101, 13, 20}
)

type Game struct {
	// runes   []rune
	// text    string
	runes         []rune
	counter       int
	showRegisters bool
	showZeroPage  bool
	// shader        *ebiten.Shader // Shaders appear to be voodoo magic?
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

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return fmt.Errorf("quit")
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		keyboard.AppendKeys([]rune{'\n'})
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF1) {
		g.showRegisters = !g.showRegisters
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF2) {
		g.showZeroPage = !g.showZeroPage
	}

	g.runes = ebiten.AppendInputChars(g.runes[:0])
	keyboard.AppendKeys(g.runes)

	g.counter++
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()

	t := display.All()

	bound := text.BoundString(normalFont, "W")

	x := 0
	y := 0 + bound.Dy()

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(float64(x), float64(y))
	op.Filter = ebiten.FilterNearest
	op.ColorScale.ScaleWithColor(screenColor)
	text.DrawWithOptions(screen, t, normalFont, op)

	// b := cpu.PC
	// debug := fmt.Sprintf("%04x", b)

	// fmt.Printf("0x00: %02x\n", b)
	// fmt.Printf("%v", t)
	// dx := pixelWidth - bound.Dx()*5
	// dy := pixelHeight - bound.Dy()

	// op := &ebiten.DrawImageOptions{}
	// op.GeoM.Scale(scale, scale)
	// op.GeoM.Translate(float64(dx), float64(dy))
	// op.Filter = ebiten.FilterNearest

	// text.DrawWithOptions(screen, debug, normalFont, op)
	if g.showRegisters {
		cpu.DebugRegister(screen, normalFont, bound, screenHeight, screenWidth)
	}
	if g.showZeroPage {
		DebugZeroPage(screen, normalFont, bound)
	}
}

func DebugZeroPage(screen *ebiten.Image, font font.Face, bound image.Rectangle) {
	s := io.DumpString(0x00, 0xFF)

	dScale := 2.0
	x := float64(bound.Dx())
	y := float64(bound.Dy())
	x = float64(screenWidth) - ((x * dScale) * 50) // width of string
	y = ((y * dScale) * 4)                         // number of lines
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(dScale, dScale)
	op.GeoM.Translate(float64(x), float64(y))
	op.Filter = ebiten.FilterNearest
	clr := color.RGBA{0, 110, 62, 20}
	op.ColorScale.ScaleWithColor(clr)
	text.DrawWithOptions(screen, s, font, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	fmt.Printf("Go 6502... \n")

	// RAM
	ram = Memory.New(0x8000, 0x0000, false)

	// Keyboard
	keyboard = Keyboard.New(0xD010)

	// Display
	display = Display.New(0xD012, cols, rows)

	// ROM
	rom = Memory.New(0x1000, 0xF000, true)
	f, err := os.ReadFile("rom/rom.bin")
	if err != nil {
		fmt.Printf("Can't read file rom/rom.bin\n\n")
		panic(err)
	}

	devices := []*IO.Device{
		IO.NewDevice("RAM", ram, 0x0000),
		IO.NewDevice("Keyboard", keyboard, 0xD010),
		IO.NewDevice("Display", display, 0xD012),
		IO.NewDevice("ROM", rom, 0xF000),
	}
	io = IO.New(devices)

	fmt.Printf("Loading rom %04x (%v) bytes\n", len(f), len(f))
	io.LoadRom(f, 0xF000)
	io.Dump(0xF000, 0xFF)

	// fmt.Printf("Loading ram %04x (%v) bytes\n", len(f), len(f))
	io.LoadRom(f, 0x0000)
	// io.Dump(0x0000, 0xFF)

	cpu = CPU.New(
		0xfffc,     // PC
		0xFF,       // SP
		0x00,       // A
		0xf0,       // X
		0xFE,       // Y
		0b00110000, // Status
		false,      // Single Step
		true,       // DebugMode
	)

	word, _ := io.GetWord(cpu.PC)
	cpu.PC = word

	cpu.Debug()

	// fmt.Printf("ZeroPage: %04x bytes from %04x\n", 0xff, 0x0000)
	// io.Dump(0x0000, 0xff) // Zero Page

	fmt.Printf("\n\n")
	// Reset Vectors
	fmt.Printf("Reset: %04x bytes from %04x\n", 0x0f, 0xfff0)
	io.Dump(0xfff0, 0x0f)

	// ROM Head
	fmt.Printf("ROM Head: %04x bytes from %04x\n", 0xff, 0xF000)
	io.Dump(0xF000, 0xFF)

	// Start of Program?
	fmt.Printf("Program: %04x bytes from %04x\n", 0xff, cpu.PC)
	io.Dump(cpu.PC, 0xff)

	g := &Game{
		// text:    "GO6502\nv0.0.0\n\n% ",
		counter:       0,
		showRegisters: false,
		showZeroPage:  false,
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
			halted, _ := cpu.Step(io)
			if cpu.DebugMode {
				cpu.Debug()
			}
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
