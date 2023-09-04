.segment "VECTORS"
.word start     ; 0xfffa
.word start     ; 0xfffc
.word start     ; 0xfffe

; .segment "DATA"
; byte1:
;   .byte $42, $43, $44, $45, $46, $47, $48, $49, $50, $51
; byte2:
;   .byte $20
; word1:
;   .word $0420
;   .byte "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed pulvinar quam justo. Phasellus nec magna vulputate, lobortis justo non, rhoncus tortor. In dictum ac neque non mattis. Donec facilisis massa eu dolor aliquet, a sodales nulla condimentum. Proin viverra pretium euismod. Duis pretium sodales lacus ut pretium."

.segment "CODE"

start:
; EXPECTED FINAL RESULTS: $0210 = FF
; (any other number will be the
;  test that failed)

; initialize:
	LDA #$00
	STA $0210
	; store each test's expected
	LDA #$55
	STA $0200
	LDA #$AA
	STA $0201
	LDA #$FF
	STA $0202
	LDA #$6E
	STA $0203
	LDA #$42
	STA $0204
	LDA #$33
	STA $0205
	LDA #$9D
	STA $0206
	LDA #$7F
	STA $0207
	LDA #$A5
	STA $0208
	LDA #$1F
	STA $0209
	LDA #$CE
	STA $020A
	LDA #$29
	STA $020B
	LDA #$42
	STA $020C
	LDA #$6C
	STA $020D
	LDA #$42
	STA $020E
	; 0000200 55aa ff6e 4233 9d7f a51f ce29 426c 4200


; expected result: $022A = 0x55
test00:
	.byte $FF      ; DebugConsole
	LDA #85				 ; $55
	LDX #42        ; $2A
	LDY #115       ; $73
	STA $81
	LDA #$01
	STA $61
	LDA #$7E
	LDA $81        ; $55
	STA $0910
	LDA #$7E
	LDA $0910      ; $55
	STA $56,X      ; 0080 = 55
	LDA #$7E
	LDA $56,X
	STY $60					; 0060 = 73
	.byte $FF       ; DebugConsole ; $4070
	STA ($60),Y     ; 01e6 = 55
	LDA #$7E
	LDA ($60),Y     ; $55
	STA $07ff,X     ; 0829 = 55
	LDA #$7E
	LDA $07ff,X
	STA $07ff,Y     ; 0872 = 55
	LDA #$7E
	LDA $07ff,Y
	STA ($36,X)     ; 0060 = 55
	LDA #$7E
	LDA ($36,X)     ; correct up to here?
	.byte $FF       ; Debug Console
	STX $50
	LDX $60
	LDY $50
	STX $0913
	LDX #$22
	LDX $0913
	STY $0914
	LDY #$99
	LDY $0914
	STY $2D,X
	STX $77,Y
	LDY #$99
	LDY $2D,X
	LDX #$22
	LDX $77,Y
	LDY #$99
	LDY $08A0,X
	LDX #$22
	LDX $08A1,Y
	STA $0200,X
	.byte $FF