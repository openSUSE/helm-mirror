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
	"log"
	"path"
	"path/filepath"

	"github.com/openSUSE/helm-mirror/formatter"
	"github.com/openSUSE/helm-mirror/service"

	"github.com/spf13/cobra"
)

var (
	//IgnoreErrors ignores errors in processing charts
	IgnoreErrors bool
	imagesFile   string
	output       string
	target       string
)

const imagesDesc = `Extract all the images of the Helm Chart or
the Helm Charts in the folder provided. This command dumps
the images on 'stdout' by default, for more options check
'output' flag. Example:

  - helm mirror inspect-images /tmp/helm
  - helm mirror inspect-images /tmp/helm/app.tgz

The [folder|tgzfile] has to be a full path.
`

const outputDesc = `choose an output for the list of images. Options:

- stdout: prints all images on stdout
- file: outputs all images to a file. (View file-name flag)
- json: outputs all images to a file in JSON format. (View file-name flag)
- yaml: outputs all images to a file in YAML format. (View file-name flag)
- skopeo: outputs all images to a file in YAML format to be used as input
  to Skopeo Sync. (View file-name flag)

Usage:

	- helm mirror inspect-images /tmp/helm --output stdout
	- helm mirror inspect-images /tmp/helm -o stdout
	- helm mirror inspect-images /tmp/helm -o file
	- helm mirror inspect-images /tmp/helm -o json
	- helm mirror inspect-images /tmp/helm -o yaml

`

const fileNameDesc = `set the name of the output file.

Usage:

	- helm mirror inspect-images /tmp/helm -o file --file-name images.txt
	- helm mirror inspect-images /tmp/helm -o json --file-name images.json
	- helm mirror inspect-images /tmp/helm -o yaml --file-name images.yaml
`

// inspectImagesCmd represents the images command
var inspectImagesCmd = &cobra.Command{
	Use:   "inspect-images [folder|tgzfile]",
	Short: "Extract all the images of the Helm Charts.",
	Long:  imagesDesc,
	Args:  validateInspectImagesArgs,
	RunE:  runInspectImages,
}

func init() {
	inspectImagesCmd.PersistentFlags().BoolVarP(&IgnoreErrors, "ignore-errors", "i", false, "ignores errors whiles processing charts. (Exit Code: 2)")
	inspectImagesCmd.PersistentFlags().StringVarP(&output, "output", "o", "stdout", outputDesc)
	inspectImagesCmd.PersistentFlags().StringVar(&imagesFile, "file-name", "images.out", fileNameDesc)
	rootCmd.AddCommand(inspectImagesCmd)
}

func validateInspectImagesArgs(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		logger.Print("error: requires at least one arg to execute")
		return errors.New("error: requires at least one arg")
	}
	if !path.IsAbs(args[0]) {
		logger.Printf("error: please provide a full path for [folder|tgzfile]: `%s`", args[0])
		return errors.New("error: please provide a full path for [folder|tgzfile]")
	}
	return nil
}

func resolveFormatter(output string, fileName string, l *log.Logger) (formatter.Formatter, error) {
	imagesFile, err := filepath.Abs(fileName)
	if err != nil {
		l.Print("error: geting working directory")
		return nil, err
	}
	var t formatter.Type
	switch output {
	case "file":
		t = formatter.FileType
	case "yaml":
		t = formatter.YamlType
	case "json":
		t = formatter.JSONType
	case "skopeo":
		t = formatter.SkopeoType
	default:
		t = formatter.StdoutType
	}
	return formatter.NewFormatter(t, imagesFile, l), nil
}

func runInspectImages(cmd *cobra.Command, args []string) error {
	target = args[0]
	formatter, err := resolveFormatter(output, imagesFile, logger)
	if err != nil {
		return err
	}
	imagesService := service.NewImagesService(target, Verbose, IgnoreErrors, formatter, logger)
	err = imagesService.Images()
	if err != nil {
		return err
	}
	return nil
}
