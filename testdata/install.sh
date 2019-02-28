#!/bin/bash

set -e

# Test data installation

BASEPATH=$(dirname $0)
TESTDATAPATH="$BASEPATH"

# repositories commands
commands=(git hg)
paths=("$TESTDATAPATH/git" "$TESTDATAPATH/hg")
repos=("$TESTDATAPATH/git/.git" "$TESTDATAPATH/hg/.hg")
urls=("https://kalyabin@bitbucket.org/kalyabin/yii2-git-view-testing.git" "https://kalyabin@bitbucket.org/kalyabin/yii2-hg-view-testing")
fakes=("$TESTDATAPATH/.git" "$TESTDATAPATH/.hg")

for i in ${!commands[@]}; do
	PROJECTPATH="${paths[$i]}"
	REPOPATH="${repos[$i]}"
	URL="${urls[$i]}"
	FAKE="${fakes[$i]}"
	
	CLONE_CMD=""
	NAME=""
	if [ "${commands[$i]}" == "git" ]; then
		CLONE_CMD="git clone $URL $PROJECTPATH"
		NAME="GIT"
	fi

	if [ "${commands[$i]}" == "hg" ]; then
		CLONE_CMD="hg clone $URL $PROJECTPATH"
		NAME="Mercurial"
	fi

	## Installing project path
	if [ ! -d "$PROJECTPATH" ]; then
		echo "Create $NAME repository path: $PORJECTPATH"
		mkdir $PROJECTPATH
	fi

	## Cloning repository
	if [ ! -d "$REPOPATH" ]; then
		echo "Clone $NAME repository to: $PROJECTPATH"
		eval $CLONE_CMD
	fi

	if [ -e $FAKE ] && [ -d $FAKE ]; then
		printf "$FAKE should be a file.\nPlease remove directory and re-install testing data.\n"
		exit 1
	fi

	if [ ! -e $FAKE ]; then
		echo "Install fake $NAME repository to: $FAKE"
		touch $FAKE
	fi
done

## Test git branches in travis
cd testdata/git && git branch -a -v && git --no-pager log --format=%H%n%P%n%an%n%ae%n%ad%n%s -n 1 --skip=0 --branches=*branch1*
