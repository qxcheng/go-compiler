    .text
.LC0:
    .string "%d\n"
printint:
	pushq   %rbp
	movq    %rsp, %rbp
	subq    $16, %rsp
	movl    %edi, -4(%rbp)
	movl    -4(%rbp), %eax
	movl    %eax, %esi
	leaq	.LC0(%rip), %rdi
	movl	$0, %eax
	call	printf@PLT
	nop
	leave
	ret
	
	.globl  main
	.type   main, @function
main:
	pushq   %rbp
	movq	%rsp, %rbp

	.comm	a,8,8
	.comm	b,8,8
	.comm	c,8,8
	movq	$3, %r8
	movq	%r8, a(%rip)
	movq	$4, %r9
	movq	%r9, b(%rip)
	movq	$4, %r10
	movq	%r10, c(%rip)
	movq	a(%rip), %r11
	movq	b(%rip), %r12
	cmpq	%r12, %r11
	jge	L0
	movq	a(%rip), %r8
	movq	%r8, %rdi
	call	printint
	movq	a(%rip), %r8
	movq	b(%rip), %r9
	addq	%r8, %r9
	movq	%r9, %rdi
	call	printint
	jmp	L1
L0:
	movq	b(%rip), %r8
	movq	%r8, %rdi
	call	printint
L1:
	movq	b(%rip), %r8
	movq	c(%rip), %r9
	cmpq	%r9, %r8
	jne	L2
	movq	c(%rip), %r8
	movq	$1, %r9
	addq	%r8, %r9
	movq	%r9, b(%rip)
	movq	b(%rip), %r8
	movq	c(%rip), %r10
	cmpq	%r10, %r8
	jle	L3
	movq	b(%rip), %r8
	movq	%r8, %rdi
	call	printint
	jmp	L4
L3:
	movq	c(%rip), %r8
	movq	%r8, %rdi
	call	printint
L4:
L2:
	movq	$10, %r8
	movq	%r8, c(%rip)
L5:
	movq	a(%rip), %r9
	movq	c(%rip), %r10
	cmpq	%r10, %r9
	jge	L6
	movq	a(%rip), %r8
	movq	%r8, %rdi
	call	printint
	movq	a(%rip), %r8
	movq	$1, %r9
	addq	%r8, %r9
	movq	%r9, a(%rip)
	jmp	L5
L6:

    movl	$0, %eax
	popq	%rbp
	ret
