package sg

import (
	"encoding/json"

	"github.com/evict/cody-security-review/logging"
)

type FileDiff struct {
	OldPath string
	NewPath string
	Hunks   []Hunk
}

type DiffResponse struct {
	Data struct {
		Node struct {
			Comparison struct {
				FileDiffs struct {
					Nodes []FileDiff
				}
			}
		}
	}
}

type Hunk struct {
	Body string
}

func FetchFile(branch string, file string) string {

	query := `
query FileContent($repoName: String!, $revision: String!, $path: String!) {
  repository(name: $repoName) {
    commit(rev: $revision) {
      blob(path: $path) {
        content
      }
    }
  }
}`

	variables := map[string]interface{}{
		"repoName": "github.com/sourcegraph/sourcegraph",
		"revision": branch,
		"path":     file,
	}

	// Construct the request
	reqBody, err := json.Marshal(map[string]interface{}{
		"query":     query,
		"variables": variables,
	})

	if err != nil {
		panic(err)
	}

	// DO WE WANT THIS FROM MAIN?
	logging.Debug("fetching file from branch")
	httpClient.AddHeader("Content-Type", "application/json")
	resp, err := httpClient.PostRequest("https://sourcegraph.sourcegraph.com/.api/graphql", reqBody)

	if resp.StatusCode != 200 {
		panic(err)
	}

	type QueryResult struct {
		Repository struct {
			Commit struct {
				Blob struct {
					Content string `json:"content"`
				} `json:"blob"`
			} `json:"commit"`
		} `json:"repository"`
	}

	var result QueryResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		panic(err)
	}

	return result.Repository.Commit.Blob.Content

}

func GetComparisonDiff(base string, head string) []FileDiff {

	if httpClient.Client == nil {
		initClient()
	}

	query := `
query ComparisonDiff($nodeID: ID!, $base: String!, $head: String!) {
  node(id: $nodeID) {
	... on Repository {
	  comparison(base: $base, head: $head) {
		fileDiffs {
		  nodes {
			oldPath
			newPath
			hunks {
			  body 
			}
		  }
		}
	  }
	}
  }
}`

	variables := map[string]interface{}{
		"nodeID": "UmVwb3NpdG9yeTozOTk=",
		"base":   base,
		"head":   head,
	}

	// Construct the request
	reqBody, err := json.Marshal(map[string]interface{}{
		"query":     query,
		"variables": variables,
	})

	if err != nil {
		panic(err)
	}

	logging.Debug("fetching diffs from branch")
	httpClient.AddHeader("Content-Type", "application/json")
	resp, err := httpClient.PostRequest("https://sourcegraph.sourcegraph.com/.api/graphql", reqBody)

	if resp.StatusCode != 200 {
		panic(err)
	}

	var dr DiffResponse

	err = json.NewDecoder(resp.Body).Decode(&dr)
	if err != nil {
		panic(err)
	}

	var diffs []FileDiff

	for _, node := range dr.Data.Node.Comparison.FileDiffs.Nodes {
		diff := FileDiff{
			OldPath: node.OldPath,
			NewPath: node.NewPath,
			Hunks:   make([]Hunk, len(node.Hunks)),
		}

		for i, hunk := range node.Hunks {
			diff.Hunks[i] = Hunk{
				Body: hunk.Body,
			}
		}

		diffs = append(diffs, diff)
	}

	return diffs

}
