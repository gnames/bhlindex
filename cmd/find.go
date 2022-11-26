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
	"github.com/gnames/bhlindex/internal/io/finderio"
	"github.com/gnames/bhlindex/internal/io/loaderio"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// findCmd represents the find command
var findCmd = &cobra.Command{
	Use:   "find",
	Short: "Detects scientific names in BHL",
	Long: `The 'find' command traverses Biodiversity Heritage Library (BHL) files
and folders. It discovers scientific names in text pages, generates
metadata that describes locations of the names in BHL and saves results
to a database.

This command does not do verification of detected scientific names. For
verification use 'bhlindex verify' command next.`,
	Run: func(cmd *cobra.Command, _ []string) {
		withoutConfirm, _ := cmd.Flags().GetBool("without-confirm")
		if withoutConfirm {
			opts = append(opts, config.OptWithoutConfirm(true))
		}
		cfg := config.New(opts...)

		if !cfg.WithoutConfirm {
			fmt.Println("\nPreviously generated data will be lost.")
			fmt.Println("Do you want to proceed? (y/N)")
			var confirm string
			fmt.Scanln(&confirm)
			if confirm != "y" {
				os.Exit(0)
			}
		}

		db := dbio.NewWithInit(cfg)
		ldr := loaderio.New(cfg, db)
		fdr := finderio.New(cfg, db)

		bhli := bhlindex.New(cfg)

		err := bhli.FindNames(ldr, fdr)
		if err != nil {
			fmt.Fprint(os.Stderr, "\r")
			err = fmt.Errorf("FindNames %w", err)
			log.Fatal().Err(err).Msg("Name-fidning failed")
		}
	},
}

func init() {
	rootCmd.AddCommand(findCmd)
}
