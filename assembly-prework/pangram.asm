            section         .text
            global      pangram
pangram:    mov         rax, 0                  ; initialize return value to zero
;            mov         rcx, 1                  ; setting step to 1
;            movzx       r11, byte [rdi]         ; initialize r11 for comparison comparison
;loop:
;            add         rax, rcx
;test:
;            inc         rdi                     ; increment character pointer
;            movzx       r11, byte [rdi]         ; pull new char for comparison
;            cmp         r11, terminal_char      ; compare to terminal char
;            jne         loop                    ; if not at null character, jump to top of the loop

done:	    ret

;                    section     .data
;terminal_char:       equ         0                      ; WHY DOES PROGRAM RETURN ZERO IF THIS IS NOT SET TO NULL? (e.g., the newline char 10)?
;
;                    section     .rodata
;jump_table          .align



;get char
;reduce to char - a ascii value
;bitshift
;check to see if equal


