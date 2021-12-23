package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"
)

func main() {
	// 显示以下字符需要安装支持CJK扩展G区的字体，比如：SourceHanSerifTC-Regular.otf
	PrintRune('𰻞')
	PrintRune('𠀾')
	PrintRune('里')
	println("\U00030ede\U0002003E\U000091cc\u91cc")
}

func PrintRune(c rune) {
	println("Character: ", string(c))
	println("Go Style String: ", RuneToGoString(c))
	println("Java Style String: ", RuneToJavaString(c))
	println("Unicode String: ", RuneToUnicodeString(c))
	println("UTF-8 Hex: ", RuneToUtf8Hex(c))
	println("UTF-16 Hex: ", RuneToUtf16Hex(c))
	println("Code Point: ", int32(c))
	println("Code Point Hex: ", fmt.Sprintf("0x%X", c))
	println("Is Han: ", unicode.Is(unicode.Han, c))
	println()
}

func RuneToGoString(r rune) string {
	return strconv.QuoteRuneToASCII(r)
}

func RuneToJavaString(r rune) string {
	r1, r2 := utf16.EncodeRune(r)
	return fmt.Sprintf("'\\u%x\\u%x'", r1, r2)
}

func RuneToUtf16Hex(r rune) string {
	r1, r2 := utf16.EncodeRune(r)
	str := fmt.Sprintf("%X %X", r1, r2)
	return str
}

func RuneToUtf8Hex(r rune) string {
	var b []byte = make([]byte, 4)
	cout := utf8.EncodeRune(b, r)
	str := ""
	for i := 0; i < cout; i++ {
		str += fmt.Sprintf("%X", b[i])
	}
	return str
}

func RuneToUnicodeString(r rune) string {
	// return "U+" + strings.ToUpper(strconvFormatUint(uint64(r), 16))
	return fmt.Sprintf("U+%X", r)
}

func UnicodeStringToRune(s string) (rune, error) {
	v, err := strconv.ParseUint(strings.TrimPrefix(s, "U+"), 16, 32)
	if err != nil {
		return 0, err
	}
	return rune(v), nil
}
