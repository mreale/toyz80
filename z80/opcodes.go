package z80

// mode is the instruction's addressing mode.
type mode int

const (
	invalid   mode = iota // invalid addressing mode
	immediate             // source only
	immediateExtended
	modifiedPageZero
	relative
	extended
	indexed
	register
	implied
	registerIndirect
	bit
	condition
	displacement
	indirect
)

// opcode describes an instruction.
type opcode struct {
	noCycles  uint64   // cycles instruction costs
	mnemonic  []string // human readable mnemonic z80 = 0, 8080 = 1
	noBytes   uint16   // no of bytes per instruction
	src       mode     // source addressing mode
	srcR      []string // source disassembly cheat
	dst       mode     // destination addressing mode
	dstR      []string // destination disassembly cheat
	multiByte bool     // use alternative array instead for lookup
}

// opcodes are all possible instructions 16 bit instructions.  We just throw
// mmemory at the problem.
var (
	opcodesCB = []opcode{
		0x27: {
			mnemonic: []string{"sla", ""},
			dst:      register,
			dstR:     []string{"a"},
			noBytes:  2,
			noCycles: 8,
		},
		0x3f: {
			mnemonic: []string{"srl", ""},
			dst:      register,
			dstR:     []string{"a"},
			noBytes:  2,
			noCycles: 8,
		},
	}
	opcodesED = []opcode{
		0x44: {
			mnemonic: []string{"neg", ""},
			noBytes:  2,
			noCycles: 8,
		},
	}
	opcodesDD = []opcode{
		0x23: {
			mnemonic: []string{"inc"},
			dst:      register,
			dstR:     []string{"ix"},
			noBytes:  2,
			noCycles: 10,
		},
		0x2b: {
			mnemonic: []string{"dec"},
			dst:      register,
			dstR:     []string{"ix"},
			noBytes:  2,
			noCycles: 10,
		},
		0xe1: opcode{
			mnemonic: []string{"pop"},
			dst:      register,
			dstR:     []string{"ix"},
			noBytes:  2,
			noCycles: 14,
		},
		0xe5: opcode{
			mnemonic: []string{"push"},
			dst:      register,
			dstR:     []string{"ix"},
			noBytes:  2,
			noCycles: 15,
		},
	}
	opcodesFD = []opcode{
		0x23: {
			mnemonic: []string{"inc"},
			dst:      register,
			dstR:     []string{"iy"},
			noBytes:  2,
			noCycles: 10,
		},
		0x2b: {
			mnemonic: []string{"dec"},
			dst:      register,
			dstR:     []string{"iy"},
			noBytes:  2,
			noCycles: 10,
		},
		0xe1: opcode{
			mnemonic: []string{"pop"},
			dst:      register,
			dstR:     []string{"iy"},
			noBytes:  2,
			noCycles: 14,
		},
		0xe5: opcode{
			mnemonic: []string{"push"},
			dst:      register,
			dstR:     []string{"iy"},
			noBytes:  2,
			noCycles: 15,
		},
	}
)

// Opcodes are all possible instructions 8 bit instructions.
var (
	opcodes = []opcode{
		// 0x00 nop
		opcode{
			mnemonic: []string{"nop", "nop"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x01 ld bc,nn
		opcode{
			mnemonic: []string{"ld", "lxi"},
			dst:      register,
			dstR:     []string{"bc", "b"},
			src:      immediateExtended,
			srcR:     []string{"", ""},
			noBytes:  3,
			noCycles: 10,
		},
		// 0x02 ld (bc),a
		opcode{
			mnemonic: []string{"ld", "stax"},
			dst:      registerIndirect,
			dstR:     []string{"bc", "b"},
			src:      register,
			srcR:     []string{"a", ""},
			noBytes:  1,
			noCycles: 7,
		},
		// 0x03 inc bc
		opcode{
			mnemonic: []string{"inc", "inx"},
			dst:      register,
			dstR:     []string{"bc", "b"},
			noBytes:  1,
			noCycles: 6,
		},
		// 0x04 inc b
		opcode{
			mnemonic: []string{"inc", "inr"},
			dst:      register,
			dstR:     []string{"b", "b"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x05 dec b
		opcode{
			mnemonic: []string{"dec", "dcr"},
			dst:      register,
			dstR:     []string{"b", "b"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x06 ld b,n
		opcode{
			mnemonic: []string{"ld", "mvi"},
			dst:      register,
			dstR:     []string{"b", "b"},
			src:      immediate,
			noBytes:  2,
			noCycles: 7,
		},
		// 0x07 rlca
		opcode{
			mnemonic: []string{"rlca", "rlc"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x08 ex af,af'
		opcode{
			mnemonic: []string{"ex"},
			dst:      register,
			dstR:     []string{"af"},
			src:      register,
			srcR:     []string{"af'"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x09 add hl,bc
		opcode{
			mnemonic: []string{"add", "dad"},
			dst:      register,
			dstR:     []string{"hl", "b"},
			src:      register,
			srcR:     []string{"bc", ""},
			noBytes:  1,
			noCycles: 11,
		},
		// 0x0a ld a,(bc)
		opcode{
			mnemonic: []string{"ld", "ldax"},
			dst:      register,
			dstR:     []string{"a", ""},
			src:      registerIndirect,
			srcR:     []string{"bc", "b"},
			noBytes:  1,
			noCycles: 7,
		},
		// 0x0b dec bc
		opcode{
			mnemonic: []string{"dec", "dcx"},
			dst:      register,
			dstR:     []string{"bc", "d"},
			noBytes:  1,
			noCycles: 6,
		},
		// 0x0c inc c
		opcode{
			mnemonic: []string{"inc", "inr"},
			dst:      register,
			dstR:     []string{"c", "c"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x0d dec c
		opcode{
			mnemonic: []string{"dec", "dcr"},
			dst:      register,
			dstR:     []string{"c", "c"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x0e ld c,n
		opcode{
			mnemonic: []string{"ld", "mvi"},
			dst:      register,
			dstR:     []string{"c", "c"},
			src:      immediate,
			noBytes:  2,
			noCycles: 7,
		},
		// 0x0f rrca
		opcode{
			mnemonic: []string{"rrca", "rrc"},
			noBytes:  1,
			noCycles: 4,
		},

		// 0x10
		opcode{},
		// 0x11 ld de,nn
		opcode{
			mnemonic: []string{"ld", "lxi"},
			dst:      register,
			dstR:     []string{"de", "d"},
			src:      immediateExtended,
			noBytes:  3,
			noCycles: 10,
		},
		// 0x12 ld (de),a
		opcode{
			mnemonic: []string{"ld", "stax"},
			dst:      registerIndirect,
			dstR:     []string{"de", "d"},
			src:      register,
			srcR:     []string{"a", ""},
			noBytes:  1,
			noCycles: 7,
		},
		// 0x13 inc de
		opcode{
			mnemonic: []string{"inc", "inx"},
			dst:      register,
			dstR:     []string{"de", "d"},
			noBytes:  1,
			noCycles: 6,
		},
		// 0x14 inc d
		opcode{
			mnemonic: []string{"inc", "inr"},
			dst:      register,
			dstR:     []string{"d", "d"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x15 dec d
		opcode{
			mnemonic: []string{"dec", "dcr"},
			dst:      register,
			dstR:     []string{"d", "d"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x16 ld d,n
		opcode{
			mnemonic: []string{"ld", "mvi"},
			dst:      register,
			dstR:     []string{"d", "d"},
			src:      immediate,
			noBytes:  2,
			noCycles: 7,
		},
		// 0x17
		opcode{},
		// 0x18
		opcode{
			mnemonic: []string{"jr", ""},
			dst:      displacement,
			noBytes:  2,
			noCycles: 12,
		},
		// 0x19 add hl,de
		opcode{
			mnemonic: []string{"add", "dad"},
			dst:      register,
			dstR:     []string{"hl", "d"},
			src:      registerIndirect,
			srcR:     []string{"de", ""},
			noBytes:  1,
			noCycles: 11,
		},
		// 0x1a ld a,(de)
		opcode{
			mnemonic: []string{"ld", "ldax"},
			dst:      register,
			dstR:     []string{"a", ""},
			src:      registerIndirect,
			srcR:     []string{"de", "d"},
			noBytes:  1,
			noCycles: 7,
		},
		// 0x1b dec de
		opcode{
			mnemonic: []string{"dec", "dcx"},
			dst:      register,
			dstR:     []string{"de", "d"},
			noBytes:  1,
			noCycles: 6,
		},
		// 0x1c inc e
		opcode{
			mnemonic: []string{"inc", "inr"},
			dst:      register,
			dstR:     []string{"e", "e"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x1d dec e
		opcode{
			mnemonic: []string{"dec", "dcr"},
			dst:      register,
			dstR:     []string{"e", "e"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x1e ld e,n
		opcode{
			mnemonic: []string{"ld", "mvi"},
			dst:      register,
			dstR:     []string{"e", "e"},
			src:      immediate,
			noBytes:  2,
			noCycles: 7,
		},
		// 0x1f
		opcode{
			mnemonic: []string{"rra", "rar"},
			noBytes:  1,
			noCycles: 4,
		},

		// 0x20
		opcode{},
		// 0x21 ld hl,nn
		opcode{
			mnemonic: []string{"ld", "lxi"},
			dst:      register,
			dstR:     []string{"hl", "h"},
			src:      immediateExtended,
			noBytes:  3,
			noCycles: 10,
		},
		// 0x22 ld (nn),hl
		opcode{
			mnemonic: []string{"ld", "shld"},
			dst:      extended,
			src:      register,
			srcR:     []string{"hl", ""},
			noBytes:  3,
			noCycles: 16,
		},
		// 0x23 inc de
		opcode{
			mnemonic: []string{"inc", "inx"},
			dst:      register,
			dstR:     []string{"hl", "h"},
			noBytes:  1,
			noCycles: 6,
		},
		// 0x24 inc h
		opcode{
			mnemonic: []string{"inc", "inr"},
			dst:      register,
			dstR:     []string{"h", "h"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x25 dec h
		opcode{
			mnemonic: []string{"dec", "dcr"},
			dst:      register,
			dstR:     []string{"h", "h"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x26 ld h,n
		opcode{
			mnemonic: []string{"ld", "mvi"},
			dst:      register,
			dstR:     []string{"h", "h"},
			src:      immediate,
			noBytes:  2,
			noCycles: 7,
		},
		// 0x27
		opcode{},
		// 0x28 // jr z,d
		opcode{
			mnemonic: []string{"jr", ""},
			dst:      condition,
			dstR:     []string{"z", ""},
			src:      displacement,
			noBytes:  2,
			noCycles: 12, // XXX or 7
		},
		// 0x29 add hl,hl
		opcode{
			mnemonic: []string{"add", "dad"},
			dst:      register,
			dstR:     []string{"hl", "h"},
			src:      register,
			srcR:     []string{"hl", ""},
			noBytes:  1,
			noCycles: 11,
		},
		// 0x2a ld hl,(nn)
		opcode{
			mnemonic: []string{"ld", "lhld"},
			dst:      register,
			dstR:     []string{"hl", ""},
			src:      extended,
			noBytes:  3,
			noCycles: 16,
		},
		// 0x2b dec hl
		opcode{
			mnemonic: []string{"dec", "dcx"},
			dst:      register,
			dstR:     []string{"hl", "h"},
			noBytes:  1,
			noCycles: 6,
		},
		// 0x2c inc l
		opcode{
			mnemonic: []string{"inc", "inr"},
			dst:      register,
			dstR:     []string{"l", "l"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x2d dec l
		opcode{
			mnemonic: []string{"dec", "dcr"},
			dst:      register,
			dstR:     []string{"l", "l"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x2e ld l,n
		opcode{
			mnemonic: []string{"ld", "mvi"},
			dst:      register,
			dstR:     []string{"l", "l"},
			src:      immediate,
			noBytes:  2,
			noCycles: 7,
		},
		// 0x2f
		opcode{
			mnemonic: []string{"cpl", "cma"},
			noBytes:  1,
			noCycles: 4,
		},

		// 0x30
		opcode{},
		// 0x31 ld sp,nn
		opcode{
			mnemonic: []string{"ld", "lxi"},
			dst:      register,
			dstR:     []string{"sp", "sp"},
			src:      immediateExtended,
			noBytes:  3,
			noCycles: 10,
		},
		// 0x32 ld (nn),a
		opcode{
			mnemonic: []string{"ld", "sta"},
			dst:      immediateExtended,
			src:      register,
			srcR:     []string{"a", ""},
			noBytes:  3,
			noCycles: 13,
		},
		// 0x33 inc de
		opcode{
			mnemonic: []string{"inc", "inx"},
			dst:      register,
			dstR:     []string{"sp", "h"},
			noBytes:  1,
			noCycles: 6,
		},
		// 0x34 inc (hl)
		opcode{
			mnemonic: []string{"inc", "inr"},
			dst:      registerIndirect,
			dstR:     []string{"hl", "m"},
			noBytes:  1,
			noCycles: 11,
		},
		// 0x35 dec (hl)
		opcode{
			mnemonic: []string{"dec", "dec"},
			dst:      registerIndirect,
			dstR:     []string{"hl", "m"},
			noBytes:  1,
			noCycles: 11,
		},
		// 0x36 ld (hl),n
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      registerIndirect,
			dstR:     []string{"hl", "m"},
			src:      immediate,
			noBytes:  2,
			noCycles: 10,
		},
		// 0x37
		opcode{
			mnemonic: []string{"scf", "stc"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x38
		opcode{},
		// 0x39 add hl,sp
		opcode{
			mnemonic: []string{"add", "dad"},
			dst:      register,
			dstR:     []string{"hl", "sp"},
			src:      registerIndirect,
			srcR:     []string{"sp", ""},
			noBytes:  1,
			noCycles: 11,
		},
		// 0x3a ld a,(nn)
		opcode{
			mnemonic: []string{"ld", "lda"},
			dst:      register,
			dstR:     []string{"a", ""},
			src:      extended,
			noBytes:  3,
			noCycles: 7,
		},
		// 0x3b dec sp
		opcode{
			mnemonic: []string{"dec", "dcx"},
			dst:      register,
			dstR:     []string{"sp", "sp"},
			noBytes:  1,
			noCycles: 6,
		},
		// 0x3c
		opcode{
			mnemonic: []string{"inc", "inr"},
			dst:      register,
			dstR:     []string{"a", "a"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x3d dec a
		opcode{
			mnemonic: []string{"dec", "dcr"},
			dst:      register,
			dstR:     []string{"a", "a"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x3e ld a,n
		opcode{
			mnemonic: []string{"ld", "mvi"},
			dst:      register,
			dstR:     []string{"a", "a"},
			src:      immediate,
			noBytes:  2,
			noCycles: 7,
		},
		// 0x3f
		opcode{
			mnemonic: []string{"ccf", "cmc"},
			noBytes:  1,
			noCycles: 4,
		},

		// 0x40 ld b,b
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"b", "b"},
			src:      register,
			srcR:     []string{"b", "b"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x41 ld b,c
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"b", "b"},
			src:      register,
			srcR:     []string{"c", "c"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x42 ld b,d
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"b", "b"},
			src:      register,
			srcR:     []string{"d", "d"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x43 ld b,e
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"b", "b"},
			src:      register,
			srcR:     []string{"e", "e"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x44 ld b,h
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"b", "b"},
			src:      register,
			srcR:     []string{"h", "h"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x45 ld b,l
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"b", "b"},
			src:      register,
			srcR:     []string{"l", "l"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x46 ld b,(hl)
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"b", "b"},
			src:      registerIndirect,
			srcR:     []string{"hl", "m"},
			noBytes:  1,
			noCycles: 7,
		},
		// 0x47 ld b,a
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"b", "b"},
			src:      register,
			srcR:     []string{"a", "a"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x48 ld c,b
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"c", "c"},
			src:      register,
			srcR:     []string{"b", "b"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x49 ld c,c
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"c", "c"},
			src:      register,
			srcR:     []string{"c", "c"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x4a ld c,d
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"c", "c"},
			src:      register,
			srcR:     []string{"d", "d"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x4b ld c,e
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"c", "c"},
			src:      register,
			srcR:     []string{"e", "e"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x4c ld c,h
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"c", "c"},
			src:      register,
			srcR:     []string{"h", "h"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x4d
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"c", "c"},
			src:      register,
			srcR:     []string{"l", "l"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x4e
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"c", "c"},
			src:      registerIndirect,
			srcR:     []string{"hl", "m"},
			noBytes:  1,
			noCycles: 7,
		},
		// 0x4f ld c,a
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"c", "c"},
			src:      register,
			srcR:     []string{"a", "a"},
			noBytes:  1,
			noCycles: 4,
		},

		// 0x50 ld d,b
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"d", "d"},
			src:      register,
			srcR:     []string{"b", "b"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x51 ld d,c
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"d", "d"},
			src:      register,
			srcR:     []string{"c", "c"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x52 ld d,d
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"d", "d"},
			src:      register,
			srcR:     []string{"d", "d"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x53 ld d,e
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"d", "d"},
			src:      register,
			srcR:     []string{"e", "e"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x54 ld d,h
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"d", "d"},
			src:      register,
			srcR:     []string{"h", "h"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x55 ld d,l
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"d", "d"},
			src:      register,
			srcR:     []string{"l", "l"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x56 ld d,(hl)
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"d", "d"},
			src:      registerIndirect,
			srcR:     []string{"hl", "m"},
			noBytes:  1,
			noCycles: 7,
		},
		// 0x57 ld d,a
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"d", "d"},
			src:      register,
			srcR:     []string{"a", "a"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x58 ld e,b
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"e", "e"},
			src:      register,
			srcR:     []string{"b", "b"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x59 ld e,c
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"e", "e"},
			src:      register,
			srcR:     []string{"c", "c"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x5a ld e,d
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"e", "e"},
			src:      register,
			srcR:     []string{"d", "d"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x5b ld e,e
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"e", "e"},
			src:      register,
			srcR:     []string{"e", "e"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x5c ld e,h
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"e", "e"},
			src:      register,
			srcR:     []string{"h", "h"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x5d
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"e", "e"},
			src:      register,
			srcR:     []string{"l", "l"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x5e ld e,(hl)
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"e", "e"},
			src:      registerIndirect,
			srcR:     []string{"hl", "m"},
			noBytes:  1,
			noCycles: 7,
		},
		// 0x5f ld e,a
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"e", "e"},
			src:      register,
			srcR:     []string{"a", "a"},
			noBytes:  1,
			noCycles: 4,
		},

		// 0x60 ld h,b
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"h", "h"},
			src:      register,
			srcR:     []string{"b", "b"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x61 ld h,c
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"h", "h"},
			src:      register,
			srcR:     []string{"c", "c"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x62 ld h,d
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"h", "h"},
			src:      register,
			srcR:     []string{"d", "d"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x63 ld h,e
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"h", "h"},
			src:      register,
			srcR:     []string{"e", "e"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x64 ld h,h
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"h", "h"},
			src:      register,
			srcR:     []string{"h", "h"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x65 ld h,l
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"h", "h"},
			src:      register,
			srcR:     []string{"l", "l"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x66 ld h,(hl)
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"h", "h"},
			src:      registerIndirect,
			srcR:     []string{"hl", "m"},
			noBytes:  1,
			noCycles: 7,
		},
		// 0x67 ld h,a
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"h", "h"},
			src:      register,
			srcR:     []string{"a", "a"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x68 ld l,b
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"l", "l"},
			src:      register,
			srcR:     []string{"b", "b"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x69 ld l,c
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"l", "l"},
			src:      register,
			srcR:     []string{"c", "c"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x6a ld l,d
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"l", "l"},
			src:      register,
			srcR:     []string{"d", "d"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x6b ld l,e
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"l", "l"},
			src:      register,
			srcR:     []string{"e", "e"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x6c ld l,h
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"l", "l"},
			src:      register,
			srcR:     []string{"h", "h"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x6d ld l,l
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"l", "l"},
			src:      register,
			srcR:     []string{"l", "l"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x6e ld l,(hl)
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"l", "l"},
			src:      registerIndirect,
			srcR:     []string{"hl", "m"},
			noBytes:  1,
			noCycles: 7,
		},
		// 0x6f ld l,a
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"l", "l"},
			src:      register,
			srcR:     []string{"a", "a"},
			noBytes:  1,
			noCycles: 4,
		},

		// 0x70 ld (hl),b
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      registerIndirect,
			dstR:     []string{"hl", "hl"},
			src:      register,
			srcR:     []string{"b", "b"},
			noBytes:  1,
			noCycles: 7,
		},
		// 0x71 ld (hl),c
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      registerIndirect,
			dstR:     []string{"hl", "hl"},
			src:      register,
			srcR:     []string{"c", "c"},
			noBytes:  1,
			noCycles: 7,
		},
		// 0x72 ld (hl),d
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      registerIndirect,
			dstR:     []string{"hl", "hl"},
			src:      register,
			srcR:     []string{"d", "d"},
			noBytes:  1,
			noCycles: 7,
		},
		// 0x73 ld (hl),e
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      registerIndirect,
			dstR:     []string{"hl", "hl"},
			src:      register,
			srcR:     []string{"e", "e"},
			noBytes:  1,
			noCycles: 7,
		},
		// 0x74 ld (hl),h
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      registerIndirect,
			dstR:     []string{"hl", "hl"},
			src:      register,
			srcR:     []string{"h", "h"},
			noBytes:  1,
			noCycles: 7,
		},
		// 0x75 ld (hl),l
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      registerIndirect,
			dstR:     []string{"hl", "hl"},
			src:      register,
			srcR:     []string{"l", "l"},
			noBytes:  1,
			noCycles: 7,
		},
		// 0x76 halt
		opcode{
			mnemonic: []string{"halt", "hlt"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x77 ld (hl),a
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      registerIndirect,
			dstR:     []string{"hl", "hl"},
			src:      register,
			srcR:     []string{"a", "a"},
			noBytes:  1,
			noCycles: 7,
		},
		// 0x78 ld a,b
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"a", "a"},
			src:      register,
			srcR:     []string{"b", "b"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x79 ld a,c
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"a", "a"},
			src:      register,
			srcR:     []string{"c", "c"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x7a ld a,d
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"a", "a"},
			src:      register,
			srcR:     []string{"d", "d"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x7b ld a,e
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"a", "a"},
			src:      register,
			srcR:     []string{"e", "e"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x7c ld a,h
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"a", "a"},
			src:      register,
			srcR:     []string{"h", "h"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x7d ld a,l
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"a", "a"},
			src:      register,
			srcR:     []string{"l", "l"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x7e ld a,(hl)
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"a", "a"},
			src:      registerIndirect,
			srcR:     []string{"hl", "m"},
			noBytes:  1,
			noCycles: 7,
		},
		// 0x7f ld a,a
		opcode{
			mnemonic: []string{"ld", "mov"},
			dst:      register,
			dstR:     []string{"a", "a"},
			src:      register,
			srcR:     []string{"a", "a"},
			noBytes:  1,
			noCycles: 4,
		},

		// 0x80 add a,b
		opcode{
			mnemonic: []string{"add", "add"},
			dst:      register,
			dstR:     []string{"a", "b"},
			src:      register,
			srcR:     []string{"b"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x81 add a,c
		opcode{
			mnemonic: []string{"add", "add"},
			dst:      register,
			dstR:     []string{"a", "c"},
			src:      register,
			srcR:     []string{"c"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x82 add a,d
		opcode{
			mnemonic: []string{"add", "add"},
			dst:      register,
			dstR:     []string{"a", "d"},
			src:      register,
			srcR:     []string{"d"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x83 add a,e
		opcode{
			mnemonic: []string{"add", "add"},
			dst:      register,
			dstR:     []string{"a", "e"},
			src:      register,
			srcR:     []string{"e"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x84 add a,h
		opcode{
			mnemonic: []string{"add", "add"},
			dst:      register,
			dstR:     []string{"a", "h"},
			src:      register,
			srcR:     []string{"h"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x85 add a,l
		opcode{
			mnemonic: []string{"add", "add"},
			dst:      register,
			dstR:     []string{"a", "l"},
			src:      register,
			srcR:     []string{"l"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x86 add a,(hl)
		opcode{
			mnemonic: []string{"add", "add"},
			dst:      register,
			dstR:     []string{"a", "a"},
			src:      register,
			srcR:     []string{"a"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x85 add a,a
		opcode{
			mnemonic: []string{"add", "add"},
			dst:      register,
			dstR:     []string{"a", "l"},
			src:      register,
			srcR:     []string{"l"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x88 adc a,b
		opcode{
			mnemonic: []string{"add", "adc"},
			dst:      register,
			dstR:     []string{"a", "b"},
			src:      register,
			srcR:     []string{"b"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x89 adc a,c
		opcode{
			mnemonic: []string{"add", "adc"},
			dst:      register,
			dstR:     []string{"a", "c"},
			src:      register,
			srcR:     []string{"c"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x8a adc a,d
		opcode{
			mnemonic: []string{"add", "adc"},
			dst:      register,
			dstR:     []string{"a", "d"},
			src:      register,
			srcR:     []string{"d"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x8b adc a,e
		opcode{
			mnemonic: []string{"add", "adc"},
			dst:      register,
			dstR:     []string{"a", "e"},
			src:      register,
			srcR:     []string{"e"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x8c adc a,h
		opcode{
			mnemonic: []string{"add", "adc"},
			dst:      register,
			dstR:     []string{"a", "h"},
			src:      register,
			srcR:     []string{"h"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x8d adc a,l
		opcode{
			mnemonic: []string{"add", "adc"},
			dst:      register,
			dstR:     []string{"a", "l"},
			src:      register,
			srcR:     []string{"l"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x8e adc a,(hl)
		opcode{
			mnemonic: []string{"add", "adc"},
			dst:      register,
			dstR:     []string{"a", "m"},
			src:      register,
			srcR:     []string{"(hl)"},
			noBytes:  1,
			noCycles: 7,
		},
		// 0x8f adc a,a
		opcode{
			mnemonic: []string{"add", "adc"},
			dst:      register,
			dstR:     []string{"a", "l"},
			src:      register,
			srcR:     []string{"l"},
			noBytes:  1,
			noCycles: 4,
		},

		// 0x90 sub a,b
		opcode{
			mnemonic: []string{"sub", "sub"},
			dst:      register,
			dstR:     []string{"a", "a"},
			src:      register,
			srcR:     []string{"b", "b"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x91 sub a,c
		opcode{
			mnemonic: []string{"sub", "sub"},
			dst:      register,
			dstR:     []string{"a", "a"},
			src:      register,
			srcR:     []string{"c", "c"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x92 0x91 sub a,d
		opcode{
			mnemonic: []string{"sub", "sub"},
			dst:      register,
			dstR:     []string{"a", "a"},
			src:      register,
			srcR:     []string{"d", "d"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x93 sub a,e
		opcode{
			mnemonic: []string{"sub", "sub"},
			dst:      register,
			dstR:     []string{"a", "a"},
			src:      register,
			srcR:     []string{"e", "e"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x94 sub a,h
		opcode{
			mnemonic: []string{"sub", "sub"},
			dst:      register,
			dstR:     []string{"a", "a"},
			src:      register,
			srcR:     []string{"h", "h"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x95 sub a,l
		opcode{
			mnemonic: []string{"sub", "sub"},
			dst:      register,
			dstR:     []string{"a", "a"},
			src:      register,
			srcR:     []string{"l", "l"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x96 sub a,(hl)
		opcode{
			mnemonic: []string{"sub", "sub"},
			dst:      register,
			dstR:     []string{"a", "a"},
			src:      registerIndirect,
			srcR:     []string{"hl", "m"},
			noBytes:  1,
			noCycles: 7,
		},
		// 0x97 sub a,a
		opcode{
			mnemonic: []string{"sub", "sub"},
			dst:      register,
			dstR:     []string{"a", "a"},
			src:      register,
			srcR:     []string{"a", "a"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x98 sbc a,b
		opcode{
			mnemonic: []string{"sbc", "sbc"},
			dst:      register,
			dstR:     []string{"a", "a"},
			src:      register,
			srcR:     []string{"b", "b"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x99 sbc a,c
		opcode{
			mnemonic: []string{"sbc", "sbc"},
			dst:      register,
			dstR:     []string{"a", "a"},
			src:      register,
			srcR:     []string{"c", "c"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x9a 0x91 sbc a,d
		opcode{
			mnemonic: []string{"sbc", "sbc"},
			dst:      register,
			dstR:     []string{"a", "a"},
			src:      register,
			srcR:     []string{"d", "d"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x9b sbc a,e
		opcode{
			mnemonic: []string{"sbc", "sbc"},
			dst:      register,
			dstR:     []string{"a", "a"},
			src:      register,
			srcR:     []string{"e", "e"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x9c sbc a,h
		opcode{
			mnemonic: []string{"sbc", "sbc"},
			dst:      register,
			dstR:     []string{"a", "a"},
			src:      register,
			srcR:     []string{"h", "h"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x9d sbc a,l
		opcode{
			mnemonic: []string{"sbc", "sbc"},
			dst:      register,
			dstR:     []string{"a", "a"},
			src:      register,
			srcR:     []string{"l", "l"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0x9e sbc a,(hl)
		opcode{
			mnemonic: []string{"sbc", "sbc"},
			dst:      register,
			dstR:     []string{"a", "a"},
			src:      registerIndirect,
			srcR:     []string{"hl", "m"},
			noBytes:  1,
			noCycles: 7,
		},
		// 0x9f sbc a
		opcode{
			mnemonic: []string{"sbc", "sbc"},
			dst:      register,
			dstR:     []string{"a", "a"},
			src:      register,
			srcR:     []string{"a", "a"},
			noBytes:  1,
			noCycles: 4,
		},

		// 0xa0 and a,b
		opcode{
			mnemonic: []string{"and", "ana"},
			dst:      register,
			dstR:     []string{"a", "a"},
			src:      register,
			srcR:     []string{"b", "b"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0xa1 and a,c
		opcode{
			mnemonic: []string{"and", "ana"},
			dst:      register,
			dstR:     []string{"a", "a"},
			src:      register,
			srcR:     []string{"c", "c"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0xa2 and a,d
		opcode{
			mnemonic: []string{"and", "ana"},
			dst:      register,
			dstR:     []string{"a", "a"},
			src:      register,
			srcR:     []string{"d", "d"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0xa3 and a,e
		opcode{
			mnemonic: []string{"and", "ana"},
			dst:      register,
			dstR:     []string{"a", "a"},
			src:      register,
			srcR:     []string{"e", "e"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0xa4 and a,h
		opcode{
			mnemonic: []string{"and", "ana"},
			dst:      register,
			dstR:     []string{"a", "a"},
			src:      register,
			srcR:     []string{"h", "h"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0xa5 adn a,l
		opcode{
			mnemonic: []string{"and", "ana"},
			dst:      register,
			dstR:     []string{"a", "a"},
			src:      register,
			srcR:     []string{"l", "l"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0xa6 and a,(hl)
		opcode{
			mnemonic: []string{"and", "ana"},
			dst:      register,
			dstR:     []string{"a", "a"},
			src:      registerIndirect,
			srcR:     []string{"hl", "m"},
			noBytes:  1,
			noCycles: 7,
		},
		// 0xa7 and a
		opcode{
			mnemonic: []string{"and", "ana"},
			dst:      register,
			dstR:     []string{"a", "a"},
			src:      register,
			srcR:     []string{"a", "a"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0xa8 xor a,b
		opcode{
			mnemonic: []string{"xor", "xra"},
			dst:      register,
			dstR:     []string{"a", "b"},
			src:      register,
			srcR:     []string{"b", ""},
			noBytes:  1,
			noCycles: 4,
		},
		// 0xa9 xor a,c
		opcode{
			mnemonic: []string{"xor", "xra"},
			dst:      register,
			dstR:     []string{"a", "c"},
			src:      register,
			srcR:     []string{"c", ""},
			noBytes:  1,
			noCycles: 4,
		},
		// 0xaa xor a,d
		opcode{
			mnemonic: []string{"xor", "xra"},
			dst:      register,
			dstR:     []string{"a", "d"},
			src:      register,
			srcR:     []string{"d", ""},
			noBytes:  1,
			noCycles: 4,
		},
		// 0xab xor a,e
		opcode{
			mnemonic: []string{"xor", "xra"},
			dst:      register,
			dstR:     []string{"a", "e"},
			src:      register,
			srcR:     []string{"e", ""},
			noBytes:  1,
			noCycles: 4,
		},
		// 0xac xor a,h
		opcode{
			mnemonic: []string{"xor", "xra"},
			dst:      register,
			dstR:     []string{"a", "h"},
			src:      register,
			srcR:     []string{"h", ""},
			noBytes:  1,
			noCycles: 4,
		},
		// 0xad xor a,l
		opcode{
			mnemonic: []string{"xor", "xra"},
			dst:      register,
			dstR:     []string{"a", "h"},
			src:      register,
			srcR:     []string{"h", ""},
			noBytes:  1,
			noCycles: 4,
		},
		// 0xae xor a,(hl)
		opcode{
			mnemonic: []string{"xor", "xra"},
			dst:      register,
			dstR:     []string{"a", "h"},
			src:      registerIndirect,
			srcR:     []string{"hl", "m"},
			noBytes:  1,
			noCycles: 7,
		},
		// 0xaf xor a,a
		opcode{
			mnemonic: []string{"xor", "xra"},
			dst:      register,
			dstR:     []string{"a", "a"},
			src:      register,
			srcR:     []string{"a", ""},
			noBytes:  1,
			noCycles: 4,
		},

		// 0xb0 or b
		opcode{
			mnemonic: []string{"or", "ora"},
			dst:      register,
			dstR:     []string{"b", "b"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0xb1 or c
		opcode{
			mnemonic: []string{"or", "ora"},
			dst:      register,
			dstR:     []string{"c", "c"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0xb2 or d
		opcode{
			mnemonic: []string{"or", "ora"},
			dst:      register,
			dstR:     []string{"d", "d"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0xb3 or e
		opcode{
			mnemonic: []string{"or", "ora"},
			dst:      register,
			dstR:     []string{"e", "e"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0xb4 or h
		opcode{
			mnemonic: []string{"or", "ora"},
			dst:      register,
			dstR:     []string{"h", "h"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0xb5 or l
		opcode{
			mnemonic: []string{"or", "ora"},
			dst:      register,
			dstR:     []string{"l", "l"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0xb6 or (hl)
		opcode{
			mnemonic: []string{"or", "ora"},
			dst:      registerIndirect,
			dstR:     []string{"hl", "m"},
			noBytes:  1,
			noCycles: 7,
		},
		// 0xb7 or a
		opcode{
			mnemonic: []string{"or", "ora"},
			dst:      register,
			dstR:     []string{"a", "a"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0xb8 cp b
		opcode{
			mnemonic: []string{"cp", "cmp"},
			dst:      register,
			dstR:     []string{"b", "b"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0xb9 cp c
		opcode{
			mnemonic: []string{"cp", "cmp"},
			dst:      register,
			dstR:     []string{"c", "c"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0xba cp d
		opcode{
			mnemonic: []string{"cp", "cmp"},
			dst:      register,
			dstR:     []string{"d", "d"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0xbb cp e
		opcode{
			mnemonic: []string{"cp", "cmp"},
			dst:      register,
			dstR:     []string{"e", "e"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0xbc cp h
		opcode{
			mnemonic: []string{"cp", "cmp"},
			dst:      register,
			dstR:     []string{"h", "h"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0xbd cp l
		opcode{
			mnemonic: []string{"cp", "cmp"},
			dst:      register,
			dstR:     []string{"l", "l"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0xbe
		opcode{
			mnemonic: []string{"cp", "cmp"},
			dst:      registerIndirect,
			dstR:     []string{"hl", "m"},
			noBytes:  1,
			noCycles: 7,
		},
		// 0xbf
		opcode{
			mnemonic: []string{"cp", "cmp"},
			dst:      register,
			dstR:     []string{"a", "a"},
			noBytes:  1,
			noCycles: 4,
		},

		// 0xc0 ret nz
		opcode{
			mnemonic: []string{"ret", "rnz"},
			dst:      condition,
			dstR:     []string{"nz", ""},
			noBytes:  1,
			noCycles: 5,
		},
		// 0xc1 pop bc
		opcode{
			mnemonic: []string{"pop", "pop"},
			dst:      register,
			dstR:     []string{"bc", "b"},
			noBytes:  1,
			noCycles: 10,
		},
		// 0xc2
		opcode{
			mnemonic: []string{"jp", "jnz"},
			dst:      condition,
			dstR:     []string{"nz", ""},
			src:      immediateExtended,
			noBytes:  3,
			noCycles: 10,
		},
		// 0xc3 jp, nn
		opcode{
			mnemonic: []string{"jp", "jmp"},
			dst:      immediateExtended,
			noBytes:  3,
			noCycles: 10,
		},
		// 0xc4 call nz
		opcode{
			mnemonic: []string{"call", "call"},
			dst:      condition,
			dstR:     []string{"nz", ""},
			src:      immediateExtended,
			noBytes:  3,
			noCycles: 10,
		},
		// 0xc5 push bc
		opcode{
			mnemonic: []string{"push", "push"},
			dst:      register,
			dstR:     []string{"bc", "b"},
			noBytes:  1,
			noCycles: 11,
		},
		// 0xc6 add a,i
		opcode{
			mnemonic: []string{"add", "adi"},
			dst:      register,
			dstR:     []string{"a", ""},
			src:      immediate,
			noBytes:  2,
			noCycles: 7,
		},
		// 0xc7 rst $00
		opcode{
			mnemonic: []string{"rst", "rst"},
			dst:      implied,
			dstR:     []string{"$00", "0"},
			noBytes:  1,
			noCycles: 11,
		},
		// 0xc8 ret z
		opcode{
			mnemonic: []string{"ret", "rz"},
			dst:      condition,
			dstR:     []string{"z", ""},
			noBytes:  1,
			noCycles: 5,
		},
		// 0xc9 ret
		opcode{
			mnemonic: []string{"ret", "ret"},
			noBytes:  1,
			noCycles: 10,
		},
		// 0xca jp z,nn
		opcode{
			mnemonic: []string{"jp", "jz"},
			dst:      condition,
			dstR:     []string{"z", ""},
			src:      immediateExtended,
			noBytes:  3,
			noCycles: 10,
		},
		// 0xcb
		opcode{
			multiByte: true,
		},
		// 0xcc
		opcode{},
		// 0xcd call nn
		opcode{
			mnemonic: []string{"call", "call"},
			dst:      immediateExtended,
			noBytes:  3,
			noCycles: 17,
		},
		// 0xce adc a,i
		opcode{
			mnemonic: []string{"adc", "aci"},
			dst:      register,
			dstR:     []string{"a"},
			src:      immediate,
			noBytes:  2,
			noCycles: 7,
		},
		// 0xcf rst $08
		opcode{
			mnemonic: []string{"rst", "rst"},
			dst:      implied,
			dstR:     []string{"$08", "1"},
			noBytes:  1,
			noCycles: 11,
		},

		// 0xd0 ret nc
		opcode{
			mnemonic: []string{"ret", "rnc"},
			dst:      condition,
			dstR:     []string{"nc", ""},
			noBytes:  1,
			noCycles: 5,
		},
		// 0xd1 pop de
		opcode{
			mnemonic: []string{"pop", "pop"},
			dst:      register,
			dstR:     []string{"de", "d"},
			noBytes:  1,
			noCycles: 10,
		},
		// 0xd2 jp nc,nn
		opcode{
			mnemonic: []string{"jp", "jnc"},
			dst:      condition,
			dstR:     []string{"nc", ""},
			src:      immediateExtended,
			noBytes:  3,
			noCycles: 10,
		},
		// 0xd3 out (n),a
		opcode{
			mnemonic: []string{"out", "out"},
			dst:      indirect,
			src:      register,
			srcR:     []string{"a", ""},
			noBytes:  2,
			noCycles: 11,
		},
		// 0xd4
		opcode{},
		// 0xd5 push de
		opcode{
			mnemonic: []string{"push", "push"},
			dst:      register,
			dstR:     []string{"de", "d"},
			noBytes:  1,
			noCycles: 11,
		},
		// 0xd6 sub i
		opcode{
			mnemonic: []string{"sub", "sub"},
			dst:      immediate,
			noBytes:  2,
			noCycles: 7,
		},
		// 0xd7 rst $10
		opcode{
			mnemonic: []string{"rst", "rst"},
			dst:      implied,
			dstR:     []string{"$10", "2"},
			noBytes:  1,
			noCycles: 11,
		},
		// 0xd8 ret c
		opcode{
			mnemonic: []string{"ret", "rc"},
			dst:      condition,
			dstR:     []string{"c", ""},
			noBytes:  1,
			noCycles: 5,
		},
		// 0xd9
		opcode{},
		// 0xda jp c,nn
		opcode{
			mnemonic: []string{"jp", "jc"},
			dst:      condition,
			dstR:     []string{"c", ""},
			src:      immediateExtended,
			noBytes:  3,
			noCycles: 10,
		},
		// 0xdb in a,(n)
		opcode{
			mnemonic: []string{"in", "in"},
			dst:      register,
			dstR:     []string{"a", ""},
			src:      indirect,
			noBytes:  2,
			noCycles: 11,
		},
		// 0xdc
		opcode{},
		// 0xdd z80 multi byte
		opcode{
			multiByte: true,
		},
		// 0xde
		opcode{},
		// 0xdf rst $18
		opcode{
			mnemonic: []string{"rst", "rst"},
			dst:      implied,
			dstR:     []string{"$18", "3"},
			noBytes:  1,
			noCycles: 11,
		},

		// 0xe0
		opcode{},
		// 0xe1 pop hl
		opcode{
			mnemonic: []string{"pop", "pop"},
			dst:      register,
			dstR:     []string{"hl", "h"},
			noBytes:  1,
			noCycles: 10,
		},
		// 0xe2 jp po,nn
		opcode{
			mnemonic: []string{"jp", "jpo"},
			dst:      condition,
			dstR:     []string{"po", ""},
			src:      immediateExtended,
			noBytes:  3,
			noCycles: 10,
		},
		// 0xe3 (sp),hl
		opcode{
			mnemonic: []string{"ex", "xthl"},
			dst:      registerIndirect,
			dstR:     []string{"sp", ""},
			src:      register,
			srcR:     []string{"hl"},
			noBytes:  1,
			noCycles: 19,
		},
		// 0xe4
		opcode{},
		// 0xe5 push hl
		opcode{
			mnemonic: []string{"push", "push"},
			dst:      register,
			dstR:     []string{"hl", "h"},
			noBytes:  1,
			noCycles: 11,
		},
		// 0xe6 and n
		opcode{
			mnemonic: []string{"and", "ani"},
			dst:      immediate,
			noBytes:  2,
			noCycles: 7,
		},
		// 0xe7 rst $20
		opcode{
			mnemonic: []string{"rst", "rst"},
			dst:      implied,
			dstR:     []string{"$20", "4"},
			noBytes:  1,
			noCycles: 11,
		},
		// 0xe8
		opcode{},
		// 0xe9 jp (hl)
		opcode{
			mnemonic: []string{"jp", "pchl"},
			dst:      registerIndirect, // FUCK YOU ZILOG
			dstR:     []string{"hl", ""},
			noBytes:  1,
			noCycles: 4,
		},
		// 0xea jp pe,nn
		opcode{
			mnemonic: []string{"jp", "jpe"},
			dst:      condition,
			dstR:     []string{"pe", ""},
			src:      immediateExtended,
			noBytes:  3,
			noCycles: 10,
		},
		// 0xeb ex de,hl
		opcode{
			mnemonic: []string{"ex", "xchg"},
			dst:      register,
			dstR:     []string{"de", ""},
			src:      register,
			srcR:     []string{"hl", ""},
			noBytes:  1,
			noCycles: 4,
		},
		// 0xec
		opcode{},
		// 0xed z80 multi byte
		opcode{
			multiByte: true,
		},
		// 0xee
		opcode{},
		// 0xef rst $28
		opcode{
			mnemonic: []string{"rst", "rst"},
			dst:      implied,
			dstR:     []string{"$28", "5"},
			noBytes:  1,
			noCycles: 11,
		},

		// 0xf0 ret p
		opcode{
			mnemonic: []string{"ret", "rp"},
			dst:      condition,
			dstR:     []string{"p", ""},
			noBytes:  1,
			noCycles: 5,
		},
		// 0xf1 pop af
		opcode{
			mnemonic: []string{"pop", "pop"},
			dst:      register,
			dstR:     []string{"af", "psw"},
			noBytes:  1,
			noCycles: 10,
		},
		// 0xf2 jp p,nn
		opcode{
			mnemonic: []string{"jp", "jp"},
			dst:      condition,
			dstR:     []string{"p", ""},
			src:      immediateExtended,
			noBytes:  3,
			noCycles: 10,
		},
		// 0xf3 di
		opcode{
			mnemonic: []string{"di", "di"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0xf4
		opcode{},
		// 0xf5
		opcode{
			mnemonic: []string{"push", "push"},
			dst:      register,
			dstR:     []string{"af", "psw"},
			noBytes:  1,
			noCycles: 11,
		},
		// 0xf6 or n
		opcode{
			mnemonic: []string{"or", "ori"},
			dst:      immediate,
			noBytes:  2,
			noCycles: 7,
		},
		// 0xf7 rst $30
		opcode{
			mnemonic: []string{"rst", "rst"},
			dst:      implied,
			dstR:     []string{"$30", "4"},
			noBytes:  1,
			noCycles: 11,
		},
		// 0xf8
		opcode{},
		// 0xf9
		opcode{},
		// 0xfa jp m,nn
		opcode{
			mnemonic: []string{"jp", "jm"},
			dst:      condition,
			dstR:     []string{"m", ""},
			src:      immediateExtended,
			noBytes:  3,
			noCycles: 10,
		},
		// 0xfb ei
		opcode{
			mnemonic: []string{"ei", "ei"},
			noBytes:  1,
			noCycles: 4,
		},
		// 0xfc
		opcode{},
		// 0xfd z80 multi byte
		opcode{
			multiByte: true,
		},
		// 0xfe
		opcode{
			mnemonic: []string{"cp", "cmp"},
			dst:      immediate,
			noBytes:  2,
			noCycles: 7,
		},

		// 0xff rst $38
		opcode{
			mnemonic: []string{"rst", "rst"},
			dst:      implied,
			dstR:     []string{"$38", "4"},
			noBytes:  1,
			noCycles: 11,
		},
	}
)
