/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/coverup-io/collector/client"
	"github.com/spf13/cobra"
)

var (
	filename string
	lang     string
	format   string
	repo     string
	commit   string
	key      string
	noFail   bool
	env      string
)

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload a coverage report to coverup.io",
	RunE: func(cmd *cobra.Command, args []string) error {

		if format == "" {
			switch lang {
			case "go", "golang":
				format = "golang"
			case "js", "javascript":
				format = "lcov"
			default:
				return reportFailure(errors.New("either --format or --lang required"))
			}
		}

		file, err := os.Open(filename)
		if err != nil {
			return reportFailure(fmt.Errorf("failed to open file: %w", err))
		}

		cl, err := client.New(client.WithEnvironment(env))
		if err != nil {
			return err
		}

		err = cl.UploadCoverage(key, repo, commit, format, file)
		if err != nil {
			return reportFailure(fmt.Errorf("failed to upload coverage report: %w", err))
		}

		return nil
	},
}

func reportFailure(err error) error {
	if noFail {
		log.Println(err)
		os.Exit(0)
	}
	return err
}

func init() {
	rootCmd.AddCommand(uploadCmd)
	uploadCmd.Flags().StringVarP(&filename, "file", "f", "", "Location of your coverage file")
	uploadCmd.Flags().StringVarP(&commit, "commit", "c", "", "Commit hash that you have tested")
	uploadCmd.Flags().StringVarP(&lang, "lang", "l", "", "Programming language your tests are using")
	uploadCmd.Flags().StringVarP(&format, "format", "m", "", "Format of your coverage report")
	uploadCmd.Flags().StringVarP(&repo, "repo", "r", "", "Fully qualified repo name, e.g. github.com/your-org/project")
	uploadCmd.Flags().StringVarP(&key, "key", "k", "", "Your API key (see app.coverup.io for details)")
	uploadCmd.Flags().BoolVar(&noFail, "no-fail", false, "Your API key (see app.coverup.io for details)")
	uploadCmd.Flags().StringVar(&env, "env", "production", "Override environment (dev only)")

	uploadCmd.MarkFlagRequired("file")
	uploadCmd.MarkFlagRequired("commit")
	uploadCmd.MarkFlagRequired("key")
	uploadCmd.MarkFlagRequired("repo")

	uploadCmd.Flags().MarkHidden("env")

}
