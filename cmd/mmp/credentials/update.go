package credentials

import (
	"flag"
	"fmt"
	"path/filepath"

	"github.com/nlamot/nero/pkg/mmp/credentials"
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a credential for mmp in a namespace",
	Long:  `Update a credential for mmp in a namespace`,
	Run: func(cmd *cobra.Command, args []string) {
		var kubeconfig *string
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		flag.Parse()

		// use the current context in kubeconfig
		config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			panic(err.Error())
		}
		err = credentials.Update(config, "dev-mri", "test", "testje")
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Credentials updated")
		}
	},
}
