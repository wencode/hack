
TEXT ·getprocaddress(SB), 7, $0-32
	JMP	syscall·getprocaddress(SB)

TEXT ·loadlibrary(SB), 7, $0-8
	JMP	syscall·loadlibrary(SB)
