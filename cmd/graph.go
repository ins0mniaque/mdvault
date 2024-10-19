package cmd

import (
	"fmt"
	"log"
	"mdvault/vault"
	"path/filepath"

	"github.com/spf13/cobra"
)

var graphFormat string
var graphIsolatedVertices bool
var graphLinks bool
var graphTags bool

var graphCmd = &cobra.Command{
	Use:     "graph",
	Aliases: []string{"g"},
	Short:   "Generate vault graph",
	Long:    "Generate vault graph in Markdown, Mermaid or GraphViz DOT format",
	Run: func(cmd *cobra.Command, args []string) {
		v := vault.NewVault(vaultDir)
		v.Load()

		entries := v.Entries()

		if graphFormat == "markdown" {
			println("```mermaid")
			printMermaidGraph(entries)
			println("```")
		} else if graphFormat == "mermaid" {
			printMermaidGraph(entries)
		} else if graphFormat == "dot" {
			printDotGraph(filepath.Base(v.Dir()), entries)
		} else {
			log.Fatalf("Invalid format: %s. Available formats: markdown|mermaid|dot", graphFormat)
		}
	},
}

func init() {
	graphCmd.Flags().StringVarP(&graphFormat, "format", "f", "markdown", "Output format: markdown|mermaid|dot")
	graphCmd.Flags().BoolVarP(&graphIsolatedVertices, "isolated", "i", false, "Output format: markdown|mermaid|dot")
	graphCmd.Flags().BoolVarP(&graphLinks, "links", "l", true, "Output format: markdown|mermaid|dot")
	graphCmd.Flags().BoolVarP(&graphTags, "tags", "t", true, "Output format: markdown|mermaid|dot")

	rootCmd.AddCommand(graphCmd)
}

func printMermaidGraph(entries map[string]*vault.Entry) {
	println("graph TD")
	printLinkEdges("  %s-->%s\n", entries)
	if graphTags {
		printTagEdges("  %s-->%s\n", entries)
	}
	if graphIsolatedVertices {
		printLinkVertices("  %s\n", entries)
	}
	if graphLinks {
		// TODO: Remove isolated vertices if graphIsolatedVertices is false
		printLinkVertices("  click %[1]s %[1]q %[1]q\n", entries)
	}
}

func printDotGraph(name string, entries map[string]*vault.Entry) {
	fmt.Printf("digraph %s {\n", name)
	printLinkEdges("%q -> %q\n", entries)
	if graphTags {
		printTagEdges("%q -> %q\n", entries)
	}
	if graphIsolatedVertices {
		printLinkVertices("%q\n", entries)
	}
	println("}")
}

func printLinkVertices(format string, entries map[string]*vault.Entry) {
	for vertex := range entries {
		fmt.Printf(format, vertex)
	}
}

func printLinkEdges(format string, entries map[string]*vault.Entry) {
	for source, metadata := range entries {
		for _, target := range metadata.Links {
			fmt.Printf(format, source, target)
		}
	}
}

func printTagEdges(format string, entries map[string]*vault.Entry) {
	for link, metadata := range entries {
		for _, tag := range metadata.Tags {
			fmt.Printf(format, "#"+tag, link)
		}
	}
}
