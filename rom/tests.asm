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
	; .byte $FF, $FE     ; DebugConsole
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
	; .byte $FF, $FE      ; DebugConsole ; $4070
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
	; .byte $FF,       ; Debug Console
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
	; .byte $FF, ; DebugConsole
	LDY #$99
	LDY $2D,X ;
	LDX #$22
	LDX $77,Y ; 55?
	LDY #$99 ;
	LDY $08A0,X
	LDX #$22
	LDX $08A1,Y
	STA $0200,X
	.byte $FF, $00 ; expected result: $022A = 0x55
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
	LDA #231 ; $E7
	STA $11
	LDA #57 ; $39
	STA $12
	LDA $99 ; 01101000 ; $68
	AND $10 ; 10111001
	;      -> 00101000
	ORA $11 ; 11100111
	;      -> 11101111
	EOR $12 ; 00111001
	;      -> 11010110
	;       = 11010110

	; zpx
	LDX #16 ; $10
	STA $99 ; D6 at $99
	LDA #188 ; $BC
	STA $20
	LDA #49 ; $31
	STA $21
	LDA #32 ; $20
	STA $22
	LDA $99   ; 11010110 ($D6 in A)
	AND $10,X ; 10111001
	;        -> 10010000
	ORA $11,X
	EOR $12,X ; $95 in A

	; abs
	STA $99
	LDA #111 ; $6F
	STA $0110
	LDA #60 ; $3C
	STA $0111
	LDA #39 ; $27
	STA $0112 ; //
	LDA $99 ; $95
	AND $0110
	ORA $0111
	EOR $0112 ; $1A in A

	; abx
	STA $99
	LDA #138 ; $8A
	STA $0120
	LDA #71 ; $47
	STA $0121
	LDA #143 ; $8F
	STA $0122
	LDA $99 ; $1A
	; X == $10 (16)
	AND $0110,X ; //
	ORA $0111,X ; //
	EOR $0112,X ; //

	; aby
	LDY #32 ; $20
	STA $99 ; $C0 at $99
	LDA #115 ; $73
	STA $0130
	LDA #42 ; $2A
	STA $0131
	LDA #241 ; $F1
	STA $0132
	LDA $99 ; $C0
	AND $0110,Y ;
	ORA $0111,Y ;
	EOR $0112,Y ;

	; idx
	STA $99 ; $9B
	LDA #112 ; $70
	STA $30
	LDA #$01 ; $01
	STA $31
	LDA #113 ; $71
	STA $32
	LDA #$01 ; $01
	STA $33
	LDA #114 ; $72
	STA $34
	LDA #$01 ; $01
	STA $35
	LDA #197 ; $C5
	STA $0170
	LDA #124 ; $7C
	STA $0171
	LDA #161 ; $A1
	STA $0172
	LDA $99 ; $9B
	AND ($20,X) ;
	ORA ($22,X) ;
	EOR ($24,X) ;

	; idy
	STA $99 ; $5C
	LDA #96 ; $60
	STA $40
	LDA #$01
	STA $41
	LDA #97 ; $61
	STA $42
	LDA #$01
	STA $43
	LDA #98 ; $62
	STA $44
	LDA #$01
	STA $45
	LDA #55 ; $37
	STA $0250
	LDA #35 ; $23
	STA $0251
	LDA #157 ; $9D
	STA $0252
	LDA $99 ; $5C
	LDY #$F0 ; $F0
	AND ($40),Y ;
	ORA ($42),Y ;
	EOR ($44),Y ;

	STA $A9 ; $AA

	.byte $FF, $01 ; expected result: $A9 = 0xAA
; CHECK test01
	LDA $A9
	CMP $0201
	BEQ test02
	LDA #$01
	STA $0210
	JMP theend

; expected result: $71 = 0xFF
test02:
	LDA #$FF
	LDX #$00

	STA $90 ; $FF
	INC $90 ; $00
	INC $90 ; $01
	LDA $90 ; $01
	LDX $90 ; $01

	STA $90,X ;
	INC $90,X ;
	LDA $90,X ; $02
	LDX $91 ; $02

	STA $0190,X ;
	INC $0192 ;
	LDA $0190,X ; $03
	LDX $0192 ; $03

	STA $0190,X ;
	INC $0190,X ;
	LDA $0190,X ; $04
	LDX $0193 ; $04

	STA $0170,X ;
	DEC $0170,X ;
	LDA $0170,X ;
	LDX $0174 ;

	STA $0170,X ;
	DEC $0173 ;
	LDA $0170,X ;
	LDX $0173 ;

	STA $70,X ;
	DEC $70,X ;
	LDA $70,X ;
	LDX $72 ;

	STA $70,X
	DEC $71 ;
	DEC $71 ; $FF

	.byte $FF, $02 ; expected result: $71 = 0xFF
; CHECK test02
	LDA $71
	CMP $0202
	BEQ test03
	LDA #$02
	STA $0210
	JMP theend

; expected result: $01DD = 0x6E
test03:
	LDA #$4B ;
	LSR ; $25
	ASL ; $4A

	STA $50 ; $4A
	ASL $50 ; $94
	ASL $50 ; $28
	LSR $50 ; $14
	LDA $50 ; $14

	LDX $50
	ORA #$C9 ; $DD
	STA $60 ; $DD
	ASL $4C,X ; $BA
	LSR $4C,X ; $5D
	LSR $4C,X ; $2E
	LDA $4C,X ; $2E

	LDX $60 ; $2E
	ORA #$41 ; $6F
	STA $012E ;
	LSR $0100,X ; $37
	LSR $0100,X ; $1B
	ASL $0100,X ; $36
	LDA $0100,X ; $36


	LDX $012E ; $41
	ORA #$81
	STA $0100,X
	LSR $0136
	LSR $0136
	ASL $0136
	LDA $0100,X

	; rol & ror

	ROL
	ROL
	ROR
	STA $70

	LDX $70
	ORA #$03
	STA $0C,X
	ROL $C0
	ROR $C0
	ROR $C0
	LDA $0C,X

	LDX $C0
	STA $D0
	ROL $75,X
	ROL $75,X
	ROR $75,X
	LDA $D0

	LDX $D0
	STA $0100,X
	ROL $01B7
	ROL $01B7
	ROL $01B7
	ROR $01B7
	LDA $0100,X

	LDX $01B7
	STA $01DD
	ROL $0100,X
	ROR $0100,X
	ROR $0100,X
	.byte $FF, $03 ; expected result: $01DD = 0x6E

; CHECK test03
	LDA $01DD
	CMP $0203
	BEQ test04
	LDA #$03
	STA $0210
	JMP theend


; expected result: $40 = 0x42
test04:
	;; location of final label
	LDA #<final; #$75
	; LDA #$E8 ;originally:#$7C
	STA $20
	LDA #>final ;#$40
	; LDA #$42 ;originally:#$02
	STA $21
	LDA #$00
	ORA #$03
	JMP jump1
	ORA #$FF ; not done
jump1:
	ORA #$30
	JSR subr
	ORA #$42
	JMP ($0020)
	ORA #$FF ; not done
subr:
	STA $30
	LDX $30
	LDA #$00
	RTS
final:
	STA $0D,X

	.byte $FF, $04 ; expected result: $40 = 0x42
; CHECK test04
	LDA $40
	CMP $0204
	BEQ test05
	LDA #$04
	STA $0210
	JMP theend

; expected result: $40 = 0x33
test05:
	LDA #$35

	TAX
	DEX
	DEX
	INX
	TXA

	TAY
	DEY
	DEY
	INY
	TYA

	TAX
	LDA #$20
	TXS
	LDX #$10
	TSX
	TXA

	STA $40
	.byte $FF, $05 ; expected result: $40 = 0x33
; CHECK test05
	LDA $40
	CMP $0205
	BEQ test06
	LDA #$05
	STA $0210
	JMP theend


; expected result: $30 = 9D
test06:
	.byte $FF, $FE
; RESET TO CARRY FLAG = 0
	ROL

	LDA #$6A ;
	STA $50 ;
	LDA #$6B ;
	STA $51 ;
	LDA #$A1 ;
	STA $60 ;
	LDA #$A2 ;
	STA $61 ;

	LDA #$FF ;
	ADC #$FF ; $FE
	ADC #$FF ; $FD
	SBC #$AE ; $4F

	STA $40 ;
	LDX $40 ; $4F?
	ADC $00,X ; $00 + $4F = $00?
	SBC $01,X ; ??

	ADC $60
	SBC $61

	STA $0120
	LDA #$4D
	STA $0121
	LDA #$23
	ADC $0120
	SBC $0121

	STA $F0
	LDX $F0
	LDA #$64
	STA $0124
	LDA #$62
	STA $0125
	LDA #$26
	ADC $0100,X
	SBC $0101,X

	STA $F1
	LDY $F1
	LDA #$E5
	STA $0128
	LDA #$E9
	STA $0129
	LDA #$34
	ADC $0100,Y
	SBC $0101,Y

	STA $F2
	LDX $F2
	LDA #$20
	STA $70
	LDA #$01
	STA $71
	LDA #$24
	STA $72
	LDA #$01
	STA $73
	ADC ($41,X)
	SBC ($3F,X)

	STA $F3
	LDY $F3
	LDA #$DA
	STA $80
	LDA #$00
	STA $81
	LDA #$DC
	STA $82
	LDA #$00
	STA $83
	LDA #$AA
	ADC ($80),Y
	SBC ($82),Y
	STA $30

	.byte $FF, $06; expected result: $30 = 9D
; CHECK test06
	LDA $30
	CMP $0206
	BEQ test07
	LDA #$06
	STA $0210
	JMP theend

test07:
	LDA #0
	LDX #0
	LDY #0

theend:
	; EXPECTED FINAL RESULTS: $0210 = FF
  .byte $FF, $FF    ; DebugConsole
	JMP theend