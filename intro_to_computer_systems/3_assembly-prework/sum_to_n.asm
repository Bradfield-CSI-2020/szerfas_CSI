; two approaches:


; use sum_to_n = n(n+1)/2
; initialize n
; move n to return location and call x
; increment x
; multiply n by x, storing the result in x
; divide x by 2, storin gthe result in x
; return x

            section     .text
            global      sum_to_n            ; implements solution sum_to_n = n(n+1)/2
sum_to_n:   mov         rax, rdi            ; assuming n to be found in register rdi, pass to rax where it will be returned
            inc         rdi                 ; increment value in rax
            mul         rdi                 ; multiply rdi by rax (mul automatically multiplies by eax and stores in eax
            sar         rax, 1              ; divide eax by 2 using arithmetic right bit shift (couldn't find a general immediate-divide-by-constant operator)
	        ret



; use loop, either do-loop or jump-to-middle strategies
; store n
; initialize return value to 0
; initialize counter to 1
; increment add counter to
;            section     .text
;            global      sum_to_n            ; implements solution sum_to_n = n(n+1)/2
;sum_to_n:   mov         rbx, 0              ; initialize counter to 0
;            mov         rax, 0              ; initialize return value to zero
;guarded_do: cmp         rbx, rdi            ; compare counter to n; if equal we want to exit - using guarded do rather than jump-to-middle strategy to get used to potential optimization
;            je          done                ; if less, jump to complete
;loop:       inc         rbx                 ; increase the counter
;            add         rax, rbx            ; add counter to return value and store in return value
;            cmp         rbx, rdi
;            jl          loop                ; if still not at n, continue loop
;done:       ret


