package student

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func AsciiWeb(art, style string) string {
	args := make([]string, 2)
	args[0] = art
	args[1] = style

	if len(args) == 0 {
		os.Exit(0)
	}
	filename, args := GetFlag("--output=", args)
	fontStyle := "standard.txt"
	if len(args) > 1 {
		fontStyle = args[1] + ".txt"
	}

	file, err := os.Open(fontStyle)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	defer func() {
		if err = file.Close(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()

	ascii := GetASCII(file)
	buf := make([]string, 8)
	//asciiChar := 0
	str := strings.Split(args[0], "\\n")
	res := ""
	for _, a := range str {
		for _, n := range a {
			if n < 32 || n > 126 {
				fmt.Println("Error: Message is not validss")
				// os.Exit(1)
			}
			//asciiChar = int(rune(n)) - 32 //given template has ascii chars from 32 to 126 (95 chars)
			buf = AddCh(buf, ascii[n])
		}
		for i := range buf {
			res += buf[i] + "\n"
		}
		buf = make([]string, 8)
	}
	if filename == "" {
		return Prints(res)
	} else {
		ioutil.WriteFile(filename, []byte(res), 0655)
	}
	os.Exit(0)
	return Prints(res)
}

//AddCh function adds characters one by one
func AddCh(buf, new []string) []string {
	for i := range buf {
		buf[i] = buf[i] + new[i]
	}
	return buf
}

//Prints function prints ASCII message by lines
func Prints(buf string) string {
	bufLen := 0
	for range buf {
		bufLen++
	}
	str := ""
	for i := range buf {
		str = str + string(buf[i])
		// fmt.Print(string(buf[i]))
	}
	return str
}

//GetLines function separates lines given in file
func GetLines(content []byte) []string {
	lines := []string{}
	currLine := ""
	for i := 0; i < len(content); i++ {
		currLine += string(content[i])
		if content[i] == '\n' {
			lines = append(lines, currLine)
			currLine = ""
		}
	}
	return lines
}

//GetFlag function gets value of flag and deletes it from args
func GetFlag(flag string, args []string) (string, []string) {
	value := ""
	l := len(flag)
	for i := 0; i < len(args); i++ {
		if len(args[i]) > l {
			if args[i][:l] == flag {
				value = args[i][l:]
				args = append(args[:i], args[i+1:]...)
				i--
			}
		}
	}
	return value, args
}

//GetASCII fuction writes all characters to array
func GetASCII(file *os.File) map[rune][]string {
	ascii := make(map[rune][]string, 95)
	scanner := bufio.NewScanner(file)
	buf := make([]string, 8)
	charLine := 0
	asciiChar := 32
	for scanner.Scan() {
		if scanner.Text() == "" {
			charLine = 0
			buf = nil
			continue
		} else {
			buf = append(buf, scanner.Text())
			if asciiChar == 127 { //break the loop when, 96 char is read, as there are only 95 chars
				asciiChar = 0
				break
			}
			if charLine == 7 {
				ascii[rune(asciiChar)] = buf
				buf = nil
				charLine = 0
				asciiChar++
				continue
			}
			charLine++
		}
	}
	return ascii
}
