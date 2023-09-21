
.segment "CODE"
; ;	.ORG $4000
; *= $4000
ENTRY:
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
	JMP loop

UAPUTW:
	RTS

UAGETW:
	RTS

UAGET:
	RTS


.include "kernal.s"