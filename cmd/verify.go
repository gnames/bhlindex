/*
Copyright © 2022 Dmitry Mozzherin <dmozzherin@gmail.com>

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

	bhlindex "github.com/gnames/bhlindex/internal"
	"github.com/gnames/bhlindex/internal/config"
	"github.com/gnames/bhlindex/internal/io/dbio"
	"github.com/gnames/bhlindex/internal/io/verifio"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// verifyCmd represents the verify command
var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Checks if detected names do exist in biodiversity databases",
	Long: `The 'verify' command uses 'gnverifier' service to find if detected
names do exist in any of the datasets registered in 'gnverifier'.
The list of registered datasets can be found at 
'https://verifier.globalnames.org/data_sources'
The 'verify' command should be executed after 'find' command.`,
	Run: func(cmd *cobra.Command, _ []string) {
		withoutConfirm, _ := cmd.Flags().GetBool("without-confirm")
		if withoutConfirm {
			opts = append(opts, config.OptWithoutConfirm(true))
		}
		allRes, _ := cmd.Flags().GetBool("all-results")
		if allRes {
			opts = append(opts, config.OptVerifAllResults(true))
		}

		cfg := config.New(opts...)

		if !cfg.WithoutConfirm {
			fmt.Println("\nPrevious verification data will be lost.")
			fmt.Println("Do you want to proceed? (y/N)")
			var confirm string
			fmt.Scanln(&confirm)
			if confirm != "y" {
				os.Exit(0)
			}
		}

		db := dbio.New(cfg)
		bhli := bhlindex.New(cfg)
		vrf := verifio.New(cfg, db)
		err := bhli.VerifyNames(vrf)
		if err != nil {
			log.Fatal().Err(err).Msg("Name verification failed")
		}
	},
}

func init() {
	rootCmd.AddCommand(verifyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// verifyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	verifyCmd.Flags().BoolP("all-results", "a", false,
		"returns all matches instead of only best ones.")
}
