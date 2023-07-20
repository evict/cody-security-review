package sg

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/sourcegraph/cody-code-review/logging"
)

func GetCompletions(query string, file string, branch string) string {

	if httpClient.Client == nil {
		err := initClient()

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	context := getCodyContext(query, file, branch)
	// Define the initial message
	messages := []map[string]string{
		{"speaker": "human", "text": "You are Cody, an AI-powered code security expert developed by Sourcegraph. You operate in a Command Line Interface (CLI) and have comprehensive access to all files. Your primary responsibilities include: \n\n- Analyzing code differences (diffs) to identify potential security vulnerabilities. \n- Responding to inquiries concerning code security. It's fine if there are no vulnerabilities in the code, only respond with 100% certainty. \n- Clearly explaining potential security risks within the code to non-technical stakeholders. \n\nIn your responses, please adhere to the following rules: \n\n- In scenarios where you don't have access to required code, files, or repositories, maintain your character as Cody while apologizing. \n- Aim to provide responses that are concise yet comprehensive, ensuring clarity is not compromised. \n- All provided code snippets must adhere to Markdown formatting rules and be placed within triple backticks like so: ```code snippet```. \n- Only provide an answer if you are certain or can make a well-informed prediction. If you lack the necessary information to provide a response, admit the knowledge gap and indicate what additional context would allow you to provide an answer. \n- When referencing file names, repository names, or URLs, ensure their existence. \n\nYou have access to all files and are capable of analyzing code diffs and responding to security queries. I will provide the necessary code diffs and snippets for your analysis and queries. You are an expert in the OWASP Top 10 vulnerabilities and are knowledgeable in securing input/output data."},
		{"speaker": "assistant", "text": "Understood. As Cody, the AI security specialist from Sourcegraph, I am equipped to analyze code, identify potential security risks and provide succinct, comprehensible explanations about these issues. I operate within a CLI and have access to all files and diffs. If I encounter scenarios where I lack access or knowledge, I will maintain my character while admitting to these limitations. For any code snippets, I will ensure they are formatted correctly using Markdown syntax. All of my responses will aim to be concise and clear. With expertise in the OWASP Top 10 vulnerabilities and the securing of input/output data, I am ready to assist."},
		{"speaker": "human", "text": "Cody, your knowledge encompasses a wide range of application vulnerabilities, including Server Side Request Forgery (SSRF) and Command Injection. SSRF allows an attacker to induce the server to make a request to an arbitrary domain, often resulting in unauthorized actions or access to data from the server. Typical methods for conducting SSRF are using webhooks are internal HTTP clients, but anything that makes outgoing connections can be used. SSRF can also be used to contact internal services, such as with outgoing webhooks. But also for DNS enumeration and port scanning. Command Injection attacks occur when an attacker can insert arbitrary commands that are executed by the server. Can you confirm that you understand these concepts and can apply this knowledge to your analysis of code security?"},
		{"speaker": "assistant", "text": "Absolutely. I'm well-versed in application vulnerabilities, including Server Side Request Forgery (SSRF) and Command Injection. I understand that SSRF can cause a server to perform unauthorized actions or reveal sensitive data by manipulating it to send a request to a manipulated domain. Command Injection involves an attacker managing to get a server to execute arbitrary commands. This knowledge forms a key part of my approach to analyzing code security. I will make sure to check for any signs of these vulnerabilities in the code I analyze."},
	}
	messages = append(messages, context...)

	/* 	cwd, _ := os.Getwd()

	   	content, err := os.ReadFile(fmt.Sprintf("%s/%s", cwd, filename))
	   	if err != nil {
	   		panic(err)
	   	}
	   	fileP = fmt.Sprintf("%s\n%s\ncontent:\n%s\n", fileP, filename, content) */

	qmsg := []map[string]string{
		{"speaker": "human", "text": query},
		{"speaker": "assistant", "text": ""},
	}

	messages = append(messages, qmsg...)

	reqBody, _ := json.Marshal(map[string]interface{}{
		"temperature":       0.1,
		"maxTokensToSample": 2000,
		"topK":              -1,
		"topP":              -1,
		"messages":          messages,
	})

	logging.Debug("requesting completions...")
	httpClient.AddHeader("Accept", "text/event-stream")
	resp, err := httpClient.PostRequest("https://sourcegraph.sourcegraph.com/.api/completions/stream", reqBody)

	if err != nil {
		panic(err)
	}

	logging.Debug("fetching completion events")

	var lastEvent Event

	reader := bufio.NewReader(resp.Body)
	var event Event
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		line = strings.TrimSuffix(line, "\n")
		parts := strings.SplitN(line, ": ", 2)
		if len(parts) < 2 {
			continue
		}

		i := 0

		switch parts[0] {
		case "event":
			event.Name = parts[1]
		case "data":
			json.Unmarshal([]byte(parts[1]), &event.Data)

			if event.Name == "done" {
				fmt.Printf("\r\033[K")
				return lastEvent.Data.Completion
			} else {
				if i == 10 {
					fmt.Printf("\r\033[K")
					i = 0
				} else {
					fmt.Printf(".")
					i++
				}
			}

			lastEvent = event
			event = Event{}
		}
	}

	return ""
}

func returnContextPrompt(path string, reponame string, content string) map[string]string {

	return map[string]string{
		"speaker": "human",
		"text": fmt.Sprintf("Use following code snippet from file `%s` in repository `%s`:\n```go\n%s\n```",
			path,
			reponame,
			content,
		),
	}
}

func getCodyContext(cquery string, path string, branch string) []map[string]string {

	// 	// TODO MAKE SURE THIS IS FROM RIGHT BRANCH TOO

	// 	// GraphQL query
	// 	query := `
	// query GetCodyContext($repos: [ID!]!, $query: String!, $codeResultsCount: Int!, $textResultsCount: Int!) {
	//    getCodyContext(repos: $repos, query:$ query, codeResultsCount: $codeResultsCount, textResultsCount: $textResultsCount) {
	// 	   ... on FileChunkContext {
	// 		   blob {
	// 			   path
	// 			   repository {
	// 				   id
	// 				   name
	// 			   }
	// 		   }
	// 		   startLine
	// 		   endLine
	// 		   chunkContent
	// 	   }
	//    }
	// }`

	// 	// Variables
	// 	variables := map[string]interface{}{
	// 		"repos":            []string{"UmVwb3NpdG9yeTozOTk="},
	// 		"query":            cquery,
	// 		"codeResultsCount": 5,
	// 		"textResultsCount": 5,
	// 	}

	// 	// Construct the request
	// 	reqBody, err := json.Marshal(map[string]interface{}{
	// 		"query":     query,
	// 		"variables": variables,
	// 	})
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	logging.Debug("fetching context for query")
	// 	httpClient.AddHeader("Content-Type", "application/json")
	// 	resp, err := httpClient.PostRequest("https://sourcegraph.sourcegraph.com/.api/graphql", reqBody)

	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	if resp.StatusCode != 200 {
	// 		panic(resp.Body)
	// 	}

	// 	var result struct {
	// 		Data struct {
	// 			GetCodyContext []struct {
	// 				Blob struct {
	// 					Path       string
	// 					Repository struct {
	// 						Name string
	// 					}
	// 				}
	// 				StartLine    int
	// 				EndLine      int
	// 				ChunkContent string
	// 			}
	// 		}
	// 	}
	// 	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
	// 		panic(err)
	// 	}

	// 	logging.Debug("adding context to prompt")
	var messages []map[string]string

	// 	var usedFiles = make(map[string]bool)

	// 	for _, file := range result.Data.GetCodyContext {
	// 		if !usedFiles[file.Blob.Path] {
	// 			humanMessage := returnContextPrompt(file.Blob.Path,
	// 				file.Blob.Repository.Name,
	// 				file.ChunkContent)
	// 			logging.Debug(fmt.Sprintf("using file: %s", file.Blob.Path))
	// 			messages = append(messages, humanMessage)

	// 			assistantMessage := map[string]string{
	// 				"speaker": "assistant",
	// 				"text":    "Ok.",
	// 			}
	// 			messages = append(messages, assistantMessage)

	// 			usedFiles[file.Blob.Path] = true
	// 		}
	// 	}

	logging.Debug("full file not in context, fetching and adding")
	content := FetchFile(branch, path)
	logging.Debug(fmt.Sprintf("added file: %s", path))

	humanMessage := returnContextPrompt(path, "github.com/sourcegraph/sourcegraph", content)
	messages = append(messages, humanMessage)
	assistantMessage := map[string]string{
		"speaker": "assistant",
		"text":    "Ok.",
	}
	messages = append(messages, assistantMessage)

	return messages

}
