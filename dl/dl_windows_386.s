
TEXT ·getprocaddress(SB), 7, $0-8
	JMP	syscall·getprocaddress(SB)

TEXT ·loadlibrary(SB), 7, $0-4
	JMP	syscall·loadlibrary(SB)
