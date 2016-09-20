package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/concourse/bosh-io-stemcell-resource/boshio"
	"github.com/concourse/bosh-io-stemcell-resource/content"
	"github.com/concourse/bosh-io-stemcell-resource/progress"
)

const routines = 10

type concourseIn struct {
	Source struct {
		Name string
	}
	Params struct {
		Tarball          bool
		PreserveFilename bool `json:"preserve_filename"`
	}
	Version struct {
		Version string
	}
}

func main() {
	rawJSON, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalln(err)
	}

	var inRequest concourseIn
	inRequest.Params.Tarball = true

	err = json.Unmarshal(rawJSON, &inRequest)
	if err != nil {
		log.Fatalln(err)
	}

	location := os.Args[1]

	client := boshio.NewClient(progress.NewBar(), content.NewRanger(routines))

	stemcells, err := client.GetStemcells(inRequest.Source.Name)
	if err != nil {
		log.Fatalln(err)
	}

	stemcell, err := client.FilterStemcells(inRequest.Version.Version, stemcells)
	if err != nil {
		log.Fatalln(err)
	}

	dataLocations := []string{"version", "sha1", "url"}

	for _, name := range dataLocations {
		fileLocation, err := os.Create(filepath.Join(location, name))
		if err != nil {
			log.Fatalln(err)
		}
		defer fileLocation.Close()

		err = client.WriteMetadata(stemcell, name, fileLocation)
		if err != nil {
			log.Fatalln(err)
		}
	}

	if inRequest.Params.Tarball {
		err = client.DownloadStemcell(stemcell, location, inRequest.Params.PreserveFilename)
		if err != nil {
			log.Fatalln(err)
		}
	}
}