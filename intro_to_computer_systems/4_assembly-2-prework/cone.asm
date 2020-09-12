default rel

section .text
global volume
volume:

    mulss   xmm0, xmm0      ; square the radius
    mulss   xmm0, xmm1      ; multiply by height
;    mov     esi, 3
;    cvtsi2ss xmm2, esi
;    divss    xmm0, xmm2         ; divide by three
    divss   xmm0, [const]
    movss   xmm2, [pi]
    mulss    xmm0, xmm2
 	ret


section     .data
    pi      dd 3.141592653589793238462
    const   dd 3.0