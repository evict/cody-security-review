package sg

import (
	"fmt"
	"testing"
)

func TestGetComparisonDiff(t *testing.T) {

	diffs := GetComparisonDiff("main", "security/sonarcloud-buildkite")

	for _, d := range diffs {
		fmt.Println(d.NewPath)
		fmt.Println(d.OldPath)

		for _, h := range d.Hunks {
			fmt.Println(h)
		}
	}
}
