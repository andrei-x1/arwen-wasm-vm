package wasmer

const (
	OpcodeUnreachable = iota
	OpcodeNop
	OpcodeBlock
	OpcodeLoop
	OpcodeIf
	OpcodeElse
	OpcodeEnd
	OpcodeBr
	OpcodeBrIf
	OpcodeBrTable
	OpcodeReturn
	OpcodeCall
	OpcodeCallIndirect
	OpcodeDrop
	OpcodeSelect
	OpcodeGetLocal
	OpcodeSetLocal
	OpcodeTeeLocal
	OpcodeGetGlobal
	OpcodeSetGlobal
	OpcodeI32Load
	OpcodeI64Load
	OpcodeF32Load
	OpcodeF64Load
	OpcodeI32Load8S
	OpcodeI32Load8U
	OpcodeI32Load16S
	OpcodeI32Load16U
	OpcodeI64Load8S
	OpcodeI64Load8U
	OpcodeI64Load16S
	OpcodeI64Load16U
	OpcodeI64Load32S
	OpcodeI64Load32U
	OpcodeI32Store
	OpcodeI64Store
	OpcodeF32Store
	OpcodeF64Store
	OpcodeI32Store8
	OpcodeI32Store16
	OpcodeI64Store8
	OpcodeI64Store16
	OpcodeI64Store32
	OpcodeMemorySize
	OpcodeMemoryGrow
	OpcodeI32Const
	OpcodeI64Const
	OpcodeF32Const
	OpcodeF64Const
	OpcodeRefNull
	OpcodeRefIsNull
	OpcodeI32Eqz
	OpcodeI32Eq
	OpcodeI32Ne
	OpcodeI32LtS
	OpcodeI32LtU
	OpcodeI32GtS
	OpcodeI32GtU
	OpcodeI32LeS
	OpcodeI32LeU
	OpcodeI32GeS
	OpcodeI32GeU
	OpcodeI64Eqz
	OpcodeI64Eq
	OpcodeI64Ne
	OpcodeI64LtS
	OpcodeI64LtU
	OpcodeI64GtS
	OpcodeI64GtU
	OpcodeI64LeS
	OpcodeI64LeU
	OpcodeI64GeS
	OpcodeI64GeU
	OpcodeF32Eq
	OpcodeF32Ne
	OpcodeF32Lt
	OpcodeF32Gt
	OpcodeF32Le
	OpcodeF32Ge
	OpcodeF64Eq
	OpcodeF64Ne
	OpcodeF64Lt
	OpcodeF64Gt
	OpcodeF64Le
	OpcodeF64Ge
	OpcodeI32Clz
	OpcodeI32Ctz
	OpcodeI32Popcnt
	OpcodeI32Add
	OpcodeI32Sub
	OpcodeI32Mul
	OpcodeI32DivS
	OpcodeI32DivU
	OpcodeI32RemS
	OpcodeI32RemU
	OpcodeI32And
	OpcodeI32Or
	OpcodeI32Xor
	OpcodeI32Shl
	OpcodeI32ShrS
	OpcodeI32ShrU
	OpcodeI32Rotl
	OpcodeI32Rotr
	OpcodeI64Clz
	OpcodeI64Ctz
	OpcodeI64Popcnt
	OpcodeI64Add
	OpcodeI64Sub
	OpcodeI64Mul
	OpcodeI64DivS
	OpcodeI64DivU
	OpcodeI64RemS
	OpcodeI64RemU
	OpcodeI64And
	OpcodeI64Or
	OpcodeI64Xor
	OpcodeI64Shl
	OpcodeI64ShrS
	OpcodeI64ShrU
	OpcodeI64Rotl
	OpcodeI64Rotr
	OpcodeF32Abs
	OpcodeF32Neg
	OpcodeF32Ceil
	OpcodeF32Floor
	OpcodeF32Trunc
	OpcodeF32Nearest
	OpcodeF32Sqrt
	OpcodeF32Add
	OpcodeF32Sub
	OpcodeF32Mul
	OpcodeF32Div
	OpcodeF32Min
	OpcodeF32Max
	OpcodeF32Copysign
	OpcodeF64Abs
	OpcodeF64Neg
	OpcodeF64Ceil
	OpcodeF64Floor
	OpcodeF64Trunc
	OpcodeF64Nearest
	OpcodeF64Sqrt
	OpcodeF64Add
	OpcodeF64Sub
	OpcodeF64Mul
	OpcodeF64Div
	OpcodeF64Min
	OpcodeF64Max
	OpcodeF64Copysign
	OpcodeI32WrapI64
	OpcodeI32TruncSF32
	OpcodeI32TruncUF32
	OpcodeI32TruncSF64
	OpcodeI32TruncUF64
	OpcodeI64ExtendSI32
	OpcodeI64ExtendUI32
	OpcodeI64TruncSF32
	OpcodeI64TruncUF32
	OpcodeI64TruncSF64
	OpcodeI64TruncUF64
	OpcodeF32ConvertSI32
	OpcodeF32ConvertUI32
	OpcodeF32ConvertSI64
	OpcodeF32ConvertUI64
	OpcodeF32DemoteF64
	OpcodeF64ConvertSI32
	OpcodeF64ConvertUI32
	OpcodeF64ConvertSI64
	OpcodeF64ConvertUI64
	OpcodeF64PromoteF32
	OpcodeI32ReinterpretF32
	OpcodeI64ReinterpretF64
	OpcodeF32ReinterpretI32
	OpcodeF64ReinterpretI64
	OpcodeI32Extend8S
	OpcodeI32Extend16S
	OpcodeI64Extend8S
	OpcodeI64Extend16S
	OpcodeI64Extend32S
	OpcodeI32TruncSSatF32
	OpcodeI32TruncUSatF32
	OpcodeI32TruncSSatF64
	OpcodeI32TruncUSatF64
	OpcodeI64TruncSSatF32
	OpcodeI64TruncUSatF32
	OpcodeI64TruncSSatF64
	OpcodeI64TruncUSatF64
	OpcodeMemoryInit
	OpcodeDataDrop
	OpcodeMemoryCopy
	OpcodeMemoryFill
	OpcodeTableInit
	OpcodeElemDrop
	OpcodeTableCopy
	OpcodeTableGet
	OpcodeTableSet
	OpcodeTableGrow
	OpcodeTableSize
	OpcodeWake
	OpcodeI32Wait
	OpcodeI64Wait
	OpcodeFence
	OpcodeI32AtomicLoad
	OpcodeI64AtomicLoad
	OpcodeI32AtomicLoad8U
	OpcodeI32AtomicLoad16U
	OpcodeI64AtomicLoad8U
	OpcodeI64AtomicLoad16U
	OpcodeI64AtomicLoad32U
	OpcodeI32AtomicStore
	OpcodeI64AtomicStore
	OpcodeI32AtomicStore8
	OpcodeI32AtomicStore16
	OpcodeI64AtomicStore8
	OpcodeI64AtomicStore16
	OpcodeI64AtomicStore32
	OpcodeI32AtomicRmwAdd
	OpcodeI64AtomicRmwAdd
	OpcodeI32AtomicRmw8UAdd
	OpcodeI32AtomicRmw16UAdd
	OpcodeI64AtomicRmw8UAdd
	OpcodeI64AtomicRmw16UAdd
	OpcodeI64AtomicRmw32UAdd
	OpcodeI32AtomicRmwSub
	OpcodeI64AtomicRmwSub
	OpcodeI32AtomicRmw8USub
	OpcodeI32AtomicRmw16USub
	OpcodeI64AtomicRmw8USub
	OpcodeI64AtomicRmw16USub
	OpcodeI64AtomicRmw32USub
	OpcodeI32AtomicRmwAnd
	OpcodeI64AtomicRmwAnd
	OpcodeI32AtomicRmw8UAnd
	OpcodeI32AtomicRmw16UAnd
	OpcodeI64AtomicRmw8UAnd
	OpcodeI64AtomicRmw16UAnd
	OpcodeI64AtomicRmw32UAnd
	OpcodeI32AtomicRmwOr
	OpcodeI64AtomicRmwOr
	OpcodeI32AtomicRmw8UOr
	OpcodeI32AtomicRmw16UOr
	OpcodeI64AtomicRmw8UOr
	OpcodeI64AtomicRmw16UOr
	OpcodeI64AtomicRmw32UOr
	OpcodeI32AtomicRmwXor
	OpcodeI64AtomicRmwXor
	OpcodeI32AtomicRmw8UXor
	OpcodeI32AtomicRmw16UXor
	OpcodeI64AtomicRmw8UXor
	OpcodeI64AtomicRmw16UXor
	OpcodeI64AtomicRmw32UXor
	OpcodeI32AtomicRmwXchg
	OpcodeI64AtomicRmwXchg
	OpcodeI32AtomicRmw8UXchg
	OpcodeI32AtomicRmw16UXchg
	OpcodeI64AtomicRmw8UXchg
	OpcodeI64AtomicRmw16UXchg
	OpcodeI64AtomicRmw32UXchg
	OpcodeI32AtomicRmwCmpxchg
	OpcodeI64AtomicRmwCmpxchg
	OpcodeI32AtomicRmw8UCmpxchg
	OpcodeI32AtomicRmw16UCmpxchg
	OpcodeI64AtomicRmw8UCmpxchg
	OpcodeI64AtomicRmw16UCmpxchg
	OpcodeI64AtomicRmw32UCmpxchg
	OpcodeV128Load
	OpcodeV128Store
	OpcodeV128Const
	OpcodeI8x16Splat
	OpcodeI8x16ExtractLaneS
	OpcodeI8x16ExtractLaneU
	OpcodeI8x16ReplaceLane
	OpcodeI16x8Splat
	OpcodeI16x8ExtractLaneS
	OpcodeI16x8ExtractLaneU
	OpcodeI16x8ReplaceLane
	OpcodeI32x4Splat
	OpcodeI32x4ExtractLane
	OpcodeI32x4ReplaceLane
	OpcodeI64x2Splat
	OpcodeI64x2ExtractLane
	OpcodeI64x2ReplaceLane
	OpcodeF32x4Splat
	OpcodeF32x4ExtractLane
	OpcodeF32x4ReplaceLane
	OpcodeF64x2Splat
	OpcodeF64x2ExtractLane
	OpcodeF64x2ReplaceLane
	OpcodeI8x16Eq
	OpcodeI8x16Ne
	OpcodeI8x16LtS
	OpcodeI8x16LtU
	OpcodeI8x16GtS
	OpcodeI8x16GtU
	OpcodeI8x16LeS
	OpcodeI8x16LeU
	OpcodeI8x16GeS
	OpcodeI8x16GeU
	OpcodeI16x8Eq
	OpcodeI16x8Ne
	OpcodeI16x8LtS
	OpcodeI16x8LtU
	OpcodeI16x8GtS
	OpcodeI16x8GtU
	OpcodeI16x8LeS
	OpcodeI16x8LeU
	OpcodeI16x8GeS
	OpcodeI16x8GeU
	OpcodeI32x4Eq
	OpcodeI32x4Ne
	OpcodeI32x4LtS
	OpcodeI32x4LtU
	OpcodeI32x4GtS
	OpcodeI32x4GtU
	OpcodeI32x4LeS
	OpcodeI32x4LeU
	OpcodeI32x4GeS
	OpcodeI32x4GeU
	OpcodeF32x4Eq
	OpcodeF32x4Ne
	OpcodeF32x4Lt
	OpcodeF32x4Gt
	OpcodeF32x4Le
	OpcodeF32x4Ge
	OpcodeF64x2Eq
	OpcodeF64x2Ne
	OpcodeF64x2Lt
	OpcodeF64x2Gt
	OpcodeF64x2Le
	OpcodeF64x2Ge
	OpcodeV128Not
	OpcodeV128And
	OpcodeV128Or
	OpcodeV128Xor
	OpcodeV128Bitselect
	OpcodeI8x16Neg
	OpcodeI8x16AnyTrue
	OpcodeI8x16AllTrue
	OpcodeI8x16Shl
	OpcodeI8x16ShrS
	OpcodeI8x16ShrU
	OpcodeI8x16Add
	OpcodeI8x16AddSaturateS
	OpcodeI8x16AddSaturateU
	OpcodeI8x16Sub
	OpcodeI8x16SubSaturateS
	OpcodeI8x16SubSaturateU
	OpcodeI8x16Mul
	OpcodeI16x8Neg
	OpcodeI16x8AnyTrue
	OpcodeI16x8AllTrue
	OpcodeI16x8Shl
	OpcodeI16x8ShrS
	OpcodeI16x8ShrU
	OpcodeI16x8Add
	OpcodeI16x8AddSaturateS
	OpcodeI16x8AddSaturateU
	OpcodeI16x8Sub
	OpcodeI16x8SubSaturateS
	OpcodeI16x8SubSaturateU
	OpcodeI16x8Mul
	OpcodeI32x4Neg
	OpcodeI32x4AnyTrue
	OpcodeI32x4AllTrue
	OpcodeI32x4Shl
	OpcodeI32x4ShrS
	OpcodeI32x4ShrU
	OpcodeI32x4Add
	OpcodeI32x4Sub
	OpcodeI32x4Mul
	OpcodeI64x2Neg
	OpcodeI64x2AnyTrue
	OpcodeI64x2AllTrue
	OpcodeI64x2Shl
	OpcodeI64x2ShrS
	OpcodeI64x2ShrU
	OpcodeI64x2Add
	OpcodeI64x2Sub
	OpcodeF32x4Abs
	OpcodeF32x4Neg
	OpcodeF32x4Sqrt
	OpcodeF32x4Add
	OpcodeF32x4Sub
	OpcodeF32x4Mul
	OpcodeF32x4Div
	OpcodeF32x4Min
	OpcodeF32x4Max
	OpcodeF64x2Abs
	OpcodeF64x2Neg
	OpcodeF64x2Sqrt
	OpcodeF64x2Add
	OpcodeF64x2Sub
	OpcodeF64x2Mul
	OpcodeF64x2Div
	OpcodeF64x2Min
	OpcodeF64x2Max
	OpcodeI32x4TruncSF32x4Sat
	OpcodeI32x4TruncUF32x4Sat
	OpcodeI64x2TruncSF64x2Sat
	OpcodeI64x2TruncUF64x2Sat
	OpcodeF32x4ConvertSI32x4
	OpcodeF32x4ConvertUI32x4
	OpcodeF64x2ConvertSI64x2
	OpcodeF64x2ConvertUI64x2
	OpcodeV8x16Swizzle
	OpcodeV8x16Shuffle
	OpcodeI8x16LoadSplat
	OpcodeI16x8LoadSplat
	OpcodeI32x4LoadSplat
	OpcodeI64x2LoadSplat
)