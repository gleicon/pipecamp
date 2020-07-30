/*
Copyright © 2019 Gleicon Moraes <gleicon@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gleicon/pipecamp/summarizer"

	"github.com/spf13/cobra"
)

var sentenceCount *int
var fileName *string

func readFromFileOrStdin(filename string) (string, error) {
	var body []string
	var fileHandler *os.File
	var err error

	if filename == "" {
		fileHandler = os.Stdin
	} else {
		fileHandler, err = os.Open(filename)
		if err != nil {
			return "", err
		}
		defer fileHandler.Close()
	}
	scanner := bufio.NewScanner(fileHandler)
	for scanner.Scan() {
		body = append(body, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}
	return (strings.Join(body, " ")), nil
}

func tldrInnerCommand(cmd *cobra.Command, args []string) {
	var body string
	body, err := readFromFileOrStdin(*fileName)
	if err != nil {
		log.Println("Error reading text: ", err)
		return
	}

	summary, err := summarizer.Summarize(body, *sentenceCount)
	if err != nil {
		log.Println("Error summarizing text: ", err)
		return
	}
	fmt.Println(summary)
}

// tldrCmd represents the tldr command
var tldrCmd = &cobra.Command{
	Use:   "tldr",
	Short: "spill out a summary from a text file",
	Long: `tldr applies the LexRank algorithm to create meaningful summaries.
	The text can be read from STDIN or using the flag -f <filename> 
	The sentence count defaults to 3 but can changed with the -s <number> flag 
	and can be used to improve the text. `,
	Run: tldrInnerCommand,
}

func init() {
	rootCmd.AddCommand(tldrCmd)
	sentenceCount = tldrCmd.Flags().IntP("sentences", "s", 3, "Change sentence count")
	fileName = tldrCmd.Flags().StringP("file", "f", "", "Filename (optional)")
}
