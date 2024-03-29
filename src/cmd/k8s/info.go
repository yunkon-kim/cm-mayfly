package k8s

import (
	"fmt"

	"github.com/cm-mayfly/cm-mayfly/src/common"
	"github.com/spf13/cobra"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get information of Cloud-Migrator System",
	Long:  `Get information of Cloud-Migrator System. Information about containers and container images`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("\n[Get info for Cloud-Migrator runtimes]")
		fmt.Println()

		if common.K8sFilePath == "" {
			fmt.Println("file is required")
		} else {
			var cmdStr string

			fmt.Println("[v]Status of Cloud-Migrator Helm release")
			cmdStr = fmt.Sprintf("helm status --namespace %s %s", common.CMK8sNamespace, common.CMHelmReleaseName)
			common.SysCall(cmdStr)
			fmt.Println()
			fmt.Println("[v]Status of Cloud-Migrator pods")
			cmdStr = fmt.Sprintf("kubectl get pods -n %s", common.CMK8sNamespace)
			common.SysCall(cmdStr)
			fmt.Println()
			fmt.Println("[v]Status of Cloud-Migrator container images")
			cmdStr = `kubectl get pods -n ` + common.CMK8sNamespace + ` -o jsonpath="{..image}" |\
				tr -s '[[:space:]]' '\n' |\
				sort |\
				uniq`
			common.SysCall(cmdStr)

		}
	},
}

func init() {
	k8sCmd.AddCommand(infoCmd)

	pf := infoCmd.PersistentFlags()
	pf.StringVarP(&common.K8sFilePath, "file", "f", common.DefaultKubernetesConfig, "User-defined configuration file")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// infoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// infoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
