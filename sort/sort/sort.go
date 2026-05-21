// Package sort is responsible for sorting the given files or sort checking them (depends on
// `options.Options` parameters).
// It does this efficiently in terms of resource usage because it uses temporary files
// and a binary heap, resulting in a K-Way Merge.
package sort

import (
	"bufio"
	"container/heap"
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/sparxfort1ano/wb-level-2/sort/options"
)

type TempFiles []*os.File

// Close closes and removes all the opened temporary files.
func (tf TempFiles) Close() {
	for _, file := range tf {
		file.Close()
		//os.Remove(file.Name())
	}
}

const chunkSizeLimit = 20

// RunSort serves as a core function that collecting all the sub-functions responsible for
// reading input files into a stream, sorting the stream, converting the stream into a temporary file
// and merging these sorted files.
func RunSort(out io.Writer, opts *options.Options) error {
	if opts.IsSorted {
		err := checkFileSorted(opts)
		if err == nil {
			fmt.Fprintln(out, "the data is sorted")
			return nil
		}
		return err
	}

	var counterSize int
	var chunk []string
	var temps TempFiles
	defer temps.Close()

	for _, file := range opts.Inputs {
		reader := bufio.NewReader(file)

		for {
			str, err := reader.ReadString('\n')
			if err != nil {
				if errors.Is(err, io.EOF) {
					if str != "" {
						if !strings.HasSuffix(str, "\n") {
							str += "\n"
						}

						chunk = append(chunk, str)
						counterSize += len(str)

						if counterSize >= chunkSizeLimit {
							f, err := sortChunk(chunk, opts)
							if err != nil {
								return err
							}

							temps = append(temps, f)

							counterSize = 0
							chunk = chunk[:0]
						}
					}
					break
				}
				return fmt.Errorf("failed to read string from an input file: %w", err)
			}

			chunk = append(chunk, str)
			counterSize += len(str)

			if counterSize >= chunkSizeLimit {
				f, err := sortChunk(chunk, opts)
				if err != nil {
					return err
				}

				temps = append(temps, f)

				counterSize = 0
				chunk = chunk[:0]
			}
		}
	}

	if len(chunk) > 0 {
		f, err := sortChunk(chunk, opts)
		if err != nil {
			return err
		}
		temps = append(temps, f)
	}

	if err := printMergedSort(out, opts, temps); err != nil {
		return err
	}

	return nil
}

func checkFileSorted(opts *options.Options) error {
	if len(opts.Inputs) > 1 {
		return fmt.Errorf("2 or more files are not allowed with -c")
	}

	reader := bufio.NewReader(opts.Inputs[0])

	curr, errRead := reader.ReadString('\n')
	if errRead != nil {
		if errors.Is(errRead, io.EOF) {
			return nil
		}
		return fmt.Errorf("failed to read a string line from the buffer: %w", errRead)
	}
	for {
		prev := curr
		curr, errRead = reader.ReadString('\n')
		if errRead != nil {
			if errors.Is(errRead, io.EOF) {
				if curr != "" {
					if !strings.HasSuffix(curr, "\n") {
						curr += "\n"
					}

					cmpResult := opts.Compare(prev, curr)

					if cmpResult > 0 {
						return fmt.Errorf("file is not sorted: disorder at %s", strings.TrimSpace(curr))
					}
					if opts.Unique && cmpResult == 0 {
						return fmt.Errorf("data elements are not unique: duplicate %s", strings.TrimSpace(curr))
					}
				}
				break
			}
			return fmt.Errorf("failed to read a string line from the buffer: %w", errRead)
		}

		cmpResult := opts.Compare(prev, curr)

		if cmpResult > 0 {
			return fmt.Errorf("file is not sorted: disorder at %s", strings.TrimSpace(curr))
		}
		if opts.Unique && cmpResult == 0 {
			return fmt.Errorf("data elements are not unique: duplicate %s", strings.TrimSpace(curr))
		}
	}

	return nil
}

func sortChunk(chunk []string, opts *options.Options) (*os.File, error) {
	slices.SortFunc(chunk, opts.Compare)
	if opts.Unique {
		chunk = slices.CompactFunc(chunk, opts.Equal)
	}

	f, err := os.CreateTemp("", "sort_file_*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary file: %w", err)
	}

	writer := bufio.NewWriter(f)
	for _, line := range chunk {
		if _, err := writer.WriteString(line); err != nil {
			return nil, fmt.Errorf("failed to write line into a temporary file: %w", err)
		}
	}
	if err := writer.Flush(); err != nil {
		return nil, fmt.Errorf("failed to flush buffer to temporary file: %w", err)
	}

	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return nil, fmt.Errorf("failed to seek to start of temporary file: %w", err)
	}

	return f, nil
}

func printMergedSort(out io.Writer, opts *options.Options, temps TempFiles) error {
	pq := NewPriorityQueue(
		make([]*MinHeapItem, 0, len(temps)),
		opts,
	)

	for _, temp := range temps {
		reader := bufio.NewReader(temp)

		line, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				if line != "" {
					item := NewMinHeapItem(line, reader)
					pq.Items = append(pq.Items, item)
				}
				continue
			}
			return fmt.Errorf("failed to read a string line from the buffer: %w", err)
		}

		item := NewMinHeapItem(line, reader)
		pq.Items = append(pq.Items, item)
	}

	heap.Init(pq)

	var prevLine string
	var isFirst = true

	for pq.Len() > 0 {
		item := heap.Pop(pq).(*MinHeapItem)

		isDuplicate := false
		if opts.Unique && !isFirst {
			if opts.Equal(item.Line, prevLine) {
				isDuplicate = true
			}
		}

		if !isDuplicate {
			fmt.Fprint(out, item.Line)
			prevLine = item.Line
			isFirst = false
		}

		nextLine, err := item.Reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				if nextLine != "" {
					item.Line = nextLine
					heap.Push(pq, item)
				}
				continue
			}
			return fmt.Errorf("failed to read a string line from the buffer: %w", err)
		}

		item.Line = nextLine
		heap.Push(pq, item)
	}

	return nil
}
