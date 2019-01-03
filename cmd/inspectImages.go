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
	"strings"

	"github.com/openSUSE/helm-mirror/formatter"
	"github.com/openSUSE/helm-mirror/service"
	"github.com/spf13/cobra"
)

var (
	output string
	target string
)

const imagesDesc = `Extract all the images of the Helm Chart or
the Helm Charts in the folder provided. This command dumps
the images on 'stdout' by default, for more options check
'output' flag. Example:

  - helm mirror inspect-images /tmp/helm
  - helm mirror inspect-images /tmp/helm/app.tgz

The [folder|tgzfile] has to be a full path.
`

const outputDesc = `choose an output for the list of images and specify
the file name, if not specified 'images.out' will be the default.
Options:

- file: outputs all images to a file
- json: outputs all images to a file in JSON format
- skopeo: outputs all images to a file in YAML format
  to be used as source file with the 'skopeo sync' command.
  For more information refer to the 'skopeo sync'
  documentation at https://github.com/SUSE/skopeo/blob/sync/docs/skopeo.1.md#skopeo-sync
- stdout: prints all images to standard output
- yaml: outputs all images to a file in YAML format

Usage:

	- helm mirror inspect-images /tmp/helm --output stdout
	- helm mirror inspect-images /tmp/helm -o stdout
	- helm mirror inspect-images /tmp/helm -o file=filename
	- helm mirror inspect-images /tmp/helm -o json=filename.json
	- helm mirror inspect-images /tmp/helm -o yaml=filename.yaml
	- helm mirror inspect-images /tmp/helm -o skopeo=filename.yaml

`

// inspectImagesCmd represents the images command
var inspectImagesCmd = &cobra.Command{
	Use:   "inspect-images [folder|tgzfile]",
	Short: "Extract all the container images listed in each chart.",
	Long:  imagesDesc,
	Args:  validateInspectImagesArgs,
	RunE:  runInspectImages,
}

func init() {

	inspectImagesCmd.PersistentFlags().StringVarP(&output, "output", "o", "stdout", outputDesc)
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

func resolveFormatter(output string, l *log.Logger) (formatter.Formatter, error) {
	a := strings.Split(output, "=")
	imagesFile := "images.out"
	if len(a) > 1 {
		imagesFile = a[1]
	}
	imagesFile, err := filepath.Abs(imagesFile)
	if err != nil {
		l.Print("error: geting working directory")
		return nil, err
	}
	var t formatter.Type
	switch a[0] {
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
	formatter, err := resolveFormatter(output, logger)
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
