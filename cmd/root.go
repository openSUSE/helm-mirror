// Copyright © 2018 openSUSE opensuse-project@opensuse.org
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/openSUSE/helm-mirror/service"
	"github.com/spf13/cobra"
	"k8s.io/helm/pkg/repo"
)

var (
	//Verbose defines if the command is being run with verbose mode
	Verbose  bool
	folder   string
	repoURL  *url.URL
	flags    = log.Ldate | log.Lmicroseconds | log.Lshortfile
	prefix   = "helm-mirror: "
	logger   *log.Logger
	username string
	password string
	caFile   string
	certFile string
	keyFile  string
)

const rootDesc = `Mirror Helm Charts from an index file into a local folder.

For example:

helm mirror https://yourorg.com/charts /yourorg/charts

This will download the index file and the charts into
the folder indicated.

The index file is a yaml that contains a list of
charts in this format. Example:

	apiVersion: v1
	entries:
	  chart:
	  - apiVersion: 1.0.0
	    created: 2018-08-08T00:00:00.00000000Z
	    description: A Helm chart for your application
	    digest: 3aa68d6cb66c14c1fcffc6dc6d0ad8a65b90b90c10f9f04125dc6fcaf8ef1b20
	    name: chart
	    urls:
	    - https://kubernetes-charts.yourorganization.com/chart-1.0.0.tgz
	  chart2:
	  - apiVersion: 1.0.0
	    created: 2018-08-08T00:00:00.00000000Z
	    description: A Helm chart for your application
	    digest: 7ae62d60b61c14c1fcffc6dc670e72e62b91b91c10f9f04125dc67cef2ef0b21
	    name: chart
	    urls:
	    - https://kubernetes-charts.yourorganization.com/chart2-1.0.0.tgz

This will download these charts

	https://kubernetes-charts.yourorganization.com/chart-1.0.0.tgz
	https://kubernetes-charts.yourorganization.com/chart2-1.0.0.tgz

into your destination folder.`

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mirror [Repo URL] [Destination Folder]",
	Short: "Mirror Helm Charts from an index file into a local folder.",
	Long:  rootDesc,
	Args:  validateRootArgs,
	RunE:  runRoot,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	logger = log.New(os.Stdout, prefix, flags)
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")

	rootCmd.Flags().StringVar(&username, "username", "", "chart repository username")
	rootCmd.Flags().StringVar(&password, "password", "", "chart repository password")
	rootCmd.Flags().StringVar(&caFile, "ca-file", "", "verify certificates of HTTPS-enabled servers using this CA bundle")
	rootCmd.Flags().StringVar(&certFile, "cert-file", "", "identify HTTPS client using this SSL certificate file")
	rootCmd.Flags().StringVar(&keyFile, "key-file", "", "identify HTTPS client using this SSL key file")
	rootCmd.AddCommand(newVersionCmd())
}

func validateRootArgs(cmd *cobra.Command, args []string) error {
	if len(args) < 2 {
		if len(args) == 1 && args[0] == "help" {
			return nil
		}
		logger.Printf("error: requires at least two args to execute")
		return errors.New("error: requires at least two args to execute")
	}
	url, err := url.Parse(args[0])
	if err != nil {
		logger.Printf("error: not a valid URL for index file: %s", err)
		return err
	}

	if !strings.Contains(url.Scheme, "http") {
		logger.Printf("error: not a valid URL protocol: `%s`", url.Scheme)
		return errors.New("error: not a valid URL protocol")
	}
	if !path.IsAbs(args[1]) {
		logger.Printf("error: please provide a full path for destination folder: `%s`", args[1])
		return errors.New("error: please provide a full path for destination folder")
	}
	return nil
}

func runRoot(cmd *cobra.Command, args []string) error {
	repoURL, err := url.Parse(args[0])
	if err != nil {
		logger.Printf("error: not a valid URL for index file: %s", err)
		return err
	}
	folder = args[1]
	err = os.MkdirAll(folder, 0744)
	if err != nil {
		logger.Printf("error: cannot create destination folder: %s", err)
		return err
	}

	config := repo.Entry{
		Name:     folder,
		URL:      repoURL.String(),
		Username: username,
		Password: password,
		CAFile:   caFile,
		CertFile: certFile,
		KeyFile:  keyFile,
	}
	getService := service.NewGetService(config, Verbose, IgnoreErrors, logger)
	err = getService.Get()
	if err != nil {
		return err
	}
	return nil
}
