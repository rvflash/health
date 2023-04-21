package health

import (
	"context"
	"fmt"
	"os"
)

// CreateFileCheck is a Check dedicated to create a file, write a content inside it and delete it.
// Useful to verify the capacity to write into a directory or access to an NFS mount.
func CreateFileCheck(dir, pattern string) Check {
	return func(_ context.Context) error {
		f, err := os.CreateTemp(dir, pattern)
		if err != nil {
			return fmt.Errorf("creating: %w", err)
		}
		_, err = f.Write([]byte("content"))
		if err != nil {
			return fmt.Errorf("writing: %w", err)
		}
		err = f.Close()
		if err != nil {
			return fmt.Errorf("closing: %w", err)
		}
		return os.Remove(f.Name())
	}
}
