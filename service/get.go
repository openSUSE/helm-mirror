package service

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"k8s.io/helm/pkg/getter"
	"k8s.io/helm/pkg/helm/environment"
	"k8s.io/helm/pkg/repo"
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
	newChartHost string
}

// NewGetService return a new instace of GetService
func NewGetService(config repo.Entry, verbose bool, ignoreErrors bool, logger *log.Logger, newChartHost string) GetServiceInterface {
	return &GetService{
		config:       config,
		verbose:      verbose,
		ignoreErrors: ignoreErrors,
		logger:       logger,
		newChartHost: newChartHost,
	}
}

//Get methods downloads the index file and the Helm charts to the working directory.
func (g *GetService) Get() error {
	chartRepo, err := repo.NewChartRepository(&g.config, getter.All(environment.EnvSettings{}))
	if err != nil {
		return err
	}

	err = chartRepo.DownloadIndexFile(g.config.Name + "/downloaded-index.yaml")
	if err != nil {
		return err
	}

	err = chartRepo.Load()
	if err != nil {
		return err
	}

	charts := chartRepo.IndexFile.Entries
	var errs []string
	for n, c := range charts {
		for _, cc := range c {
			for _, u := range cc.URLs {
				b, err := chartRepo.Client.Get(u)
				if err != nil {
					errs = append(errs, err.Error())
				}
				err = writeFile(g.config.Name+"/"+n+"-"+cc.Version+".tgz", b.Bytes(), g.logger)
				if err != nil {
					errs = append(errs, err.Error())
				}
			}
		}
	}

	err = prepareIndexFile(g.config.Name, g.config.URL, g.newChartHost, g.logger)
	if err != nil {
		errs = append(errs, err.Error())
	}

	if len(errs) > 0 && !g.ignoreErrors {
		return errors.New(strings.Join(errs, "\n"))
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

func prepareIndexFile(folder string, URL string, newChartHost string, log *log.Logger) error {
	indexContent, err := ioutil.ReadFile(folder + "/downloaded-index.yaml")
	if err != nil {
		return err
	}
	content := string(indexContent)
	if newChartHost != "" {
		content = strings.Replace(content, URL, newChartHost, -1)
	}

	err = writeFile(folder+"/index.yaml", []byte(content), log)
	if err != nil {
		return err
	}

	err = os.RemoveAll(folder + "/downloaded-index.yaml")
	if err != nil {
		return err
	}

	return nil
}
