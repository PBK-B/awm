package cli

import (
	"fmt"

	"github.com/pbk-b/awm/internal/version"
)

func cmdVersion(args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("usage: awm version")
	}
	fmt.Println(version.String())
	return nil
}
