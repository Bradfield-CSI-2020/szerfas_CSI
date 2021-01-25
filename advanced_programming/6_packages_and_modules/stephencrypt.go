package stephencrypt

//const KEY = `STEPHEN`

func StephenCrypt(clearText string) string {
	var result string
	//var keyIndex int
	//for i, v := range clearText {
	for _, v := range clearText {
		//keyIndex = i % len(KEY)
		//result = result + string(v + int32(KEY[keyIndex]))
		result = result + string(v + 1)
	}
	return result
}

func StephenDecrypt(cipherText string) string {
	var result string
	for _, v := range cipherText {
		result = result + string(v - 1)
	}
	return result
}

//func main() {
//	input := "aaa"
//	fmt.Printf("input text: %s\n", input)
//	fmt.Printf("ciphered text: %s\n", StephenCrypt(input))
//	fmt.Printf("decciphereed text: %s\n", StephenDecrypt(StephenCrypt(input)))
//}