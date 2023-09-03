package CPU

type OpCode uint8

/**
I = immediate
INX = X, indirect
INY = indirect, Y
ZP = Zero Page
ZPX = Zero Page, X-indexed
ZPY = Zerp Page, Y-Indexed
X = X-index, indirect
A = absolute
AX = absolute, X index
AY = absolute, Y index
R = relative
*/

const (
	BRK     OpCode = 0x00
	ORA_X          = 0x01
	ORA_ZP         = 0x05
	ASL_ZP         = 0x06
	PHP            = 0x08
	ORA_I          = 0x09
	ASL            = 0x0A
	ORA_A          = 0x0D
	ASL_A          = 0x0E
	BPL_R          = 0x10
	ORA_INY        = 0x11
	ORA_ZPX        = 0x15
	ASL_ZPX        = 0x16
	CLC            = 0x18
	ORA_AY         = 0x19
	ORA_AX         = 0x1D
	ASL_AX         = 0x1E

	SEC     = 0x38
	CLI     = 0x58
	ADC_INX = 0x61
	ADC_ZP  = 0x65
	ADC_I   = 0x69
	JMP_IN  = 0x6C
	ADC_A   = 0x6D
	ADC_INY = 0x71
	ADC_ZPX = 0x75
	SEI     = 0x78
	ADC_AY  = 0x79
	ADC_AX  = 0x7D

	STA_INX = 0x81
	STA_ZP  = 0x85
	DEY     = 0x88
	TXA     = 0x8A
	STA_A   = 0x8D
	STA_INY = 0x91
	STA_ZPX = 0x95
	TYA     = 0x98
	STA_AY  = 0x99
	STA_AX  = 0x9D

	// skip a few
	JMP_A = 0x4C
	// skip a few
	LDY_I   = 0xA0
	LDA_IX  = 0xA1
	LDX_I   = 0xA2
	LDY_ZP  = 0xA4
	LDA_ZP  = 0xA5
	LDX_ZP  = 0xA6
	TAY     = 0xA8
	LDA_I   = 0xA9
	TAX     = 0xAA
	LDY_A   = 0xAC
	LDA_A   = 0xAD
	LDX_A   = 0xAE
	BCS_R   = 0xB0
	LDA_INY = 0xB1
	LDY_ZPX = 0xB4
	LDA_ZPX = 0xB5
	LDX_ZPY = 0xB6
	CLV     = 0xB8
	LDA_AY  = 0xB9
	LDY_AX  = 0xBC
	LDA_AX  = 0xBD
	LDX_AY  = 0xBE
	DEX     = 0xCA
	DEC_ZP  = 0xC6
	INY     = 0xC8
	DEC_A   = 0xCE
	DEC_ZPX = 0xD6
	CLD     = 0xD8
	DEC_AX  = 0xDE
	INC_ZP  = 0xE6
	INX     = 0xE8
	INC_A   = 0xEE
	INC_ZPX = 0xF6
	SED     = 0xF8
	INC_AX  = 0xFE
)
