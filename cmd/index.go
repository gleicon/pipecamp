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

	"github.com/gleicon/pipecamp/summarizer"

	searchengine "github.com/gleicon/pipecamp/search"
	"github.com/spf13/cobra"
)

var baseDir *string
var datapath string

// indexCmd represents the index command
var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "Creates a searchable index",
	Long:  `Given $HOME or a basedir, index all files in a way that they can be seqrched with the search command. `,
	Run:   indexInnerCommand,
}

func indexInnerCommand(cmd *cobra.Command, args []string) {
	var sm *summarizer.PersistentSummarizer
	var err error
	baseDir = rootCmd.Flags().StringP("basedir", "d", ".", "directory to index")
	if sm, err = summarizer.NewPersistentSummarizer(summarizerpath, 3); err != nil {
		fmt.Println(err)
		return
	}
	se = searchengine.NewSearchEngine(datapath, sm)
	fmt.Println("Indexing " + args[0] + " at " + datapath)
	se.AddDocuments(args[0])
}

func init() {
	rootCmd.AddCommand(indexCmd)
}
