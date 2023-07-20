package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/evict/cody-security-review/sg"
)

func main() {

	var branch string

	if len(os.Args) > 1 {
		branch = os.Args[1]
	} else {
		fmt.Println("error: branch argument missing")
		os.Exit(1)
	}

	diffs := sg.GetComparisonDiff("main", branch)

	preamble := "Does the diff negatively impact the security of the application? I provided you with the full content of the file as well, for more thorough review. Is there a vulnerability here? If so, how do i mitigate it? Add a relevant code sample where possible to improve the security of the code, but never respond with another diff. If you think there are security vulnerabilities in this code, reply exactly with REVIEW:YES. Diff: "

	var path string

	query := preamble
	for _, d := range diffs {

		path = d.NewPath

		for _, h := range d.Hunks {
			ext := filepath.Ext(d.NewPath)

			switch ext {
			case ".java", ".py", ".c", ".cpp", ".go", ".js", ".ts", ".cs", ".rb", ".php", ".tsx":

				// exclude tests for go, should go into a config file
				if strings.Contains(d.NewPath, "_test.go") {
					continue
				}

				dbody := fmt.Sprintf("git diff -- a/%s\n%s", d.NewPath, (h.Body))
				fmt.Println("```")
				fmt.Println(dbody)
				fmt.Println("```")
				var response string
				fmt.Printf("Do you want Cody to review this diff? (y/n) ")
				fmt.Scanln(&response)
				fmt.Println()

				if response == "y" {
					r := sg.GetCompletions(query+dbody, path, branch)
					fmt.Println(r)
				}

			}
		}
	}
}
