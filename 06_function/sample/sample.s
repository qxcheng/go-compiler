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

	.text
	.globl	myfunc
	.type	myfunc, @function
myfunc:
	pushq	%rbp
	movq	%rsp, %rbp
	.comm	num,8,8
	movq	%rdi, num(%rip)
	movq	num(%rip), %r8
	movq	$1, %r9
	addq	%r8, %r9
	movq	%r9, %rax
	jmp	L0
L0:
    popq	%rbp
	ret

	.text
	.globl	main
	.type	main, @function
main:
	pushq	%rbp
	movq	%rsp, %rbp
	.comm	a,8,8
	.comm	b,8,8
	movq	$0, %r8
	movq	%r8, a(%rip)
	movq	$10, %r9
	movq	%r9, b(%rip)
L2:
	movq	a(%rip), %r10
	movq	b(%rip), %r11
	cmpq	%r11, %r10
	jge	L3
	movq	a(%rip), %r8
	movq	$5, %r9
	cmpq	%r9, %r8
	jge	L4
	movq	a(%rip), %r8
	movq	%r8, %rdi
	call	myfunc
	movq	%rax, %r9
	movq	%r9, %rdi
	call	printint
	jmp	L5
L4:
	movq	b(%rip), %r8
	movq	%r8, %rdi
	call	myfunc
	movq	%rax, %r9
	movq	%r9, %rdi
	call	printint
L5:
	movq	a(%rip), %r8
	movq	$1, %r9
	addq	%r8, %r9
	movq	%r9, a(%rip)
	jmp	L2
L3:
L1:
    popq	%rbp
	ret
