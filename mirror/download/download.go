// Package download provides network utilities for fetching remote files.
package download

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// DownloadFile fetches a resource from the specified URL and saves it to a local file.
func DownloadFile(url, filename string) error {
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create a request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	response, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make http-request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("response error")
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, response.Body); err != nil {
		return fmt.Errorf("failed to copy overload to file: %w", err)
	}

	return nil
}
