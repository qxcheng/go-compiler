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
	movq	$3, %r8
	movq	%r8, a(%rip)
	movq	$4, %r9
	movq	$5, %r10
	imulq	%r9, %r10
	movq	%r10, b(%rip)
	movq	a(%rip), %r9
	movq	b(%rip), %r11
	addq	%r9, %r11
	movq	%r11, %rdi
	call	printint
	movq	a(%rip), %r9
	movq	b(%rip), %r11
	subq	%r11, %r9
	movq	%r9, %rdi
	call	printint

    movl	$0, %eax
	popq	%rbp
	ret
