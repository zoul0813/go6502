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
  ; ;; internal test
	; LDA #$FF ; X
	; STA $f0 ; X
	; LDA #$50 ; X
	; LDX #$11 ; X
	; LDY #$22 ; X
	; STA $00 ; X
	; LDA #$51 ; X
	; STA $00,X ; X
	; LDA #$52 ; X
	; STA $0022 ; X
	; ; .byte $FF
	; LDA #$53 ; X
	; STA $0020,X ; X
	; LDA #$54 ; X
	; STA $0030,Y ; X
	; LDA #$55 ; X
	; LDX #$00 ; X
	; STA ($00,X)
	; LDA #$56
	; LDY #$01
	; STA ($00),Y

	; .byte $FF ; DebugConsole
	; ;; quit here


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
	; .byte $FF      ; DebugConsole
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
	; .byte $FF       ; DebugConsole ; $4070
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
	; .byte $FF       ; Debug Console
	STX $50
	LDX $60
	LDY $50					;
	STX $0913
	LDX #$22        ;
	LDX $0913       ;
	STY $0914
	LDY #$99
	LDY $0914       ;
	STY $2D,X       ;
	STX $77,Y       ; 00a1 = 55
	; .byte $FF ; DebugConsole
	LDY #$99
	LDY $2D,X ;
	LDX #$22
	LDX $77,Y ; 55?
	LDY #$99 ;
	LDY $08A0,X
	LDX #$22
	LDX $08A1,Y
	STA $0200,X
	; .byte $FF ; DebugConsole
; CHECK test00:
	LDA $022A
	CMP $0200
	BEQ test00pass
	JMP theend
test00pass:
	LDA #$FE
	STA $0210


; expected result: $A9 = 0xAA
test01:
	; imm
	LDA #85  ; 01010101
	AND #83  ; 01010011
	;       -> 01010001
	ORA #56  ; 00111000
	;       -> 01111001
	EOR #17  ; 00010001
	;       -> 01101000 = 68


	; zpg
	STA $99
	LDA #185
	STA $10
	LDA #231
	STA $11
	LDA #57
	STA $12
	LDA $99 ; 01101000
	AND $10 ; 10111001
	;      -> 00101000
	ORA $11 ; 11100111
	;      -> 11101111
	EOR $12 ; 00111001
	;      -> 11010110
	;       = 11010110

	.byte $FF ; DebugConsole
	; zpx
	LDX #16
	STA $99
	LDA #188
	STA $20
	LDA #49
	STA $21
	LDA #32
	STA $22
	LDA $99   ; 11010110 (d6 in A)
	AND $10,X ; 10111001
	;        -> 10010000
	ORA $11,X
	EOR $12,X

	; abs
	STA $99
	LDA #111
	STA $0110
	LDA #60
	STA $0111
	LDA #39
	STA $0112
	LDA $99
	AND $0110
	ORA $0111
	EOR $0112

	; abx
	STA $99
	LDA #138
	STA $0120
	LDA #71
	STA $0121
	LDA #143
	STA $0122
	LDA $99
	AND $0110,X
	ORA $0111,X
	EOR $0112,X

	; aby
	LDY #32
	STA $99
	LDA #115
	STA $0130
	LDA #42
	STA $0131
	LDA #241
	STA $0132
	LDA $99
	AND $0110,Y
	ORA $0111,Y
	EOR $0112,Y

	; idx
	STA $99
	LDA #112
	STA $30
	LDA #$01
	STA $31
	LDA #113
	STA $32
	LDA #$01
	STA $33
	LDA #114
	STA $34
	LDA #$01
	STA $35
	LDA #197
	STA $0170
	LDA #124
	STA $0171
	LDA #161
	STA $0172
	LDA $99
	AND ($20,X)
	ORA ($22,X)
	EOR ($24,X)

	; idy
	STA $99
	LDA #96
	STA $40
	LDA #$01
	STA $41
	LDA #97
	STA $42
	LDA #$01
	STA $43
	LDA #98
	STA $44
	LDA #$01
	STA $45
	LDA #55
	STA $0250
	LDA #35
	STA $0251
	LDA #157
	STA $0252
	LDA $99
	LDY #$F0
	AND ($40),Y
	ORA ($42),Y
	EOR ($44),Y

	STA $A9

; CHECK test01
	.byte $FF
	LDA $A9
	CMP $0201
	BEQ test02
	LDA #$01
	STA $0210
	JMP theend

test02:
	LDA #$00
	LDX #$00
	LDY #$00

theend:
  .byte $FF     ; DebugConsole
	JMP theend