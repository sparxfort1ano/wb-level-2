package command

import (
	"cmp"
	"fmt"
	"io"
	"slices"
	"text/tabwriter"
	"time"

	"github.com/tklauser/ps"
)

// ProcessStatus retrieves a list of all currently running processes on the system.
// It sorts the processes numerically by PID and writes a formatted table
// containing the PID, elapsed execution time and executable path to the
// provided output system.
func ProcessStatus(outStream io.Writer) error {
	procs, err := ps.Processes()
	if err != nil {
		return fmt.Errorf("failed to list processes: %w", err)
	}

	slices.SortFunc(procs, func(a, b ps.Process) int {
		return cmp.Compare(a.PID(), b.PID())
	})

	w := tabwriter.NewWriter(outStream, 0, 0, 1, ' ', 0)
	fmt.Fprintf(w, "PID\tTIME\tCMD\n")
	for _, p := range procs {
		fmt.Fprintf(w, "%d\t%s\t%s\n",
			p.PID(),
			time.Unix(int64(time.Since(p.CreationTime()).Seconds()), 0).UTC().Format("15:04:05"),
			p.ExecutablePath(),
		)
	}
	w.Flush()

	return nil
}
