package cmd

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var vaultDir string

var rootCmd = &cobra.Command{
	Use:   "mdvault",
	Short: "mdvault is a markdown knowledge base command-line tool",
	Long:  "mdvault is a markdown knowledge base command-line tool",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initVault)

	rootCmd.PersistentFlags().StringVar(&vaultDir, "vault", "", "vault directory")
}

func initVault() {
	if vaultDir == "" {
		wd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		vaultDir = findVaultDir(wd, ".mdvault")
		if vaultDir == "" {
			vaultDir = findVaultDir(wd, ".obsidian")
		}
		if vaultDir == "" {
			vaultDir = findVaultDir(wd, ".git")
		}
		if vaultDir == "" {
			vaultDir = wd
		}
	}
}

func isVaultDir(path string, configDirName string) bool {
	configDir := filepath.Join(path, configDirName)
	info, err := os.Stat(configDir)
	return err == nil && info.IsDir()
}

func findVaultDir(path string, configDirName string) string {
	for {
		if isVaultDir(path, configDirName) {
			return path
		}

		parentPath := filepath.Dir(path)
		println(parentPath)
		if parentPath == path {
			break
		}

		path = parentPath
	}

	return ""
}
