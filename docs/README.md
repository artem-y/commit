# commit
Simple CLI tool that finds an issue number in the branch and includes it in the commit message.

## Installation
### Install Pre-built Executable
1. In the [Releases](https://github.com/artem-y/commit/releases) section, find and open the latest release (or any other one if you want).
2. Download `bin.zip` archive
3. Unarchive it and find a folder for your machine's architecture.
4. Make sure your system allows execution of the tool. For example, on Linux/macOS, use `chmod +x <path-to-commit-executable>` command to enable the "execute" permission. Adjust other settings if needed.
5. Move the executable into one of the directories visible in `PATH` so that the command is visible and reachable from anywhere.
6. If you want to use a custom path to a global `.commit.json` file for all the repositories, you can specify it with a typealias in your `.bashrc` or `.zshrc` file etc.  
   For example:  
   ```shell
   alias commit="commit -config-path=${HOME}/.config/.commit.json"
   ```
### Install With Makefile
1. Make sure you have a compatible version of `Go` installed _(see [go.mod](https://github.com/artem-y/commit/blob/main/go.mod#L3) file)_
2. Clone the repository
3. In the root of the repository, run `make build` command
4. Find the executable in `bin` directory, move it into one of the directories visible in `PATH` so that it is visible and reachable from anywhere.
5. If you want to use a custom path to a global `.commit.json` file for all the repositories, you can specify it with a typealias in your `.bashrc` or `.zshrc` file etc.  
   For example:  
   ```shell
   alias commit="commit -config-path=${HOME}/.config/.commit.json"
   ```  

Check out the [Makefile](/Makefile) for more commands.

### Easy Installation With Makefile
There's a helper command that can simplify installation:
```shell
make install
```
1. Make sure you have a compatible version of `Go` installed _(see [go.mod](https://github.com/artem-y/commit/blob/main/go.mod#L3) file)_
2. Clone the repository
3. In the root of the repository, run `make install`
4. _(Optional)_ It will try to install to where the tool was installed previously. In case this is the first installation, it will be installed in the default Go package directory. If the tool is not visible after instllation, make sure its directory is included in `PATH`.

### Install With Go Tooling
You can do the steps described above in the [Install With Makefile](#install-with-makefile) section, replacing Step 3 (`make build` command) with just a plain Go build command:  
```shell
go build -o bin/ ./cmd/commit
```
Or you can skip even those steps and use `go install` command:
```shell
go install -ldflags "-s -w" ./cmd/commit
```
This will install the executable in the default Go package directory. If the tool is not visible after installation, make sure its directory is included in `PATH`.

For more information, see [Tutorial: Compile and install the application](https://go.dev/doc/tutorial/compile-install) 
## Usage
To commit all staged changes, use `commit` command with some commit message:
```shell
commit "Refactor core service initialization logic"
```
The result: for example, if the branch name is `312-improve-stability-of-the-core-service`, the resulting commit will have a message:

> #312: Refactor core service initialization logic
### Configuration
By default, the tool recognizes the pattern suggested by GitHub when auto-generating branches, when issue numbers are just digits: `28-update-documentation` for the issue  `#28`.  
But this can be changed by providing a `.commit.json` file with different values at the root of your repository:  
```json
{  
    "issueRegex": "ABC-[0-9]+", 
    "outputIssuePrefix": "#",
    "outputIssueSuffix": "",
    "outputStringPrefix": "",
    "outputStringSuffix": ": ",
}
```
What each setting does:
- **issueRegex**: This is how the tool determines what is the pattern to look for.
- **outputIssuePrefix**: Added before each issue.
- **outputIssueSuffix**: Added after each issue.
- **outputStringPrefix**: Precedes the generated part of the commit message.
- **outputStringSuffix**: Follows at the end of the generated part of the commit message.

The structure of the resulting commit message is as follows:  
```
<outputStringPrefix><outputIssuePrefix><issueRegex><outputIssueSuffix>, <outputIssuePrefix><issueRegex><outputIssueSuffix><outputStringSuffix> <commit message>
```
If the `.commit.json` file is not included, the tool will just fall back to its default settings (GitHub style issues). See the default values in [constants.go](/internal/helpers/constants.go) file.
### Custom Config Path
If you don't want to include the `.commit.json` file at the root of your repository, path to the config file can be passed with a `-config-path` flag like this (linux/macOS shell example):
```shell
commit -config-path=${HOME}/.config/.commit.json "Finally fix everything"
```
### Multiple Issue Numbers
If the branch has multiple issues in its name, the tool will include them all, comma-separated.  
For example, the branch named `add-tests-for-CR-127-and-CR-131-features`, the issue regex set to `[A-Z]{2}-[0-9]+`, and the "outputIssuePrefix" and "outputIssueSuffix" settings for the output set to `[` and `]:`, the generated commit message would start with the following:  
> [CR-127, CR-131]:
## Testing
So far only the core parts of the logic are covered with tests, but I will be adding more in the future.  
To run the unit tests, use the `make test` or `go test -v ./...` commands.  
There also are a few end-to-end tests that will run in the CI pipeline to make sure the tool actually works as expected with a real repository. I do not advise running those often manually, since they can potentially execute destructive commands. For example, there is at least one change the tests will try to apply to the global git setup, which is the setting of the default init branch to "main".
### Debug
To check what commit message will be generated without making the actual commit, there is a `-dry-run` flag that can be passed to the command:
```shell
commit -dry-run "Not a real commit"
```
## Contributing
Originally I wrote this tool for myself in shell and Swift and used a lot on MacOS. This repo is an attempt to make it crossplatform and an opportunity to excercise in writing Go.  
If you find it useful and see that something's wrong or missing, feel free to raise issues and contribute to the project.  
When doing so, please follow the ususal code of conduct and contribution guidelines from GitHub, and just common sense.
## License
This repository is distributed under the MIT license (see [LICENSE](/LICENSE)).  
My tool uses [go-git](https://github.com/go-git/go-git) as a 3rd party dependency to work with git directly, you might want to check it out too.
