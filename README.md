
# Overview

A CLI tool to work with JIRA.

# What it does

Given a config file in your home directory and information taken from the Git repository in which the CLI is run, it allows to create a new issue from the command line, create a branch for it and switch to it in your local repository.

# Usage

Run `jira new "my new issue"`. The first time, this will prompt you for some setup. Following executions will open your configured editor so that you can enter the description. Upon save, it creates an issue in JIRA, assigned to you with the given `--type <type>` (defaults to `Task`). It then creates a new branch following the JIRA naming conventions and pushes it to the remote.
