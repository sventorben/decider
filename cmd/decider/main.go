// Package main provides the entry point for the decider CLI.
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/sventorben/decider/internal/cli"
)

// Version information, set via ldflags at build time.
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

const defaultADRDir = "docs/adr"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "init":
		runInit(os.Args[2:])
	case "new":
		runNew(os.Args[2:])
	case "index":
		runIndex(os.Args[2:])
	case "list":
		runList(os.Args[2:])
	case "show":
		runShow(os.Args[2:])
	case "check":
		runCheck(os.Args[2:])
	case "explain":
		runExplain(os.Args[2:])
	case "version":
		printVersion()
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`decider - Git-native ADR management

Usage:
  decider <command> [options]

Commands:
  init          Initialize ADR directory structure
  new           Create a new ADR
  index         Generate/update the ADR index
  list          List ADRs with optional filters
  show          Display details of an ADR
  check         Validate ADRs or check diff applicability
  explain       Explain why ADRs apply to changed files
  version       Show version information
  help          Show this help message

Run 'decider <command> -h' for more information on a command.`)
}

func printVersion() {
	fmt.Printf("decider version %s\n", version)
	fmt.Printf("  commit: %s\n", commit)
	fmt.Printf("  built:  %s\n", date)
}

func runInit(args []string) {
	fs := flag.NewFlagSet("init", flag.ExitOnError)
	dir := fs.String("dir", defaultADRDir, "ADR directory path")
	format := fs.String("format", "text", "Output format (text|json)")
	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	outputFormat, err := cli.ParseOutputFormat(*format)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	cfg := &cli.InitConfig{
		Dir:    *dir,
		Output: cli.NewOutput(outputFormat),
	}

	if err := cli.RunInit(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func runNew(args []string) {
	fs := flag.NewFlagSet("new", flag.ExitOnError)
	dir := fs.String("dir", defaultADRDir, "ADR directory path")
	tags := fs.String("tags", "", "Comma-separated tags")
	paths := fs.String("paths", "", "Comma-separated scope paths (globs)")
	status := fs.String("status", "proposed", "Initial status")
	noIndex := fs.Bool("no-index", false, "Skip updating index")
	format := fs.String("format", "text", "Output format (text|json)")

	fs.Usage = func() {
		fmt.Println("Usage: decider new [options] <title>")
		fmt.Println()
		fmt.Println("Create a new ADR with the given title.")
		fmt.Println()
		fmt.Println("Options:")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if fs.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "error: title is required")
		fs.Usage()
		os.Exit(1)
	}

	title := strings.Join(fs.Args(), " ")

	outputFormat, err := cli.ParseOutputFormat(*format)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	var tagList []string
	if *tags != "" {
		for _, t := range strings.Split(*tags, ",") {
			tagList = append(tagList, strings.TrimSpace(t))
		}
	}

	var pathList []string
	if *paths != "" {
		for _, p := range strings.Split(*paths, ",") {
			pathList = append(pathList, strings.TrimSpace(p))
		}
	}

	cfg := &cli.NewConfig{
		Title:   title,
		Dir:     *dir,
		Tags:    tagList,
		Paths:   pathList,
		Status:  *status,
		NoIndex: *noIndex,
		Format:  outputFormat,
		Output:  cli.NewOutput(outputFormat),
	}

	if _, err := cli.RunNew(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func runIndex(args []string) {
	fs := flag.NewFlagSet("index", flag.ExitOnError)
	dir := fs.String("dir", defaultADRDir, "ADR directory path")
	check := fs.Bool("check", false, "Check if index is up-to-date (don't modify)")
	format := fs.String("format", "text", "Output format (text|json|yaml)")
	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	outputFormat, err := cli.ParseOutputFormat(*format)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	cfg := &cli.IndexConfig{
		Dir:    *dir,
		Check:  *check,
		Format: outputFormat,
		Output: cli.NewOutput(outputFormat),
	}

	_, err = cli.RunIndex(cfg)
	if err != nil {
		if *check {
			os.Exit(2) // Lint failure exit code
		}
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func runList(args []string) {
	fs := flag.NewFlagSet("list", flag.ExitOnError)
	dir := fs.String("dir", defaultADRDir, "ADR directory path")
	status := fs.String("status", "", "Filter by status")
	tag := fs.String("tag", "", "Filter by tag")
	path := fs.String("path", "", "Filter by scope path match")
	format := fs.String("format", "text", "Output format (text|json)")
	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	outputFormat, err := cli.ParseOutputFormat(*format)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	var tags []string
	if *tag != "" {
		tags = []string{*tag}
	}

	cfg := &cli.ListConfig{
		Dir:    *dir,
		Status: *status,
		Tags:   tags,
		Path:   *path,
		Format: outputFormat,
		Output: cli.NewOutput(outputFormat),
	}

	if _, err := cli.RunList(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func runShow(args []string) {
	fs := flag.NewFlagSet("show", flag.ExitOnError)
	dir := fs.String("dir", defaultADRDir, "ADR directory path")
	format := fs.String("format", "text", "Output format (text|json)")

	fs.Usage = func() {
		fmt.Println("Usage: decider show [options] <ADR-ID|number|filename>")
		fmt.Println()
		fmt.Println("Display details of a specific ADR.")
		fmt.Println()
		fmt.Println("Options:")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if fs.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "error: ADR identifier is required")
		fs.Usage()
		os.Exit(1)
	}

	outputFormat, err := cli.ParseOutputFormat(*format)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	cfg := &cli.ShowConfig{
		ID:     fs.Arg(0),
		Dir:    *dir,
		Format: outputFormat,
		Output: cli.NewOutput(outputFormat),
	}

	if _, err := cli.RunShow(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func runCheck(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "Usage: decider check <adr|diff> [options]")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Subcommands:")
		fmt.Fprintln(os.Stderr, "  adr     Validate ADR files")
		fmt.Fprintln(os.Stderr, "  diff    Find ADRs applicable to git diff")
		os.Exit(1)
	}

	subCmd := args[0]

	switch subCmd {
	case "adr":
		runCheckADR(args[1:])
	case "diff":
		runCheckDiff(args[1:])
	default:
		fmt.Fprintf(os.Stderr, "Unknown check subcommand: %s\n", subCmd)
		os.Exit(1)
	}
}

func runCheckADR(args []string) {
	fs := flag.NewFlagSet("check adr", flag.ExitOnError)
	dir := fs.String("dir", defaultADRDir, "ADR directory path")
	strict := fs.Bool("strict", false, "Treat warnings as errors (fail on missing rationale pattern)")
	format := fs.String("format", "text", "Output format (text|json)")
	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	outputFormat, err := cli.ParseOutputFormat(*format)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	cfg := &cli.CheckADRConfig{
		Dir:    *dir,
		Strict: *strict,
		Format: outputFormat,
		Output: cli.NewOutput(outputFormat),
	}

	result, err := cli.RunCheckADR(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if !result.Valid {
		os.Exit(2) // Lint failure exit code
	}
}

func runCheckDiff(args []string) {
	fs := flag.NewFlagSet("check diff", flag.ExitOnError)
	dir := fs.String("dir", defaultADRDir, "ADR directory path")
	base := fs.String("base", "", "Base ref for git diff (required)")
	format := fs.String("format", "text", "Output format (text|json)")

	fs.Usage = func() {
		fmt.Println("Usage: decider check diff --base <ref> [options]")
		fmt.Println()
		fmt.Println("Find ADRs applicable to files changed since <base>.")
		fmt.Println()
		fmt.Println("Options:")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if *base == "" {
		fmt.Fprintln(os.Stderr, "error: --base is required")
		fs.Usage()
		os.Exit(1)
	}

	outputFormat, err := cli.ParseOutputFormat(*format)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	cfg := &cli.CheckDiffConfig{
		Dir:    *dir,
		Base:   *base,
		Format: outputFormat,
		Output: cli.NewOutput(outputFormat),
	}

	if _, err := cli.RunCheckDiff(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func runExplain(args []string) {
	fs := flag.NewFlagSet("explain", flag.ExitOnError)
	dir := fs.String("dir", defaultADRDir, "ADR directory path")
	base := fs.String("base", "", "Base ref for git diff (required)")
	format := fs.String("format", "text", "Output format (text|json)")

	fs.Usage = func() {
		fmt.Println("Usage: decider explain --base <ref> [options]")
		fmt.Println()
		fmt.Println("Explain why ADRs apply to files changed since <base>.")
		fmt.Println()
		fmt.Println("Options:")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if *base == "" {
		fmt.Fprintln(os.Stderr, "error: --base is required")
		fs.Usage()
		os.Exit(1)
	}

	outputFormat, err := cli.ParseOutputFormat(*format)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	cfg := &cli.ExplainConfig{
		Dir:    *dir,
		Base:   *base,
		Format: outputFormat,
		Output: cli.NewOutput(outputFormat),
	}

	if _, err := cli.RunExplain(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
