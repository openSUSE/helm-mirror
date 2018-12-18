package service

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"k8s.io/helm/pkg/getter"
	"k8s.io/helm/pkg/helm/environment"
	"k8s.io/helm/pkg/repo"
)

const (
	downloadedFileName = "downloaded-index.yaml"
	indexFileName      = "index.yaml"
)

// GetServiceInterface defines a Get service
type GetServiceInterface interface {
	Get() error
}

// GetService structure definition
type GetService struct {
	config       repo.Entry
	verbose      bool
	ignoreErrors bool
	logger       *log.Logger
	newRootURL   string
}

// NewGetService return a new instace of GetService
func NewGetService(config repo.Entry, verbose bool, ignoreErrors bool, logger *log.Logger, newRootURL string) GetServiceInterface {
	return &GetService{
		config:       config,
		verbose:      verbose,
		ignoreErrors: ignoreErrors,
		logger:       logger,
		newRootURL:   newRootURL,
	}
}

//Get methods downloads the index file and the Helm charts to the working directory.
func (g *GetService) Get() error {
	chartRepo, err := repo.NewChartRepository(&g.config, getter.All(environment.EnvSettings{}))
	if err != nil {
		return err
	}

	downloadedIndexPath := path.Join(g.config.Name, downloadedFileName)
	err = chartRepo.DownloadIndexFile(downloadedIndexPath)
	if err != nil {
		return err
	}

	err = chartRepo.Load()
	if err != nil {
		return err
	}

	charts := chartRepo.IndexFile.Entries
	for n, c := range charts {
		for _, cc := range c {
			for _, u := range cc.URLs {
				b, err := chartRepo.Client.Get(u)
				if err != nil {
					if g.ignoreErrors {
						g.logger.Printf("WARNING: processing chart %s(%s) - %s", cc.Name, cc.Version, err)
						continue
					} else {
						return err
					}
				}
				chartFileName := fmt.Sprintf("%s-%s.tgz", n, cc.Version)
				chartPath := path.Join(g.config.Name, chartFileName)
				err = writeFile(chartPath, b.Bytes(), g.logger)
				if err != nil {
					if g.ignoreErrors {
						g.logger.Printf("WARNING: saving chart %s(%s) - %s", cc.Name, cc.Version, err)
						continue
					} else {
						return err
					}
				}
			}
		}
	}

	err = prepareIndexFile(g.config.Name, g.config.URL, g.newRootURL, g.logger)
	if err != nil {
		return err
	}
	return nil
}

func writeFile(name string, content []byte, log *log.Logger) error {
	err := ioutil.WriteFile(name, content, 0666)
	if err != nil {
		log.Printf("cannot write files %s: %s", name, err)
		return err
	}
	return nil
}

func prepareIndexFile(folder string, repoURL string, newRootURL string, log *log.Logger) error {
	downloadedPath := path.Join(folder, downloadedFileName)
	indexPath := path.Join(folder, indexFileName)
	if newRootURL != "" {
		indexContent, err := ioutil.ReadFile(downloadedPath)
		if err != nil {
			return err
		}
		content := bytes.Replace(indexContent, []byte(repoURL), []byte(newRootURL), -1)
		writeFile(downloadedPath, []byte(content), log)
	}
	return os.Rename(downloadedPath, indexPath)
}
