package version

import (
	"fmt"

	"github.com/spf13/cobra"

	"morpher-agent/internal/version"
)

// VersionCmd represents the version command.
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  `Print the version number of morpher-agent.`,
	Run:   runVersion,
}

func runVersion(_ *cobra.Command, _ []string) {
	fmt.Print(version.GetVersionInfo())
}
