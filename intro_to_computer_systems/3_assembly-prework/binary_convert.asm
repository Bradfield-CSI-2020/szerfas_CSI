; approach

; input = input string
; result = 0
; for i = 0; i < len(input); i++ {
;    result = 2 * result + input[i]
;}

;test cases:
;1   1   result = 0, input[i] = 1 --> 1  good
;10  2   result = 1, input[i] = 0 --> 2  good
;11  3   result = 1, input[i] = 1 --> 3  good
;00  0   result = 0, input[i] = 0 --> 0  good
;01  1   result = 0, input[i] = 1 --> 1  good
;111 7   result = 3, input[i] = 1 --> 7  good
;110 6   result = 3, input[i] = 0 --> 6  good
;100 4   result = 2, input[i] = 0 --> 3  good
;...
;101 5
;011 3
;010 2
;001 1
;000 0


; in assembly pseudocode
; use jump-to-middle loop structure: initialize section, loop section, test section, done section; initialize section jumps to test section which jumps to loop section if fail
; initialize: initialize output to zero,
; loop: pull from rdi memory reference and pass that to for loop logic, increment rdi pointer
; test: if rdi memory reference is
; jump to test condition
; if test condition fails, then jump up to top of loop

                    section     .text
                    global      binary_convert
binary_convert:     mov         rax, 0                  ; initialize return value to zero
                    mov         rcx, 1                  ; prepare a 1 for future use
                    movzx       r11, byte [rdi]         ; initialize r11 for comparison comparison

loop:               mov         rdx, 0                  ; initialize/reset rdx to prepare a 0 for future use
                    cmp         r11, one_char           ; check if char equal to 1 - 49 is ascii decimal value for "1"
                    cmove       rdx, rcx                ; if so, prepare to add 1 to new result
                    sal         rax, 1                  ; multiply result by 2
                    add         rax, rdx                ; add either 1 or 0 to new resul
test:
                    inc         rdi                     ; increment character pointer
                    movzx       r11, byte [rdi]         ; pull new char for comparison
                    cmp         r11, terminal_char      ; compare to terminal char
                    jne         loop                    ; if not at null character, jump to top of the loop
done:               ret

                    section     .data
one_char:            equ         49                     ; ascii number for "1"
terminal_char:       equ         0                      ; WHY DOES PROGRAM RETURN ZERO IF THIS IS NOT SET TO NULL? (e.g., the newline char 10)?


