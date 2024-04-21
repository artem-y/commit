# commit
Simple CLI tool that finds an issue number in the branch and includes it in the commit message

## Usage
To commit all staged changes, use `commit` command with some commit message:
```shell
commit "Refactor core service initialization logic"
```
The result: for example, if the branch name is `312-improve-stability-of-the-core-service`, the resulting commit will have a message:

> #312: Refactor core service initialization logic
### Configuration
By default, the tool recognizes the pattern suggested by GitHub when auto-generating branches, when issue numbers are just digits: `28-update-documentation` for the issue  `#28`.  
But this can be changed by setting different values in a `.commit.json` file at the root of your repository:  
```json
{  
    "issueRegex": "ABC-[0-9]+", // this is how the tool determines what is the pattern to look for
    "outputIssuePrefix": "#", // precedes the generated part of the commit message
    "outputIssueSuffix": ": " // follows at the end of the generated part of the commit message
}
```
The structure of the resulting commit message is as follows:  
```
<outputIssuePrefix><issueRegex><outputIssueSuffix> <commit message>
```
If the `.commit.json` file is not included, the tool will just fall back to its default settings (GitHub style issues). 
### Custom Config Path
If you don't want to include the `.commit.json` file at the root of your repository, path to the config file can be passed with a `-config-path` flag like this (linux/macOS shell example):
```shell
commit -config-path=${HOME}/.config/.commit.json "Finally fix everything"
```
### Multiple Issue Numbers
If the branch has multiple issues in its name, the tool will include them all, comma-separated.  
For example, the branch named `add-tests-for-CR-127-and-CR-131-features`, the issue regex set to `[A-Z]{2}-[0-9]+`, and the "outputIssuePrefix" and "outputIssueSuffix" settings for the output set to `[` and `]:`, the generated commit message would start with the following:  
> [CR-127, CR-131]: 
