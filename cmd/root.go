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
	"os"

	searchengine "github.com/gleicon/pipecamp/search"
	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var cfgFile string
var datadir string
var se *searchengine.SearchEngine
var datafile string
var summarizerpath string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pipecamp",
	Short: "cli file search and summarizer",
	Long: `Command line file search and summarizer. 
  pipecamp index, search and create summaries of your local files.
  The config file sites at $HOME/pipecamp/.pipecamp.yaml, 
  along with the index files created from your files.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	defaultConfigDir, _ := os.UserHomeDir()
	defaultConfigDir = defaultConfigDir + "/.pipecamp/"
	if _, err := os.Stat(defaultConfigDir); os.IsNotExist(err) {
		fmt.Println("Creating pipecamp config and data dir " + defaultConfigDir)
		os.Mkdir(defaultConfigDir, 0755)
	}
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", defaultConfigDir+".pipecamp", "config file (default is $HOME/.pipecamp/.pipecamp.yaml)")
	rootCmd.PersistentFlags().StringVar(&datadir, "datadir", defaultConfigDir, "data dir (uses config dir, default is $HOME/.pipecamp/)")
	datapath = datadir + "index.db"
	summarizerpath = datadir + "summarizer.db"
	viper.SetConfigFile(cfgFile)
	viper.AddConfigPath(cfgFile)
	viper.SetConfigName(".pipecamp")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
