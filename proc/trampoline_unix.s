// +build linux darwin

#include "textflag.h"

// func trampoline(a *tramargs)
TEXT ·trampoline(SB),NOSPLIT,$0-8
    MOVQ DI, BX
    MOVQ 0(BX), R12

    // calling conventions,
    // see: https://en.wikipedia.org/wiki/X86_calling_conventions
    MOVQ 8(BX), DI
    MOVQ 16(BX), SI
    MOVQ 24(BX), DX
    MOVQ 32(BX), CX
    MOVQ 40(BX), R8
    MOVQ 48(BX), R9

    CALL R12

    LEAQ 56(BX), CX
    MOVQ AX, (CX)

    RET

TEXT ·trampoline_addr(SB),NOSPLIT,$0-8
    LEAQ ·trampoline(SB), AX
    MOVQ AX, addr+0(FP)

    RET

