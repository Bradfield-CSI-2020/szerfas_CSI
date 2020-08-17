; approach

;get char
;reduce to char - a ascii value
;bit-set value of the char - a ascii
;check to see if equal to expected output (67108863 in decimal, or 3FFFFFF in hex, which is 26 1's in binary)

                    section     .text
                    global      pangram

pangram:
                    mov         rax, 0                  ; initialize return value to zero
                    mov         rcx, 1                  ; setting step to 1
                    movzx       r11, byte [rdi]         ; pull string to compare
                    jmp         process_first_char

process_str_loop:
                    inc         rdi
                    movzx       r11, byte [rdi]         ; pull new char for comparison

process_first_char:
                    cmp         r11, terminal_char      ; compare to terminal char
                    je          done
                    sub         r11, 65                 ; set value to 1 if A, 2 if B, and so on
                    cmp         r11, 0
                    jl          process_str_loop        ; if less than A, out of range, so repeat
                    cmp         r11, 26                 ; if less than 26 at this point, it's a cap letter ready for processing
                    jl          add_count
                    cmp         r11, 57                 ; if greater than 57, out of range
                    jg          process_str_loop
                    cmp         r11, 32                 ; if 32 <= char <= 57, then lower-case
                    jge         lower_case

add_count:
                    bts         rax, r11
                    jmp         process_str_loop

lower_case:
                    sub         r11, 32                 ; set equal to 1 if a, 2 if b, and so on
                    jmp         add_count


done:
                    mov         r10, 0
                    cmp         rax, 67108863           ; compare to 26 1's in binary
                    cmovne       rax, r10
                    ret

                    section     .data
terminal_char:      equ         0






