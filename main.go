package main

import (
	"fmt"
	"io"
	"os"

	"github.com/dustin/go-humanize"
)

func run(args []string, out io.Writer) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: sizeconv '1.5 GiB'")
	}
	size, err := humanize.ParseBytes(args[0])
	if err != nil {
		return fmt.Errorf("invalid size: %w", err)
	}
	_, err = fmt.Fprintf(out, "%d bytes\n", size)
	return err
}

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
