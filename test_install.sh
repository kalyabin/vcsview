#!/bin/bash

set -e

# Test data installation

BASEPATH=$(dirname $0)
NOREPOPATH="$BASEPATH/testdata"
FAKEGITPATH="$NOREPOPATH/.git"
FAKEHGPATH="$NOREPOPATH/.hg"
GITPATH="$BASEPATH/testdata/git"
GITREPO="$GITPATH/.git"
GITURL="https://kalyabin@bitbucket.org/kalyabin/yii2-git-view-testing.git"
HGPATH="$BASEPATH/testdata/hg"
HGREPO="$HGPATH/.hg"
HGURL="https://kalyabin@bitbucket.org/kalyabin/yii2-hg-view-testing"

## Installing GIT Repository
if [ ! -d "$GITPATH" ]; then
	echo "Create git repository path: $GITPATH"
	mkdir $GITPATH
fi

if [ ! -d "$GITREPO" ]; then
	echo "Clone git repository to: $GITPATH"
	git clone $GITURL $GITPATH
fi

## Installing Mercurial Repository
if [ ! -d "$HGPATH" ]; then
	echo "Create Mercurial repository path: $HGPATH"
	mkdir $HGPATH
fi

if [ ! -d "$HGREPO" ]; then
	echo "Clone Mercurial repository to: $HGPATH"
	hg clone --insecure $HGURL $HGPATH
fi

## Installing empty files to check fake repositories
if [ -e "$FAKEGITPATH" ] && [ -d "$FAKEGITPATH" ]; then
	printf "$FAKEGITPATH should be a file.\nPlease remove directory and re-install testing data.\n"
	exit 1
fi
if [ -e "$FAKEHGPATH" ] && [ -d "$FAKEHGPATH" ]; then
	printf "$FAKEHGPATH should be a file.\nPlease remove directory and re-install testing data.\n"
	exit 1
fi

if [ ! -e "$FAKEGITPATH" ]; then
	echo "Install fake Git repository to: $FAKEGITPATH"
	touch $FAKEGITPATH
fi
if [ ! -e "$FAKEHGPATH" ]; then
	echo "Install fake Mercurial repository to: $FAKEHGPATH"
	touch $FAKEHGPATH
fi
