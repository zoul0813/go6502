.segment "VECTORS"
.word reset     ; 0xfffa
.word reset     ; 0xfffc
.word reset     ; 0xfffe

.segment "DATA"
byte1:
  .byte $42, $43, $44, $45, $46, $47, $48, $49, $50, $51
byte2:
  .byte $20
word1:
  .word $0420
  .byte "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed pulvinar quam justo. Phasellus nec magna vulputate, lobortis justo non, rhoncus tortor. In dictum ac neque non mattis. Donec facilisis massa eu dolor aliquet, a sodales nulla condimentum. Proin viverra pretium euismod. Duis pretium sodales lacus ut pretium."

.segment "CODE"

reset:
:
  LDA #$0F
:
  ADC #$01
  BEQ :--
  ADC #$01
  ADC #$01
  ADC #$01
  JMP :-