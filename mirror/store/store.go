package store

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"sync"

	"github.com/sparxfort1ano/wb-level-2/mirror/download"
	"github.com/sparxfort1ano/wb-level-2/mirror/parse"
	"golang.org/x/net/html"
)

type Store struct {
	baseHost string
	dirName  string

	maxDepth int

	visited map[string]struct{}
	sem     chan struct{}
	mu      sync.Mutex
	wg      sync.WaitGroup
}

func NewStore(
	baseHost string,
	dirName string,
	maxDepth int,
	maxConcurrency int,
	baseURL *url.URL,
) *Store {
	return &Store{
		baseHost: baseHost,
		dirName:  dirName,
		sem:      make(chan struct{}, maxConcurrency),
		maxDepth: maxDepth,
		visited: map[string]struct{}{
			baseURL.String(): {},
		},
	}
}

func (s *Store) Wait() {
	s.wg.Wait()
}

func (s *Store) Store(urls []string, currDepth int, baseURL *url.URL) error {
	if currDepth >= s.maxDepth {
		return nil
	}

	for _, urlStr := range urls {
		u, err := url.Parse(urlStr)
		if err != nil {
			return fmt.Errorf("failed to parse url: %w", err)
		}

		fullURL := baseURL.ResolveReference(u)
		if fullURL.Host != s.baseHost {
			continue
		}

		fullURLStr := fullURL.String()
		s.mu.Lock()
		if _, ok := s.visited[fullURLStr]; ok {
			s.mu.Unlock()
			continue
		}
		s.visited[fullURLStr] = struct{}{}
		s.mu.Unlock()

		s.wg.Add(1)
		s.sem <- struct{}{}

		var errFromFunc error
		go func(targetURL *url.URL, currDepth int) {
			defer s.wg.Done()
			defer func() {
				<-s.sem
			}()

			dirPath := path.Join(s.dirName, path.Dir(fullURL.Path))
			if err := os.MkdirAll(dirPath, 0750); err != nil {
				errFromFunc = fmt.Errorf("failed to create a directory: %w", err)
				return
			}

			var ext string
			if path.Ext(fullURL.Path) == "" {
				ext = "index.html"
			}
			filePath := path.Join(s.dirName, fullURL.Path, ext)
			if err := download.DownloadFile(fullURL.String(), filePath); err != nil {
				errFromFunc = err
				return
			}

			if ext := path.Ext(filePath); ext == "" || ext == ".html" {
				if ext == "" {
					filePath = path.Join(filePath, "index.html")
				}

				file, err := os.Open(filePath)
				if err != nil {
					errFromFunc = fmt.Errorf("failed to open file: %w", err)
					return
				}

				doc, err := html.Parse(file)
				if err != nil {
					errFromFunc = fmt.Errorf("failed to parse HTML file: %w", err)
					return
				}
				file.Close()

				var urls []string
				parse.ExtractHTML(doc, &urls)

				currDepth++
				if err := s.Store(urls, currDepth, fullURL); err != nil {
					errFromFunc = fmt.Errorf("failed to store a file: %w", err)
					return
				}
			}
		}(fullURL, currDepth)

		if errFromFunc != nil {
			return err
		}
	}

	return nil
}
