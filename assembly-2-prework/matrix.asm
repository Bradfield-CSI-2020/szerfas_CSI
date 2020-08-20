section .text
global index
index:
	; rdi: matrix
	; rsi: rows
	; rdx: cols
	; rcx: rindex
	; r8: cindex
    ; increment matrix pointer down X rows * size of each row, then number of desired colums, and dereference.
    ; size of each row is # col * size of col (AKA size of data type)

    imul rcx, rdx            ; each row contains all columns
    lea rax, [rdi + rcx*4]   ; move down the desired number of rows (scale 4 = size of int). use lea instead of move so that we don't de-reference twice
    mov rax, [rax + r8*4]    ; move down desired number of columns and derence
	ret
