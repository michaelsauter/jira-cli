package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func createBranch(branchName string) {
	fmt.Println("creating branch")
	cmd := exec.Command("git", "checkout", "-b", branchName)
	err := cmd.Run()
	fmt.Println(err)
}

func pushBranch() {
	fmt.Println("pushing branch")
	cmd := exec.Command("git", "push")
	err := cmd.Run()
	fmt.Println(err)
}

func readProject() string {
	out, _ := exec.Command("git", "config", "remote.origin.url").Output()
	repo := strings.Split(string(out), "/")
	return strings.ToUpper(repo[len(repo)-2])
}
