.segment "VECTORS"
.word reset     ; 0xfffa
.word reset     ; 0xfffc
.word reset     ; 0xfffe

.segment "DATA"
byte1:
  .byte $42
byte2:
  .byte $20
word1:
  .word $0420

.segment "CODE"

reset:
  LDA byte1
  INX
  INY
  LDA #$69
  SEC
  LDA byte2
  CLC
  JMP reset