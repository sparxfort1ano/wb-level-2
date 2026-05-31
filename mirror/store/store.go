// Package store is responsible for mirroring web-sites and storing the payload to appropriate directories.
package store

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path"
	"sync"

	"github.com/sparxfort1ano/wb-level-2/mirror/download"
	"github.com/sparxfort1ano/wb-level-2/mirror/parse"
	"golang.org/x/net/html"
)

// Store holds the configuration and state for the mirroring process.
type Store struct {
	baseHost string
	dirName  string

	maxDepth int

	visited map[string]struct{}
	sem     chan struct{}
	mu      sync.Mutex
	wg      sync.WaitGroup

	allErrs error
	errsMu  sync.Mutex
}

// NewStore creates a new instance of Store.
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

func (s *Store) Errors() error {
	return s.allErrs
}

// Store processes a slice of URLs, downloading valid assets and recursively
// traversing internal HTML pages up to the configured maximum depth.
// It executes downloads concurrently using a semaphore to limit active connections.
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

		go func(fullURL *url.URL, currDepth int) {
			defer s.wg.Done()
			defer func() {
				<-s.sem
			}()

			dirPath := path.Join(s.dirName, path.Dir(fullURL.Path))
			if err := os.MkdirAll(dirPath, 0750); err != nil {
				s.errsMu.Lock()
				s.allErrs = errors.Join(s.allErrs, fmt.Errorf("failed to create a directory: %w", err))
				s.errsMu.Unlock()
				return
			}

			var ext string
			if path.Ext(fullURL.Path) == "" {
				ext = "index.html"
			}
			filePath := path.Join(s.dirName, fullURL.Path, ext)
			if err := download.DownloadFile(fullURL.String(), filePath); err != nil {
				s.errsMu.Lock()
				s.allErrs = errors.Join(s.allErrs, err)
				s.errsMu.Unlock()
				return
			}

			if ext := path.Ext(filePath); ext == "" || ext == ".html" {
				if ext == "" {
					filePath = path.Join(filePath, "index.html")
				}

				file, err := os.Open(filePath)
				if err != nil {
					s.errsMu.Lock()
					s.allErrs = errors.Join(s.allErrs, fmt.Errorf("failed to open file: %w", err))
					s.errsMu.Unlock()
					return
				}
				defer file.Close()

				doc, err := html.Parse(file)
				if err != nil {
					s.errsMu.Lock()
					s.allErrs = errors.Join(s.allErrs, fmt.Errorf("failed to parse HTML file: %w", err))
					s.errsMu.Unlock()
					return
				}

				var urls []string
				parse.ExtractHTML(doc, &urls)

				currDepth++
				if err := s.Store(urls, currDepth, fullURL); err != nil {
					s.errsMu.Lock()
					s.allErrs = errors.Join(s.allErrs, fmt.Errorf("failed to store a file: %w", err))
					s.errsMu.Unlock()
					return
				}
			}
		}(fullURL, currDepth)
	}

	return nil
}
