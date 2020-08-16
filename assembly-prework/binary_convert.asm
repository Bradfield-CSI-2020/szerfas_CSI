
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

; examine first line in string
; if one, then multiply by
section .text
global binary_convert
binary_convert:
	ret
