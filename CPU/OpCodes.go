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

	SEC = 0x38

	// skip a few
	JMP_A = 0x4C
	// skip a few
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
	LDA_AY  = 0xB9
	LDA_AX  = 0xBD
	INY     = 0xC8
	INX     = 0xE8
)
