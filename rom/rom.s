.segment "VECTORS"
.word start     ; 0xfffa
.word start     ; 0xfffc
.word start     ; 0xfffe

.segment "CODE"
; ;	.ORG $4000
; *= $4000
start:
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