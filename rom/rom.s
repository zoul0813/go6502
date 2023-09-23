.segment "CODE"
; *= $4000
ENTRY:
	LDA #$EA
	LDX #0
	LDY #32

loop:
	TYA
	STA $0400,x
	INX
	INY
	CPY #126
	BNE :+
	LDY #32
:
	STA $00
	JMP loop

UAPUTW:
	RTS

UAGETW:
	RTS

UAGET:
	RTS

.include "wozmon.s"

; ; Interrupt Vectors
.segment "VECTORS"
	.WORD $0F00     ; NMI
	.WORD ENTRY     ; RESET
	.WORD $0000     ; BRK/IRQ