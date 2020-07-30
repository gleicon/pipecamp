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
	"fmt"

	"github.com/gleicon/pipecamp/searchengine"
	"github.com/spf13/cobra"
)

var baseDir *string
var se searchengine.SearchEngine

func indexInnerCommand(cmd *cobra.Command, args []string) {
	se.CreateOrOpenIndex(baseDir)
}

// indexCmd represents the index command
var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "Creates a searchable index",
	Long:  `Given $HOME or a basedir, index all files in a way that they can be seqrched with the search command. `,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("index called")
	},
}

func init() {
	rootCmd.AddCommand(indexCmd)
	baseDir = tldrCmd.Flags().StringP("basedir", "d", "", "directory to index down")
	se := searchengine.NewSearchEngine(indexFile)
}
