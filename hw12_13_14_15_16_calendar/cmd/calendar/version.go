package main

import (
	"fmt"
)

var (
	release   = "UNKNOWN"
	buildDate = "UNKNOWN"
	gitHash   = "UNKNOWN"
)

func getVersion() string {
	return fmt.Sprintf("release: %s - buildDate: %s - gitHash: %s\n", release, buildDate, gitHash)
}

// func printVersion() {
// 	if err := json.NewEncoder(os.Stdout).Encode(struct {
// 		Release   string
// 		BuildDate string
// 		GitHash   string
// 	}{
// 		Release:   release,
// 		BuildDate: buildDate,
// 		GitHash:   gitHash,
// 	}); err != nil {
// 		fmt.Printf("error while decode version info: %v\n", err)
// 	}
// }
