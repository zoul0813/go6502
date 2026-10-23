package main

import (
	_ "embed"
	"flag"
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"strings"
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
	ROM_HEAD            = 0x8000
	ZP_HEAD             = 0x000
	STACK_HEAD          = 0x100
	SCREEN_HEAD  uint16 = 0x400
	scale               = 4
	cols                = 40
	rows                = 24
	fontSize            = 8
	padding             = 32
	frameRate           = 60
	pixelWidth          = (fontSize * cols)
	pixelHeight         = (fontSize * rows)
	screenWidth         = pixelWidth * scale
	screenHeight        = pixelHeight * scale
)

var (
	normalFont font.Face
	io         *IO.IO
	clockSpeed = time.Nanosecond * 1000 // 1Mhz
	// clockSpeed = time.Nanosecond * 10000 // 100Hz?
	// TODO: make the clockSpeed variable with an argunment
	// clockSpeed  = time.Millisecond * 100 // 10Hz
	cpu         *CPU.CPU
	ram         *Memory.Memory
	keyboard    *Keyboard.Keyboard
	display     *Display.Display
	rom         *Memory.Memory
	screenColor = color.RGBA{4, 101, 13, 255}
)

//go:embed assets/fonts/C64_Pro_Mono-STYLE.ttf
var fontFile []byte

type Game struct {
	// runes   []rune
	// text    string
	// runes         []rune
	keys          []ebiten.Key
	counter       int
	showRegisters bool
	showZeroPage  bool
	showWozIn     bool
	showStack     bool
	elapsed       time.Duration
	// shader        *ebiten.Shader // Shaders appear to be voodoo magic?
}

func (g *Game) Update() error {
	// Keyboard input
	g.keys = inpututil.AppendJustPressedKeys(g.keys[:0])

	for _, key := range g.keys {
		switch key {
		case ebiten.KeyF1:
			g.showRegisters = !g.showRegisters
		case ebiten.KeyF2:
			g.showZeroPage = !g.showZeroPage
		case ebiten.KeyF3:
			g.showWozIn = !g.showWozIn
		case ebiten.KeyF4:
			g.showStack = !g.showStack
		case ebiten.KeyF5:
			s := io.DumpString(0x0000, 0xFFFF)
			os.WriteFile("dump.txt", []byte(s), 0644)
			return fmt.Errorf("quit")
		case ebiten.KeyF7:
			cpu.SingleStep = !cpu.SingleStep
		case ebiten.KeyF8:
			if !cpu.SingleStep {
				continue
			}
			halted, err := cpu.Step(io)
			if err != nil {
				return err
			}
			if cpu.DebugMode {
				cpu.Debug()
			}
			if halted {
				fmt.Printf("Halted: %v", cpu)
			}
		case ebiten.KeyHome:
			// reset
		case ebiten.KeyEscape:
			keyboard.AppendKey(0x1B) // ESC 27
		case ebiten.KeyEnter:
			keyboard.AppendKey(0x0D) // LF 10
		default:
			var buffer []rune
			buffer = ebiten.AppendInputChars(buffer[:0])
			for _, r := range buffer {
				b := strings.ToUpper(string(byte(r)))[0]
				fmt.Printf("Key: %02x %02x '%v' '%v'\n", r, b, b, string(b))
				keyboard.AppendKey(b)
			}
			// fmt.Printf("KeyCode: %v, %v\n", key, buffer)
			// name := ebiten.KeyName(key)
			// fmt.Printf("KeyName: %v\n", name)
			// if len(name) > 0 {
			// 	b := name[0]
			// 	fmt.Printf("Key: %02x %08b '%v'\n", b, b, string(b))
			// 	keyboard.AppendKey(b)
			// }
		}
	}

	g.counter++
	return g.advance(time.Second / frameRate)
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()

	blink := false
	if g.counter%frameRate < (frameRate / 2) {
		blink = true
	}
	t := display.All(blink)

	bound := text.BoundString(normalFont, "W")

	x := 0
	y := 0 + bound.Dy()*scale

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(float64(x), float64(y))
	op.Filter = ebiten.FilterNearest
	op.ColorScale.ScaleWithColor(screenColor)
	text.DrawWithOptions(screen, t, normalFont, op)

	if g.showRegisters {
		cpu.DebugRegister(screen, normalFont, bound, screenHeight, screenWidth)
	}

	if g.showZeroPage {
		DebugMemory(0x00, 0xFF, screen, normalFont, bound)
	}

	if g.showStack {
		DebugMemory(0x0100, 0xFF, screen, normalFont, bound)
	}

	if g.showWozIn {
		DebugMemory(0x0200, 0xFF, screen, normalFont, bound)
	}
}

func DebugMemory(start uint16, size uint16, screen *ebiten.Image, font font.Face, bound image.Rectangle) {
	s := io.DumpString(start, size)

	dScale := 2.0
	x := float64(bound.Dx())
	y := float64(bound.Dy())
	x = float64(screenWidth) - ((x * dScale) * 55) // width of string
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

// Ebitengine calls Update and Draw serially. Keep the CPU, keyboard, display,
// and debugger on that same owner instead of sharing mutable state with a ticker.
// This is an instruction-rate budget, not cycle-accurate 6502 timing.
func (g *Game) advance(elapsed time.Duration) error {
	if clockSpeed <= 0 {
		return fmt.Errorf("clock interval must be positive")
	}
	if cpu.SingleStep {
		g.elapsed = 0
		return nil
	}
	g.elapsed += elapsed
	for g.elapsed >= clockSpeed {
		g.elapsed -= clockSpeed
		if err := automaticStep(); err != nil {
			return err
		}
		if cpu.SingleStep {
			g.elapsed = 0
			break
		}
	}
	return nil
}

// automaticStep is also used after resuming from single-step mode.
func automaticStep() error {
	if cpu.SingleStep {
		return nil
	}
	halted, err := cpu.Step(io)
	cpu.Debug()
	if halted || err != nil {
		cpu.SingleStep = true
	}
	return err
}

func main() {
	fmt.Printf("Go 6502... \n")

	// set clockSpeed to value of CLI argument
	clockMultiplier := 1000 // 1000Khz
	flag.IntVar(&clockMultiplier, "clock", 1000, "Clock Speed (kHz)")
	singleStep := false
	debugMode := false
	hz := false
	flag.BoolVar(&singleStep, "single", false, "Single Step")
	flag.BoolVar(&debugMode, "debug", false, "Debug Mode")
	flag.BoolVar(&hz, "hz", false, "Set Clock to Hz")
	flag.Parse()

	// calculate the clock speed using kHz
	khz := time.Microsecond * 1_000
	if hz {
		khz = time.Microsecond * 1_000_000
	}
	d := khz / time.Duration(clockMultiplier)
	clockSpeed = time.Nanosecond * d
	fmt.Printf("Clock Mult : %v (%v)\n", clockMultiplier, d)
	fmt.Printf("Clock Speed: %v\n", clockSpeed)

	if err := initializeMachine(singleStep, debugMode); err != nil {
		log.Fatal(err)
	}

	g := &Game{
		// text:    "GO6502\nv0.0.0\n\n% ",
		counter:       0,
		showRegisters: false,
		showZeroPage:  false,
		showWozIn:     false,
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
	ebiten.SetWindowTitle("Gosho-1 (Apple 1 Emulator in Go)")
	ebiten.SetTPS(frameRate)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
