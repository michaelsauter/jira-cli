package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func createBranch(branchName string) {
	verboseMsg("creating branch")
	cmd := exec.Command("git", "checkout", "-b", branchName)
	err := cmd.Run()
	if err == nil {
		fmt.Printf("Checked out branch: %s", branchName)
	} else {
		fmt.Println("ERROR! Branch could not be created.")
	}
}

func pushBranch() {
	verboseMsg("pushing branch")
	cmd := exec.Command("git", "push")
	err := cmd.Run()
	if err != nil {
		fmt.Println("ERROR! Branch could not be pushed.")
	}
}

func guessProject() string {
	out, _ := exec.Command("git", "config", "remote.origin.url").Output()
	repo := strings.Split(string(out), "/")
	return strings.ToUpper(repo[len(repo)-2])
}
