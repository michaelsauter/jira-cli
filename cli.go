package main

import (
	"fmt"
	"github.com/alecthomas/kingpin"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

var (
	app         = kingpin.New("jira", "JIRA CLI").Interspersed(false).DefaultEnvars()
	verboseFlag = app.Flag("verbose", "Enable verbose output.").Short('v').Bool()
	projectFlag = app.Flag("project", "JIRA project").Short('p').String()

	newCommand = app.Command(
		"new",
		"Create new JIRA issue and switch to its branch",
	)
	newTypeFlag = newCommand.Flag(
		"type",
		"Issue type",
	).Short('t').Default("Task").String()

	newTitleArg = newCommand.Arg("title", "Title").String()
)

func realMain() {
	err := readConfig()
	if err != nil {
		os.Exit(1)
		return
	}

	command := kingpin.MustParse(app.Parse(os.Args[1:]))

	switch command {
	case newCommand.FullCommand():
		openEditor()
		desc := readDescription()
		project := *projectFlag
		verboseMsg("Read project '" + project + "' from CLI")
		if len(project) == 0 {
			project = guessProject()
		}
		verboseMsg("Determined project '" + project + "'")

		issue, _ := createIssue(project, *newTitleArg, *newTypeFlag, desc)

		branchName := issue.Key + "-" + strings.ToLower(*newTitleArg)
		branchName = strings.Replace(branchName, " ", "-", -1)

		createBranch(branchName)
		pushBranch()
	}
}

func isVerbose() bool {
	return *verboseFlag
}

func verboseMsg(message string) {
	if isVerbose() {
		fmt.Printf("--> %s\n", message)
	}
}

func openEditor() {
	verboseMsg("Opening editor")
	cmd := exec.Command(config.Editor, ".JIRA_DESCRIPTION")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func readDescription() string {
	verboseMsg("Reading file")
	dat, _ := ioutil.ReadFile(".JIRA_DESCRIPTION")
	verboseMsg("Deleting file")
	os.Remove(".JIRA_DESCRIPTION")
	return string(dat)
}
