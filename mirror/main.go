package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path"

	"github.com/sparxfort1ano/wb-level-2/mirror/download"
	"github.com/sparxfort1ano/wb-level-2/mirror/parse"
	"github.com/sparxfort1ano/wb-level-2/mirror/store"
	"golang.org/x/net/html"
)

func main() {
	const (
		baseDir    = "testMirrors"
		subDir     = "test2"
		baseURLStr = "https://books.toscrape.com/"
	)

	dirName := path.Join(baseDir, subDir)
	if err := stat(dirName); err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll(dirName, 0750); err != nil {
		log.Fatal(fmt.Errorf("failed to create a directory: %w", err))
	}

	baseURL, err := url.Parse(baseURLStr)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to parse url: %w", err))
	}
	indexHTML := path.Join(dirName, "index.html")

	if err := download.DownloadFile(baseURL.String(), indexHTML); err != nil {
		log.Fatal(err)
	}

	baseFile, err := os.Open(indexHTML)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to open index.html file: %w", err))
	}

	doc, err := html.Parse(baseFile)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to parse HTML file: %w", err))
	}
	baseFile.Close()

	var urls []string
	parse.ExtractHTML(doc, &urls)

	store := store.NewStore(baseURL.Host, dirName, 1, 10, baseURL)
	if err := store.Store(urls, 0, baseURL); err != nil {
		log.Fatal(err)
	}

	store.Wait()
}

func stat(dirName string) error {
	_, err := os.Stat(dirName)
	switch err {
	case nil:
		return fmt.Errorf("folder already exists")
	default:
		return nil
	}
}
