# commit
Simple CLI tool that finds an issue number in the branch and includes it in the commit message

## Usage
To commit all staged changes, use `commit` command with some commit message:
```shell
commit "Refactor core service initialization logic"
```
The result: for example, if the branch name is `312-improve-stability-of-the-core-service`, the resulting commit will have a message:

> #312: Refactor core service initialization logic
