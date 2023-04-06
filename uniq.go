package main

import (
	"fmt"
	"flag"
	"os"
	"io"
	"bufio"
	"strings"
)

func bool2int(bool_value bool) (int) {
	if bool_value {
		return 1
	} else {
		return 0
	}
}


func main() {
	// Флаги
	var cFlag = flag.Bool("c", false, "-c - подсчитать количество встречаний строки во входных данных. Вывести это число перед строкой отделив пробелом.")
	var dFlag = flag.Bool("d", false, "-d - вывести только те строки, которые повторились во входных данных.")
	var uFlag = flag.Bool("u", false, "-u - вывести только те строки, которые не повторились во входных данных.")
	var iFlag = flag.Bool("i", false, "-i - не учитывать регистр букв.")
	var fFlag = flag.Int("f", 0, "-f num_fields - не учитывать первые num_fields полей в строке. Полем в строке является непустой набор символов отделённый пробелом.")
	var sFlag = flag.Int("s", 0, "-s num_chars - не учитывать первые num_chars символов в строке. При использовании вместе c параметром -f учитываются первые символы после num_fields полей (не учитывая пробел-разделитель после последнего поля).")
	flag.Parse()

	// Проверка на количество взаимозаменяемых флагов
	var flag_sum int = bool2int(*cFlag) + bool2int(*dFlag) + bool2int(*uFlag)
	if flag_sum > 1 {
		fmt.Printf("Usage:\nuniq [-c | -d | -u] [-i] [-f num] [-s chars] [input_file [output_file]]\n\n")
		os.Exit(1)
	}

	// Файлы ввода и вывода
	input_file_name := flag.Arg(0)
	output_file_name := flag.Arg(1)

	// Входной поток
	var in io.Reader
	if input_file_name != "" {
		f, err := os.Open(input_file_name)
		if err != nil {
			fmt.Println("error opening file: err:", err)
			os.Exit(1)
		}
		defer f.Close()

		in = f
	} else {
		in = os.Stdin
	}

	// Выходной поток
	var out io.Writer
	if output_file_name != "" {
		f, err := os.Create(output_file_name)
		if err != nil {
			fmt.Println("error opening file: err:", err)
			os.Exit(1)
		}
		defer f.Close()

		out = f
	} else {
		out = os.Stdout
	}

	// Словарь строк
	lineOrig := map[string]string{}	// обрезанная строка в нижнем регистре : первая встреча этой строки
	lineCounter := map[string]int{}	// обрезанная строка в нижнем регистре : количество

	// Чтение с входного потока и обработка строк
	buf := bufio.NewScanner(in)
	for i := 0; true; i++ {
		if !buf.Scan() {
			break
		}
		// fmt.Println(buf.Text())
		line := buf.Text()

		// Обработка строки флагами -i -f -s
		lineMod := line
		if *iFlag {
			lineMod = strings.ToLower(line)
		}
		if *fFlag > 0 {
			lineModSplit := strings.Split(lineMod, " ")
			if lineSliceLen := len(lineModSplit); lineSliceLen > *fFlag {
				lineMod = strings.Join(lineModSplit[*fFlag:], " ")
			} else if lineSliceLen > 1 {
				lineMod = strings.Join(lineModSplit[lineSliceLen:], " ")
			} else {
				lineMod = strings.Join(lineModSplit, " ")
			}
		}
		if *sFlag > 0 {
			lineModSplit := strings.Split(lineMod, "")
			if lineLen := len(lineMod); lineLen > *sFlag {
				lineMod = strings.Join(lineModSplit[*sFlag:], "")
			} else {
				lineMod = strings.Join(lineModSplit[lineLen:], "")
			}
		}

		_, lineExist := lineOrig[lineMod]
		if !lineExist {
			lineOrig[lineMod] = line
			lineCounter[lineMod] = 1
			// fmt.Println(line, 1, lineExist, !lineExist)
		} else {
			lineCounter[lineMod]++
		}
		// fmt.Fprintln(out, buf.Text())
	}

	// Вывод в выходной поток
	for key, value := range lineCounter {
		if *cFlag {
			fmt.Fprintln(out, value, lineOrig[key])		// с количеством строк
		}
		if *dFlag && (value > 1) {
			fmt.Fprintln(out, lineOrig[key])			// только строки, встретившиеся более 1 раза
		}
		if *uFlag && (value == 1) {
			fmt.Fprintln(out, lineOrig[key])			// только строки, встретившиеся ровно 1 раз
		}
		if flag_sum == 0 {
			fmt.Fprintln(out, lineOrig[key])			// без флагов, все уникальные строки
		}
	}

	if err := buf.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error reading: err:", err)
	}
}
