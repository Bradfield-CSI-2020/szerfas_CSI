; approach
; store input in a preserved register
; establish base case condition
; recurse
; take value from recursion and operate:
; release register
; return


                section     .text
                global      fib
fib:
                push        rbx             ; save callee saved register
                mov         rbx, rdi        ; put input into preserved register
                mov         rax, 1          ; in case where input is <= 1, prepare to return one
                cmp         rdi, 1
                jle         base_case       ; if input is <= 1, we're in base case and will return 1; from here on out we're in recurse case(s)
                sub         rdi, 1          ; prepare to recurse on n-1
                call        fib             ; we'll now find the return value in rax
                push        rbp             ; save callee saved register
                mov         rbp, rax        ; store in a preserved register
                lea         rdi, [rbx - 2]  ; prepare to recurse on n-1
                call        fib             ; after this line, we expect to have the result of fib(n-2) in rax and fib(n-1) in rbp
                add         rax, rbp
                pop         rbp             ; restore rbp

done:
	            pop         rbx             ; takes ?? bytes off the stack and restores to rbx
	            ret

base_case:                                  ; this can be optimized but just get to work
	            cmp         rdi, 1
	            jl          zero_base_case
	            jmp         done

zero_base_case:
                mov         rax, 0
                jmp         done