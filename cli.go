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
	).Short('t').Default("Bug").String()

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

		issue, _ := createIssue(readProject(), *newTitleArg, *newTypeFlag, desc)

		branchName := issue.Key + "-" + strings.ToLower(*newTitleArg)
		branchName = strings.Replace(branchName, " ", "-", -1)

		createBranch(branchName)
		pushBranch()
	}
}

func openEditor() {
	fmt.Println("opening editor")
	cmd := exec.Command(config.Editor, ".JIRA_DESCRIPTION")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func readDescription() string {
	fmt.Println("reading file")
	dat, _ := ioutil.ReadFile(".JIRA_DESCRIPTION")
	fmt.Println("deleting file")
	os.Remove(".JIRA_DESCRIPTION")
	return string(dat)
}
