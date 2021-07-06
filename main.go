package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

func convert(r io.Reader, w io.Writer) error {
	in := bufio.NewReader(r)
	out := bufio.NewWriter(w)
	defer out.Flush()

	var collected []string
	var isCodeBlock bool
	var isHyperLink bool
	var hyperLink string
	var last rune
	for {
		c, _, err := in.ReadRune()
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}
		if c == '`' && isCodeBlock {
			isCodeBlock = false
			out.WriteRune(c)
		} else if c == '`' && !isCodeBlock {
			isCodeBlock = true
			out.WriteRune(c)
		} else if c == '[' && !isCodeBlock {
			out.WriteRune(c)
		} else if c == '(' && last == ']' && !isCodeBlock {
			hyperLink = ""
			isHyperLink = true
		} else if c == ')' && !isCodeBlock && isHyperLink {
			isHyperLink = false
			if hyperLink != "" {
				link := hyperLink
				pointer := -1
				for i, it := range collected {
					if it == link {
						pointer = i
						break
					}
				}
				if pointer == -1 {
					collected = append(collected, link)
					pointer = len(collected)
				}

				fmt.Fprintf(out, "[%d]", pointer)
			} else {
				out.WriteRune(c)
			}
		} else if !isCodeBlock && isHyperLink {
			hyperLink += fmt.Sprintf("%c", c)
		} else {
			out.WriteRune(c)
		}
		last = c
	}

	if len(collected) != 0 {
		out.WriteRune('\n')
		out.WriteRune('\n')
		for i, link := range collected {
			fmt.Fprintf(out, "[%d]:%s\n", i+1, link)
		}
	}

	return nil
}

func main() {
	err := convert(os.Stdin, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}
