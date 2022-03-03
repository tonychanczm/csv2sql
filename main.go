package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

var in = flag.String("i", "", "input csv file path")
var tableName = flag.String("t", "table_name", "insert table name")
var out = flag.String("o", "", "output sql file path")

func init() {
	flag.Parse()
}

var csvSep = []byte{','}

func main() {
	var err error
	var inFile *os.File
	if *in == "" {
		inFile = os.Stdin
	} else {
		inFile, err = os.Open(*in)
		if err != nil {
			log.Printf("cannot open file %s: %v", *in, err)
			os.Exit(1)
		}
		defer inFile.Close()
	}

	var outFile *os.File
	if *out == "" {
		outFile = os.Stdout
	} else {
		outFile, err = os.Create(*out)
		if err != nil {
			log.Printf("cannot open file %s: %v", *out, err)
			os.Exit(1)
		}
		defer outFile.Close()
	}

	br := bufio.NewReader(inFile)
	for i := 0; true; i++ {
		line, err := br.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("cannot read file %s: %v", *in, err)
			os.Exit(1)
		}
		line = bytes.Trim(line, "\n\r ")
		if i == 0 {
			fields := bytes.Split(line, csvSep)
			err = writeInsertIntoFields(outFile, *tableName, fields)
			if err != nil {
				log.Printf("cannot write file %s: %v", *out, err)
				os.Exit(1)
			}
			continue
		}
		if i != 1 {
			_, err := outFile.Write([]byte{',', '\n'})
			if err != nil {
				log.Printf("cannot write file %s: %v", *out, err)
				os.Exit(1)
			}
		}
		values := bytes.Split(line, csvSep)
		err = writeValues(outFile, values)
		if err != nil {
			log.Printf("cannot write file %s: %v", *out, err)
			os.Exit(1)
		}
	}
	_, err = outFile.Write([]byte{';', '\n'})
	if err != nil {
		log.Printf("cannot write file %s: %v", *out, err)
		os.Exit(1)
	}
}

func writeInsertIntoFields(w io.Writer, tableName string, fields [][]byte) error {
	_, err := io.WriteString(w, fmt.Sprintf("INSERT INTO `%s` (", tableName))
	if err != nil {
		return err
	}
	for i, f := range fields {
		if i != 0 {
			_, err = w.Write([]byte{',', ' '})
			if err != nil {
				return err
			}
		}
		_, err = w.Write([]byte{'`'})
		if err != nil {
			return err
		}

		_, err = w.Write(trimField(f))
		if err != nil {
			return err
		}

		_, err = w.Write([]byte{'`'})
		if err != nil {
			return err
		}
	}
	_, err = io.WriteString(w, ") \nVALUES\n")
	if err != nil {
		return err
	}
	return nil
}

func writeValues(w io.Writer, lineValues [][]byte) error {
	_, err := w.Write([]byte{'('})
	if err != nil {
		return err
	}
	for i, v := range lineValues {
		if i != 0 {
			_, err := w.Write([]byte{',', ' '})
			if err != nil {
				return err
			}
		}
		_, err := w.Write([]byte{'\''})
		if err != nil {
			return err
		}

		_, err = w.Write(trimValue(v))
		if err != nil {
			return err
		}

		_, err = w.Write([]byte{'\''})
		if err != nil {
			return err
		}
	}
	_, err = w.Write([]byte{')'})
	if err != nil {
		return err
	}
	return nil
}

func trimField(f []byte) []byte {
	return bytes.Trim(f, " 	\"\uFEFF")
}

func trimValue(f []byte) []byte {
	return bytes.Trim(f, " 	\uFEFF")
}
