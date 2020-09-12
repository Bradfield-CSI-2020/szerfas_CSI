	.section	__TEXT,__text,regular,pure_instructions
	.build_version macos, 10, 14	sdk_version 10, 14
	.globl	_pagecount              ## -- Begin function pagecount
	.p2align	4, 0x90
_pagecount:                             ## @pagecount
	.cfi_startproc
## %bb.0:
	pushq	%rbp
	.cfi_def_cfa_offset 16
	.cfi_offset %rbp, -16
	movq	%rsp, %rbp
	.cfi_def_cfa_register %rbp
	movq	%rdi, %rax
	xorl	%edx, %edx
	divq	%rsi
	popq	%rbp
	retq
	.cfi_endproc
                                        ## -- End function
	.section	__TEXT,__literal16,16byte_literals
	.p2align	4               ## -- Begin function main
LCPI1_0:
	.long	1127219200              ## 0x43300000
	.long	1160773632              ## 0x45300000
	.long	0                       ## 0x0
	.long	0                       ## 0x0
LCPI1_1:
	.quad	4841369599423283200     ## double 4503599627370496
	.quad	4985484787499139072     ## double 1.9342813113834067E+25
	.section	__TEXT,__literal8,8byte_literals
	.p2align	3
LCPI1_2:
	.quad	4696837146684686336     ## double 1.0E+6
LCPI1_3:
	.quad	4741671816366391296     ## double 1.0E+9
LCPI1_4:
	.quad	4711630319722168320     ## double 1.0E+7
	.section	__TEXT,__text,regular,pure_instructions
	.globl	_main
	.p2align	4, 0x90
_main:                                  ## @main
	.cfi_startproc
## %bb.0:
	pushq	%rbp
	.cfi_def_cfa_offset 16
	.cfi_offset %rbp, -16
	movq	%rsp, %rbp
	.cfi_def_cfa_register %rbp
	pushq	%r15
	pushq	%r14
	pushq	%r13
	pushq	%r12
	pushq	%rbx
	subq	$24, %rsp
	.cfi_offset %rbx, -56
	.cfi_offset %r12, -48
	.cfi_offset %r13, -40
	.cfi_offset %r14, -32
	.cfi_offset %r15, -24
	xorl	%r15d, %r15d
	movl	$8, %r13d
	movl	$1, %r12d
	movl	$10000000, %r14d        ## imm = 0x989680
	callq	_clock
	movq	%rax, -64(%rbp)         ## 8-byte Spill
	movabsq	$-6148914691236517205, %r8 ## imm = 0xAAAAAAAAAAAAAAAB
	leaq	l___const.main.msizes(%rip), %r9
	leaq	l___const.main.psizes(%rip), %rdi
	xorl	%ecx, %ecx
	xorl	%ebx, %ebx
	.p2align	4, 0x90
LBB1_1:                                 ## =>This Inner Loop Header: Depth=1
	movq	%r12, %rax
	mulq	%r8
	shlq	$2, %rdx
	andq	$-8, %rdx
	leaq	(%rdx,%rdx,2), %rax
	movq	%r13, %rsi
	subq	%rax, %rsi
	movq	%r15, %rax
	mulq	%r8
	shlq	$2, %rdx
	andq	$-8, %rdx
	leaq	(%rdx,%rdx,2), %rax
	movq	%rcx, %rdx
	subq	%rax, %rdx
	movl	(%rdx,%rdi), %eax
	addl	(%r9,%rdx), %eax
	addl	%ebx, %eax
	movl	(%rsi,%rdi), %ebx
	addl	(%r9,%rsi), %ebx
	addl	%eax, %ebx
	addq	$16, %r13
	addq	$2, %r12
	addq	$16, %rcx
	addq	$2, %r15
	addl	$-2, %r14d
	jne	LBB1_1
## %bb.2:
	callq	_clock
	movq	%rax, -56(%rbp)         ## 8-byte Spill
	movl	$8, %r14d
	movl	$1, %r13d
	movl	$10000000, %r12d        ## imm = 0x989680
	xorl	%r15d, %r15d
	callq	_clock
	leaq	l___const.main.psizes(%rip), %r9
	leaq	l___const.main.msizes(%rip), %r8
	movabsq	$-6148914691236517205, %rdi ## imm = 0xAAAAAAAAAAAAAAAB
	movq	%rax, -48(%rbp)         ## 8-byte Spill
	xorl	%ecx, %ecx
	.p2align	4, 0x90
LBB1_3:                                 ## =>This Inner Loop Header: Depth=1
	movq	%r13, %rax
	mulq	%rdi
	shlq	$2, %rdx
	andq	$-8, %rdx
	leaq	(%rdx,%rdx,2), %rax
	movq	%r14, %rsi
	subq	%rax, %rsi
	movq	%r15, %rax
	mulq	%rdi
	shlq	$2, %rdx
	andq	$-8, %rdx
	leaq	(%rdx,%rdx,2), %rax
	movq	%rcx, %rdx
	subq	%rax, %rdx
	movl	(%rdx,%r9), %eax
	addl	(%r8,%rdx), %eax
	addl	%ebx, %eax
	movl	(%rsi,%r9), %ebx
	addl	(%r8,%rsi), %ebx
	addl	%eax, %ebx
	addq	$16, %r14
	addq	$2, %r13
	addq	$16, %rcx
	addq	$2, %r15
	addl	$-2, %r12d
	jne	LBB1_3
## %bb.4:
	callq	_clock
	movq	-64(%rbp), %rcx         ## 8-byte Reload
	subq	-56(%rbp), %rcx         ## 8-byte Folded Reload
	subq	-48(%rbp), %rcx         ## 8-byte Folded Reload
	addq	%rax, %rcx
	movq	%rcx, %xmm0
	punpckldq	LCPI1_0(%rip), %xmm0 ## xmm0 = xmm0[0],mem[0],xmm0[1],mem[1]
	subpd	LCPI1_1(%rip), %xmm0
	haddpd	%xmm0, %xmm0
	divsd	LCPI1_2(%rip), %xmm0
	movsd	LCPI1_3(%rip), %xmm1    ## xmm1 = mem[0],zero
	mulsd	%xmm0, %xmm1
	divsd	LCPI1_4(%rip), %xmm1
	leaq	L_.str(%rip), %rdi
	movl	$10000000, %esi         ## imm = 0x989680
	movb	$2, %al
	callq	_printf
	movl	%ebx, %eax
	addq	$24, %rsp
	popq	%rbx
	popq	%r12
	popq	%r13
	popq	%r14
	popq	%r15
	popq	%rbp
	retq
	.cfi_endproc
                                        ## -- End function
	.section	__TEXT,__const
	.p2align	4               ## @__const.main.msizes
l___const.main.msizes:
	.quad	4294967296              ## 0x100000000
	.quad	1099511627776           ## 0x10000000000
	.quad	4503599627370496        ## 0x10000000000000

	.p2align	4               ## @__const.main.psizes
l___const.main.psizes:
	.quad	4096                    ## 0x1000
	.quad	65536                   ## 0x10000
	.quad	4294967296              ## 0x100000000

	.section	__TEXT,__cstring,cstring_literals
L_.str:                                 ## @.str
	.asciz	"%.2fs to run %d tests (%.2fns per test)\n"


.subsections_via_symbols
