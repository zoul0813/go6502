# OpCodes

List of high level opcodes, and whether they've been implemented

- [x] ADC - add with carry
- [x] AND - and (with accumulator)
- [x] ASL - arithmetic shift left
- [x] BCC - branch on carry clear
- [x] BCS - branch on carry set
- [x] BEQ - branch on equal (zero set)
- [x] BIT - bit test
- [x] BMI - branch on minus (negative set)
- [x] BNE - branch on not equal (zero clear)
- [x] BPL - branch on plus (negative clear)
- [x] BRK - break / interrupt
- [x] BVC - branch on overflow clear
- [x] BVS - branch on overflow set
- [x] CLC - clear carry
- [x] CLD - clear decimal
- [x] CLI - clear interrupt disable
- [x] CLV - clear overflow
- [x] CMP - compare (with accumulator)
- [x] CPX - compare with X
- [x] CPY - compare with Y
- [x] DEC - decrement
- [x] DEX - decrement X
- [x] DEY - decrement Y
- [x] EOR - exclusive or (with accumulator)
- [x] INC - increment
- [x] INX - increment X
- [x] INY - increment Y
- [x] JMP - jump
- [x] JSR - jump subroutine
- [x] LDA - load accumulator
- [x] LDX - load X
- [x] LDY - load Y
- [x] LSR - logical shift right
- [x] NOP - no operation
- [x] ORA - or with accumulator
- [x] PHA - push accumulator
- [x] PHP - push processor status (SR)
- [x] PLA - pull accumulator
- [x] PLP - pull processor status (SR)
- [x] ROL - rotate left
- [x] ROR - rotate right
- [x] RTI - return from interrupt
- [x] RTS - return from subroutine
- [x] SBC - subtract with carry
- [x] SEC - set carry
- [x] SED - set decimal
- [x] SEI - set interrupt disable
- [x] STA - store accumulator
- [x] STX - store X
- [x] STY - store Y
- [x] TAX - transfer accumulator to X
- [x] TAY - transfer accumulator to Y
- [x] TSX - transfer stack pointer to X
- [x] TXA - transfer X to accumulator
- [x] TXS - transfer X to stack pointer
- [x] TYA - transfer Y to accumulator

## Special OpCodes supported

- [x] 0xFF - Debug Console
