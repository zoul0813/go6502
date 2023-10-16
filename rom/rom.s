
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
	STA DSP
	STA $01
	CLC
:
	BCC :-


;;; Hello World from https://www.applefritter.com/comment/60111#comment-60111
  ; 280    00   01   02   03   04   05   06   07   08   09   0A   0B   0C   0D   0E   0F
	; .byte $A2, $0C, $BD, $8B, $02, $20, $EF, $FF, $CA, $D0, $F7, $60, $8D, $C4, $CC, $D2
	; .byte $CF, $D7, $A0, $CF, $CC, $CC, $C5, $C8

UAPUTW:
	RTS

UAGETW:
	RTS

UAGET:
	RTS

