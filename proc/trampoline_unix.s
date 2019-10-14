// +build linux darwin

#include "textflag.h"

// func trampoline(a *tramargs)
TEXT ·trampoline(SB),NOSPLIT,$8-8
    MOVQ BX, 0(SP)
    MOVQ DI, BX
    MOVQ 0(BX), AX

    // calling conventions,
    // see: https://en.wikipedia.org/wiki/X86_calling_conventions
    MOVQ 8(BX), DI
    MOVQ 16(BX), SI
    MOVQ 24(BX), DX
    MOVQ 32(BX), CX
    MOVQ 40(BX), R8
    MOVQ 48(BX), R9

    CALL AX

    LEAQ 56(BX), CX
    MOVQ AX, (CX)
    MOVQ 0(SP), BX

    RET

TEXT ·trampoline_addr(SB),NOSPLIT,$0-8
    LEAQ ·trampoline(SB), AX
    MOVQ AX, addr+0(FP)

    RET

// func realcall(a *tramargs)
TEXT ·realcall(SB),NOSPLIT,$0-8
    MOVQ a+0(FP), BX
    MOVQ 0(BX), AX
    // calling conventions,
    // see: https://en.wikipedia.org/wiki/X86_calling_conventions
    MOVQ 8(BX), DI
    MOVQ 16(BX), SI
    MOVQ 24(BX), DX
    MOVQ 32(BX), CX
    MOVQ 40(BX), R8
    MOVQ 48(BX), R9

    CALL AX

    LEAQ 56(BX), CX
    MOVQ AX, (CX)

    RET
