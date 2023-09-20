.segment "VECTORS"
.word start     ; 0xfffa
.word start     ; 0xfffc
.word start     ; 0xfffe

.segment "CODE"
; ;	.ORG $4000
; *= $4000
start:
	LDX #0

loop:
	TXA
	STA $00
	INX
	JMP loop