.segment "VECTORS"
.word start     ; 0xfffa
.word start     ; 0xfffc
.word start     ; 0xfffe

.segment "CODE"
; ;	.ORG $4000
; *= $4000
start:
	LDA #$07
	ADC #$87


	.byte $FF, $F0