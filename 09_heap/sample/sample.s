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

	.data
	.globl	c
c:	.quad	0
	.data
	.globl	d
d:	.quad	0

	.text
	.globl	myfunc
	.type	myfunc, @function
myfunc:
	pushq	%rbp
	movq	%rsp, %rbp
	addq	$-16,%rsp
	movq	%rdi, -8(%rbp)
	movq	-8(%rbp), %r8
	movq	(%r8), %r8
	movq	$1, %r9
	addq	%r8, %r9
	movq	-8(%rbp), %r8
	movq	%r9, (%r8)
	movq	-8(%rbp), %r10
	movq	(%r10), %r10
	movq	%r10, %rax
	jmp	L0
L0:
	addq	$16,%rsp
	popq	%rbp
	ret

	.text
	.globl	main
	.type	main, @function
main:
	pushq	%rbp
	movq	%rsp, %rbp
	addq	$-16,%rsp
	movq	$10, %r8
	movq	%r8, -8(%rbp)
	movq	$0, %r9
	movq	%r9, d(%rip)
	movq	d(%rip), %r10
	leaq	d(%rip), %r11
	movq	%r11, c(%rip)
	movq	c(%rip), %r12
	movq	%r12, %rdi
	call	printint
L2:
	movq	d(%rip), %r12
	movq	-8(%rbp), %r13
	cmpq	%r13, %r12
	jge	L3
	movq	d(%rip), %r8
	movq	$5, %r9
	cmpq	%r9, %r8
	jge	L4
	movq	c(%rip), %r8
	movq	%r8, %rdi
	call	myfunc
	movq	%rax, %r9
	movq	%r9, %rdi
	call	printint
	jmp	L5
L4:
	movq	$0, %r8
	movq	%r8, %rdi
	call	printint
	movq	d(%rip), %r8
	movq	$1, %r9
	addq	%r8, %r9
	movq	%r9, d(%rip)
L5:
	jmp	L2
L3:
	movq	c(%rip), %r8
	movq	(%r8), %r8
	movq	%r8, %rdi
	call	printint
L1:
    mov $8, %rdi
    call malloc
    movq	$0, %r8
    movq	%rax, %r8
    movq	$1000, (%r8)
    movq	(%r8), %rdi
    call	printint

    int $0x80

	addq	$16,%rsp
	popq	%rbp
	ret
