# README #

Timetrack is a go-based programm that offers a slim, staight-forward web-interface to track employee work hours.

## Technology ##
* golang -> can be cross compiled to many platforms
* postgresql

## Install/Build ##
Before you can use/build this project you need Go (golang.org) installed and your GOPATH set up.
In order to get this project you can either clone this repository from github or call `go get github.com/dahawk/timetrack`.

To run this software you need:
* `go get` within the repo to fetch all required dependencies (first time only)
* to create a postgresql database according to the file create.sql
* go to main.go file and set the correct values for
  * username
  * password
  * db host
  * db name
* call `go build`
* run the resulting executable file on bash, cmd, or other (depending on the OS)

The software will start without console output and will listen on port 1234 (unless changed in main.go)

## Possible Future Work ##
* Impersonate (admin can create/edit/delete entries of other users)
* basic logging (store for each entry who modified it)
* extended logging (log table that give useful info who did what on which data when)
* include expected work time per user
* based on expected work time calculate over/under performance
* import holidays to update expedted work time

## Developer/Maintainer ##

Christof Horschitz (horschitz@gmail.com)
