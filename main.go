package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	flagOutput = pflag.StringP("output", "o", ".github/CODEOWNERS", "output file. If \"-\", write to stdout")
	flagIgnore = pflag.StringArrayP("ignore", "i", []string{"/vendor/"}, "list of filepath regexes to ignore relative to the input directory")
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func parseArgs(args []string) (string, error) {
	if len(args) < 1 {
		return os.Getwd()
	}
	return args[0], nil
}

var rootCmd = &cobra.Command{
	Use:   "codeowners [flags] [path]",
	Short: "codeowners generates a GitHub CODEOWNERS file from multiple CODEOWNERS files throughout the repo.",
	Args:  cobra.MaximumNArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		root, err := parseArgs(args)
		if err != nil {
			log.Fatal(fmt.Errorf("error while parsing root dir: %w", err))
		}

		rewrittenCodeownerRules, err := RewriteCodeownersRules(root, *flagIgnore)
		if err != nil {
			log.Fatal(fmt.Errorf("error while rewriting codeowner rules in %s: %w", root, err))
		}

		if len(rewrittenCodeownerRules) == 0 {
			log.Fatal(fmt.Errorf("no CODEOWNER rules found in %s", root))
		}

		generatedCodeownersFile := GenerateCodeownersFile(rewrittenCodeownerRules)

		w, err := outWriter()
		if err != nil {
			log.Fatal(fmt.Errorf("error opening output file: %w", err))
		}
		fmt.Fprintln(w, generatedCodeownersFile)
	},
}

func outWriter() (io.Writer, error) {
	out := *flagOutput
	if out == "-" {
		return os.Stdout, nil
	}
	return os.Create(out)
}
