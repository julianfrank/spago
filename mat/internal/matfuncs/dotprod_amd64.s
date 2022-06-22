// Code generated by command: go run dotprod_asm.go -out ../../matfuncs/dotprod_amd64.s -stubs ../../matfuncs/dotprod_amd64_stubs.go -pkg matfuncs. DO NOT EDIT.

//go:build amd64 && gc && !purego

#include "textflag.h"

// func DotProdAVX32(x1 []float32, x2 []float32) float32
// Requires: AVX, FMA3, SSE
TEXT ·DotProdAVX32(SB), NOSPLIT, $0-52
	MOVQ x1_base+0(FP), AX
	MOVQ x2_base+24(FP), CX
	MOVQ x1_len+8(FP), DX
	VZEROALL

unrolledLoop14:
	CMPQ        DX, $0x00000070
	JL          unrolledLoop8
	VMOVUPS     (AX), Y2
	VMOVUPS     32(AX), Y3
	VMOVUPS     64(AX), Y4
	VMOVUPS     96(AX), Y5
	VMOVUPS     128(AX), Y6
	VMOVUPS     160(AX), Y7
	VMOVUPS     192(AX), Y8
	VMOVUPS     224(AX), Y9
	VMOVUPS     256(AX), Y10
	VMOVUPS     288(AX), Y11
	VMOVUPS     320(AX), Y12
	VMOVUPS     352(AX), Y13
	VMOVUPS     384(AX), Y14
	VMOVUPS     416(AX), Y15
	VFMADD231PS (CX), Y2, Y0
	VFMADD231PS 32(CX), Y3, Y1
	VFMADD231PS 64(CX), Y4, Y0
	VFMADD231PS 96(CX), Y5, Y1
	VFMADD231PS 128(CX), Y6, Y0
	VFMADD231PS 160(CX), Y7, Y1
	VFMADD231PS 192(CX), Y8, Y0
	VFMADD231PS 224(CX), Y9, Y1
	VFMADD231PS 256(CX), Y10, Y0
	VFMADD231PS 288(CX), Y11, Y1
	VFMADD231PS 320(CX), Y12, Y0
	VFMADD231PS 352(CX), Y13, Y1
	VFMADD231PS 384(CX), Y14, Y0
	VFMADD231PS 416(CX), Y15, Y1
	ADDQ        $0x000001c0, AX
	ADDQ        $0x000001c0, CX
	SUBQ        $0x00000070, DX
	JMP         unrolledLoop14

unrolledLoop8:
	CMPQ        DX, $0x00000040
	JL          unrolledLoop4
	VMOVUPS     (AX), Y2
	VMOVUPS     32(AX), Y3
	VMOVUPS     64(AX), Y4
	VMOVUPS     96(AX), Y5
	VMOVUPS     128(AX), Y6
	VMOVUPS     160(AX), Y7
	VMOVUPS     192(AX), Y8
	VMOVUPS     224(AX), Y9
	VFMADD231PS (CX), Y2, Y0
	VFMADD231PS 32(CX), Y3, Y1
	VFMADD231PS 64(CX), Y4, Y0
	VFMADD231PS 96(CX), Y5, Y1
	VFMADD231PS 128(CX), Y6, Y0
	VFMADD231PS 160(CX), Y7, Y1
	VFMADD231PS 192(CX), Y8, Y0
	VFMADD231PS 224(CX), Y9, Y1
	ADDQ        $0x00000100, AX
	ADDQ        $0x00000100, CX
	SUBQ        $0x00000040, DX
	JMP         unrolledLoop8

unrolledLoop4:
	CMPQ        DX, $0x00000020
	JL          unrolledLoop1
	VMOVUPS     (AX), Y2
	VMOVUPS     32(AX), Y3
	VMOVUPS     64(AX), Y4
	VMOVUPS     96(AX), Y5
	VFMADD231PS (CX), Y2, Y0
	VFMADD231PS 32(CX), Y3, Y1
	VFMADD231PS 64(CX), Y4, Y0
	VFMADD231PS 96(CX), Y5, Y1
	ADDQ        $0x00000080, AX
	ADDQ        $0x00000080, CX
	SUBQ        $0x00000020, DX
	JMP         unrolledLoop4

unrolledLoop1:
	CMPQ        DX, $0x00000008
	JL          tail
	VMOVUPS     (AX), Y2
	VFMADD231PS (CX), Y2, Y0
	ADDQ        $0x00000020, AX
	ADDQ        $0x00000020, CX
	SUBQ        $0x00000008, DX
	JMP         unrolledLoop1

tail:
	VXORPS X2, X2, X2

tailLoop:
	CMPQ        DX, $0x00000000
	JE          reduce
	VMOVSS      (AX), X3
	VFMADD231SS (CX), X3, X2
	ADDQ        $0x00000004, AX
	ADDQ        $0x00000004, CX
	DECQ        DX
	JMP         tailLoop

reduce:
	VADDPS       Y0, Y1, Y0
	VEXTRACTF128 $0x01, Y0, X1
	VADDPS       X0, X1, X0
	VADDPS       X0, X2, X0
	VHADDPS      X0, X0, X0
	VHADDPS      X0, X0, X0
	MOVSS        X0, ret+48(FP)
	RET

// func DotProdAVX64(x1 []float64, x2 []float64) float64
// Requires: AVX, FMA3, SSE2
TEXT ·DotProdAVX64(SB), NOSPLIT, $0-56
	MOVQ x1_base+0(FP), AX
	MOVQ x2_base+24(FP), CX
	MOVQ x1_len+8(FP), DX
	VZEROALL

unrolledLoop14:
	CMPQ        DX, $0x00000038
	JL          unrolledLoop8
	VMOVUPD     (AX), Y2
	VMOVUPD     32(AX), Y3
	VMOVUPD     64(AX), Y4
	VMOVUPD     96(AX), Y5
	VMOVUPD     128(AX), Y6
	VMOVUPD     160(AX), Y7
	VMOVUPD     192(AX), Y8
	VMOVUPD     224(AX), Y9
	VMOVUPD     256(AX), Y10
	VMOVUPD     288(AX), Y11
	VMOVUPD     320(AX), Y12
	VMOVUPD     352(AX), Y13
	VMOVUPD     384(AX), Y14
	VMOVUPD     416(AX), Y15
	VFMADD231PD (CX), Y2, Y0
	VFMADD231PD 32(CX), Y3, Y1
	VFMADD231PD 64(CX), Y4, Y0
	VFMADD231PD 96(CX), Y5, Y1
	VFMADD231PD 128(CX), Y6, Y0
	VFMADD231PD 160(CX), Y7, Y1
	VFMADD231PD 192(CX), Y8, Y0
	VFMADD231PD 224(CX), Y9, Y1
	VFMADD231PD 256(CX), Y10, Y0
	VFMADD231PD 288(CX), Y11, Y1
	VFMADD231PD 320(CX), Y12, Y0
	VFMADD231PD 352(CX), Y13, Y1
	VFMADD231PD 384(CX), Y14, Y0
	VFMADD231PD 416(CX), Y15, Y1
	ADDQ        $0x000001c0, AX
	ADDQ        $0x000001c0, CX
	SUBQ        $0x00000038, DX
	JMP         unrolledLoop14

unrolledLoop8:
	CMPQ        DX, $0x00000020
	JL          unrolledLoop4
	VMOVUPD     (AX), Y2
	VMOVUPD     32(AX), Y3
	VMOVUPD     64(AX), Y4
	VMOVUPD     96(AX), Y5
	VMOVUPD     128(AX), Y6
	VMOVUPD     160(AX), Y7
	VMOVUPD     192(AX), Y8
	VMOVUPD     224(AX), Y9
	VFMADD231PD (CX), Y2, Y0
	VFMADD231PD 32(CX), Y3, Y1
	VFMADD231PD 64(CX), Y4, Y0
	VFMADD231PD 96(CX), Y5, Y1
	VFMADD231PD 128(CX), Y6, Y0
	VFMADD231PD 160(CX), Y7, Y1
	VFMADD231PD 192(CX), Y8, Y0
	VFMADD231PD 224(CX), Y9, Y1
	ADDQ        $0x00000100, AX
	ADDQ        $0x00000100, CX
	SUBQ        $0x00000020, DX
	JMP         unrolledLoop8

unrolledLoop4:
	CMPQ        DX, $0x00000010
	JL          unrolledLoop1
	VMOVUPD     (AX), Y2
	VMOVUPD     32(AX), Y3
	VMOVUPD     64(AX), Y4
	VMOVUPD     96(AX), Y5
	VFMADD231PD (CX), Y2, Y0
	VFMADD231PD 32(CX), Y3, Y1
	VFMADD231PD 64(CX), Y4, Y0
	VFMADD231PD 96(CX), Y5, Y1
	ADDQ        $0x00000080, AX
	ADDQ        $0x00000080, CX
	SUBQ        $0x00000010, DX
	JMP         unrolledLoop4

unrolledLoop1:
	CMPQ        DX, $0x00000004
	JL          tail
	VMOVUPD     (AX), Y2
	VFMADD231PD (CX), Y2, Y0
	ADDQ        $0x00000020, AX
	ADDQ        $0x00000020, CX
	SUBQ        $0x00000004, DX
	JMP         unrolledLoop1

tail:
	VXORPD X2, X2, X2

tailLoop:
	CMPQ        DX, $0x00000000
	JE          reduce
	VMOVSD      (AX), X3
	VFMADD231SD (CX), X3, X2
	ADDQ        $0x00000008, AX
	ADDQ        $0x00000008, CX
	DECQ        DX
	JMP         tailLoop

reduce:
	VADDPD       Y0, Y1, Y0
	VEXTRACTF128 $0x01, Y0, X1
	VADDPD       X0, X1, X0
	VADDPD       X0, X2, X0
	VHADDPD      X0, X0, X0
	MOVSD        X0, ret+48(FP)
	RET

// func DotProdSSE32(x1 []float32, x2 []float32) float32
// Requires: SSE, SSE3
TEXT ·DotProdSSE32(SB), NOSPLIT, $0-52
	MOVQ  x1_base+0(FP), AX
	MOVQ  x2_base+24(FP), CX
	MOVQ  x1_len+8(FP), DX
	XORPS X0, X0
	XORPS X1, X1
	CMPQ  DX, $0x00000000
	JE    reduce
	MOVQ  CX, BX
	ANDQ  $0x0000000f, BX
	JZ    unrolledLoops
	XORQ  $0x0000000f, BX
	INCQ  BX
	SHRQ  $0x02, BX

alignmentLoop:
	MOVSS (AX), X2
	MULSS (CX), X2
	ADDSS X2, X0
	ADDQ  $0x00000004, AX
	ADDQ  $0x00000004, CX
	DECQ  DX
	JZ    reduce
	DECQ  BX
	JNZ   alignmentLoop

unrolledLoops:
unrolledLoop14:
	CMPQ   DX, $0x00000038
	JL     unrolledLoop8
	MOVUPS (AX), X2
	MOVUPS 16(AX), X3
	MOVUPS 32(AX), X4
	MOVUPS 48(AX), X5
	MOVUPS 64(AX), X6
	MOVUPS 80(AX), X7
	MOVUPS 96(AX), X8
	MOVUPS 112(AX), X9
	MOVUPS 128(AX), X10
	MOVUPS 144(AX), X11
	MOVUPS 160(AX), X12
	MOVUPS 176(AX), X13
	MOVUPS 192(AX), X14
	MOVUPS 208(AX), X15
	MULPS  (CX), X2
	MULPS  16(CX), X3
	MULPS  32(CX), X4
	MULPS  48(CX), X5
	MULPS  64(CX), X6
	MULPS  80(CX), X7
	MULPS  96(CX), X8
	MULPS  112(CX), X9
	MULPS  128(CX), X10
	MULPS  144(CX), X11
	MULPS  160(CX), X12
	MULPS  176(CX), X13
	MULPS  192(CX), X14
	MULPS  208(CX), X15
	ADDPS  X2, X0
	ADDPS  X3, X1
	ADDPS  X4, X0
	ADDPS  X5, X1
	ADDPS  X6, X0
	ADDPS  X7, X1
	ADDPS  X8, X0
	ADDPS  X9, X1
	ADDPS  X10, X0
	ADDPS  X11, X1
	ADDPS  X12, X0
	ADDPS  X13, X1
	ADDPS  X14, X0
	ADDPS  X15, X1
	ADDQ   $0x000000e0, AX
	ADDQ   $0x000000e0, CX
	SUBQ   $0x00000038, DX
	JMP    unrolledLoop14

unrolledLoop8:
	CMPQ   DX, $0x00000020
	JL     unrolledLoop4
	MOVUPS (AX), X2
	MOVUPS 16(AX), X3
	MOVUPS 32(AX), X4
	MOVUPS 48(AX), X5
	MOVUPS 64(AX), X6
	MOVUPS 80(AX), X7
	MOVUPS 96(AX), X8
	MOVUPS 112(AX), X9
	MULPS  (CX), X2
	MULPS  16(CX), X3
	MULPS  32(CX), X4
	MULPS  48(CX), X5
	MULPS  64(CX), X6
	MULPS  80(CX), X7
	MULPS  96(CX), X8
	MULPS  112(CX), X9
	ADDPS  X2, X0
	ADDPS  X3, X1
	ADDPS  X4, X0
	ADDPS  X5, X1
	ADDPS  X6, X0
	ADDPS  X7, X1
	ADDPS  X8, X0
	ADDPS  X9, X1
	ADDQ   $0x00000080, AX
	ADDQ   $0x00000080, CX
	SUBQ   $0x00000020, DX
	JMP    unrolledLoop8

unrolledLoop4:
	CMPQ   DX, $0x00000010
	JL     unrolledLoop1
	MOVUPS (AX), X2
	MOVUPS 16(AX), X3
	MOVUPS 32(AX), X4
	MOVUPS 48(AX), X5
	MULPS  (CX), X2
	MULPS  16(CX), X3
	MULPS  32(CX), X4
	MULPS  48(CX), X5
	ADDPS  X2, X0
	ADDPS  X3, X1
	ADDPS  X4, X0
	ADDPS  X5, X1
	ADDQ   $0x00000040, AX
	ADDQ   $0x00000040, CX
	SUBQ   $0x00000010, DX
	JMP    unrolledLoop4

unrolledLoop1:
	CMPQ   DX, $0x00000004
	JL     tailLoop
	MOVUPS (AX), X2
	MULPS  (CX), X2
	ADDPS  X2, X0
	ADDQ   $0x00000010, AX
	ADDQ   $0x00000010, CX
	SUBQ   $0x00000004, DX
	JMP    unrolledLoop1

tailLoop:
	CMPQ  DX, $0x00000000
	JE    reduce
	MOVSS (AX), X2
	MULSS (CX), X2
	ADDSS X2, X0
	ADDQ  $0x00000004, AX
	ADDQ  $0x00000004, CX
	DECQ  DX
	JMP   tailLoop

reduce:
	ADDPS  X1, X0
	HADDPS X0, X0
	HADDPS X0, X0
	MOVSS  X0, ret+48(FP)
	RET

// func DotProdSSE64(x1 []float64, x2 []float64) float64
// Requires: SSE2, SSE3
TEXT ·DotProdSSE64(SB), NOSPLIT, $0-56
	MOVQ  x1_base+0(FP), AX
	MOVQ  x2_base+24(FP), CX
	MOVQ  x1_len+8(FP), DX
	XORPD X0, X0
	XORPD X1, X1
	CMPQ  DX, $0x00000000
	JE    reduce
	MOVQ  CX, BX
	ANDQ  $0x0000000f, BX
	JZ    unrolledLoops
	MOVSD (AX), X2
	MULSD (CX), X2
	ADDSD X2, X0
	ADDQ  $0x00000008, AX
	ADDQ  $0x00000008, CX
	DECQ  DX

unrolledLoops:
unrolledLoop14:
	CMPQ   DX, $0x0000001c
	JL     unrolledLoop8
	MOVUPD (AX), X2
	MOVUPD 16(AX), X3
	MOVUPD 32(AX), X4
	MOVUPD 48(AX), X5
	MOVUPD 64(AX), X6
	MOVUPD 80(AX), X7
	MOVUPD 96(AX), X8
	MOVUPD 112(AX), X9
	MOVUPD 128(AX), X10
	MOVUPD 144(AX), X11
	MOVUPD 160(AX), X12
	MOVUPD 176(AX), X13
	MOVUPD 192(AX), X14
	MOVUPD 208(AX), X15
	MULPD  (CX), X2
	MULPD  16(CX), X3
	MULPD  32(CX), X4
	MULPD  48(CX), X5
	MULPD  64(CX), X6
	MULPD  80(CX), X7
	MULPD  96(CX), X8
	MULPD  112(CX), X9
	MULPD  128(CX), X10
	MULPD  144(CX), X11
	MULPD  160(CX), X12
	MULPD  176(CX), X13
	MULPD  192(CX), X14
	MULPD  208(CX), X15
	ADDPD  X2, X0
	ADDPD  X3, X1
	ADDPD  X4, X0
	ADDPD  X5, X1
	ADDPD  X6, X0
	ADDPD  X7, X1
	ADDPD  X8, X0
	ADDPD  X9, X1
	ADDPD  X10, X0
	ADDPD  X11, X1
	ADDPD  X12, X0
	ADDPD  X13, X1
	ADDPD  X14, X0
	ADDPD  X15, X1
	ADDQ   $0x000000e0, AX
	ADDQ   $0x000000e0, CX
	SUBQ   $0x0000001c, DX
	JMP    unrolledLoop14

unrolledLoop8:
	CMPQ   DX, $0x00000010
	JL     unrolledLoop4
	MOVUPD (AX), X2
	MOVUPD 16(AX), X3
	MOVUPD 32(AX), X4
	MOVUPD 48(AX), X5
	MOVUPD 64(AX), X6
	MOVUPD 80(AX), X7
	MOVUPD 96(AX), X8
	MOVUPD 112(AX), X9
	MULPD  (CX), X2
	MULPD  16(CX), X3
	MULPD  32(CX), X4
	MULPD  48(CX), X5
	MULPD  64(CX), X6
	MULPD  80(CX), X7
	MULPD  96(CX), X8
	MULPD  112(CX), X9
	ADDPD  X2, X0
	ADDPD  X3, X1
	ADDPD  X4, X0
	ADDPD  X5, X1
	ADDPD  X6, X0
	ADDPD  X7, X1
	ADDPD  X8, X0
	ADDPD  X9, X1
	ADDQ   $0x00000080, AX
	ADDQ   $0x00000080, CX
	SUBQ   $0x00000010, DX
	JMP    unrolledLoop8

unrolledLoop4:
	CMPQ   DX, $0x00000008
	JL     unrolledLoop1
	MOVUPD (AX), X2
	MOVUPD 16(AX), X3
	MOVUPD 32(AX), X4
	MOVUPD 48(AX), X5
	MULPD  (CX), X2
	MULPD  16(CX), X3
	MULPD  32(CX), X4
	MULPD  48(CX), X5
	ADDPD  X2, X0
	ADDPD  X3, X1
	ADDPD  X4, X0
	ADDPD  X5, X1
	ADDQ   $0x00000040, AX
	ADDQ   $0x00000040, CX
	SUBQ   $0x00000008, DX
	JMP    unrolledLoop4

unrolledLoop1:
	CMPQ   DX, $0x00000002
	JL     tailLoop
	MOVUPD (AX), X2
	MULPD  (CX), X2
	ADDPD  X2, X0
	ADDQ   $0x00000010, AX
	ADDQ   $0x00000010, CX
	SUBQ   $0x00000002, DX
	JMP    unrolledLoop1

tailLoop:
	CMPQ  DX, $0x00000000
	JE    reduce
	MOVSD (AX), X2
	MULSD (CX), X2
	ADDSD X2, X0
	ADDQ  $0x00000008, AX
	ADDQ  $0x00000008, CX
	DECQ  DX
	JMP   tailLoop

reduce:
	ADDPD  X1, X0
	HADDPD X0, X0
	MOVSD  X0, ret+48(FP)
	RET