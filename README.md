# cody-security-review
A tool that helps with security code review, powered by Sourcegraph's code AI platform.

## Building

```
go build -o cody-review main.go
```

Move to your favorite `bin` folder.

## Usage

Generate a token in your Sourcegraph account, set it as an environment value `CODY_AUTH_TOKEN`. In `sg/client.go` you can change the base URL to your Sourcegraph instance.

Invoke tool with branch name:

```
./cody-review [branch_name]

DEBUG: 2023/07/20 16:53:27 debug.go:21: 2023-07-20T15:53:27.033187Z fetching diffs from branch
git diff -- a/cmd/gitserver/server/server.go
 		s.setCloneStatusNonFatal(context.Background(), repo, cloneStatus(repoCloned(dir), false))
 	}()
 
-	cmd, err := syncer.CloneCommand(ctx, remoteURL, tmpPath)
-	if err != nil {
-		return errors.Wrap(err, "get clone command")
-	}
-	if cmd.Env == nil {
-		cmd.Env = os.Environ()
-	}

Do you want Cody to review this diff? (y/n) y

DEBUG: 2023/07/20 16:54:17 debug.go:21: 2023-07-20T15:54:17.572558Z full file not in context, fetching and adding
DEBUG: 2023/07/20 16:54:17 debug.go:21: 2023-07-20T15:54:17.573074Z fetching file from branch
DEBUG: 2023/07/20 16:54:17 debug.go:21: 2023-07-20T15:54:17.760324Z added file: cmd/gitserver/server/server.go
DEBUG: 2023/07/20 16:54:17 debug.go:21: 2023-07-20T15:54:17.760668Z requesting completions...
DEBUG: 2023/07/20 16:54:20 debug.go:21: 2023-07-20T15:54:20.703470Z fetching completion events
 Based on my review, I do not see any clear security vulnerabilities introduced in this diff. The key changes appear to be refactoring the clone logic to support using a new Cloner interface, and moving the call to setLastOutput to after error handling. 

The refactoring to support Cloner does not in itself introduce any new vulnerabilities - it just allows different clone implementations to be used in a consistent way.

The move of setLastOutput also does not appear to be a security issue. It now happens after error handling, which seems reasonable.

To summarize:

REVIEW: NO

I did not identify any definite security vulnerabilities in this diff. The changes appear to be refactoring and reorganization rather than introducing new logic. Please let me know if you would like me to examine any other specific parts of the diff or codebase. I'm happy to take another look.

```
