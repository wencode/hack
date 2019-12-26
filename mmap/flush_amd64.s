#include "textflag.h"

// func flush_cache(addr, size uintptr)
// http://www.intel.com/content/dam/www/public/us/en/documents/manuals/64-ia-32-architectures-optimization-manual.pdf
// page: 7-12
TEXT ·flush_cache(SB),NOSPLIT,$0
	XORQ CX, CX
	MOVQ addr+0(FP), R9
	MOVQ size+8(FP), SI

	//MFENCE - OBSOLETE with CLFLUSH
LOOP:
	//CLFLUSH (CX)(R9)
	BYTE $0x41; BYTE $0x0f; BYTE $0xae; BYTE $0x3c; BYTE $0x09

	ADDQ $64, CX
	CMPQ CX, SI
	JL LOOP
	//MFENCE - OBSOLETE with CLFLUSH
	RET

// func flush_cache(addr, size uintptr)
// http://www.intel.com/content/dam/www/public/us/en/documents/manuals/64-ia-32-architectures-optimization-manual.pdf
// page: 7-12
TEXT ·flush_cache_opt(SB),NOSPLIT,$0
	XORQ CX, CX
	MOVQ addr+0(FP), R9
	MOVQ size+8(FP), SI

	SFENCE
LOOP:
	//CLFLUSHOPT (CX)(R9)
	BYTE $0x41; BYTE $0x66; BYTE $0x0f; BYTE $0xae; BYTE $0x3c; BYTE $0x09

	ADDQ $64, CX
	CMPQ CX, SI
	JL LOOP
	SFENCE
	RET

