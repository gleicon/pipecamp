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
	"strings"

	"github.com/charmbracelet/glamour"
	searchengine "github.com/gleicon/pipecamp/search"
	"github.com/gleicon/pipecamp/summarizer"
	"github.com/muesli/termenv"
	"github.com/spf13/cobra"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "search <terms> - search through pre indexed documents",
	Long:  `search <term> search through documents indexed with index`,
	Run:   searchInnerCommand,
}
var renderCSV *bool

func searchInnerCommand(cmd *cobra.Command, args []string) {
	var sm *summarizer.PersistentSummarizer
	var err error

	if sm, err = summarizer.NewPersistentSummarizer(summarizerpath, 3); err != nil {
		fmt.Println(err)
		return
	}
	se = searchengine.NewSearchEngine(datapath, sm)
	terms := strings.Join(args, " ")

	// print terms, ids and summaries
	queryResults, err := se.Query(terms)
	if err != nil {
		fmt.Println(err)
		return
	}

	if *renderCSV {
		fmt.Println("document,score,summary")
		for _, result := range queryResults.Results {
			summary, err := sm.Fetch(result.ID)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Printf("%s,%f,%s,\n", result.ID, result.Score, summary)

		}
		return
	}

	termenv.ClearScreen()

	in := "# Results for " + queryResults.Query + "\n"
	in = in + "| document | score | summary |\n"
	in = in + "| --- | --- | --- |\n"

	for _, result := range queryResults.Results {
		summary, err := sm.Fetch(result.ID)
		if err != nil {
			fmt.Println(err)
			continue
		}

		in = in + fmt.Sprintf("| (%s) | %f ", result.ID, result.Score)
		if len(summary) > 40 {
			summary = summary[:40]
		}
		in = in + fmt.Sprintf("| %s |  \n", summary)

	}
	r, _ := glamour.NewTermRenderer(
		glamour.WithStandardStyle("dark"),
		glamour.WithWordWrap(120),
	)

	out, _ := r.Render(in)
	fmt.Print(out)

}

func init() {
	rootCmd.AddCommand(searchCmd)
	renderCSV = searchCmd.Flags().BoolP("csv", "c", false, "render as csv")

}
