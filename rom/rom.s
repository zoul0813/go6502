
KBD             = $D010         ;  PIA.A keyboard input
KBDCR           = $D011         ;  PIA.A keyboard control register
DSP             = $D012         ;  PIA.B display output register
DSPCR           = $D013         ;  PIA.B display control register


.segment "CODE"

.export ENTRY
; *= $4000
ENTRY:
	LDA KBD         ; Load character. B7 should be ‘1’.
	LDX #0
	LDY #32

loop:
	LDA KBD         ; Load character. B7 should be ‘1’.
	TYA
	STA $0400,x
	INX
	INY
	CPY #126
	BNE :+
	LDY #32
:
	STA $00
	STA DSP
	JMP loop

UAPUTW:
	RTS

UAGETW:
	RTS

UAGET:
	RTS

