#!/bin/bash

set -e

# Test data installation

BASEPATH=$(dirname $0)
GITPATH="$BASEPATH/testdata/git"
GITREPO="$GITPATH/.git"
GITURL="https://kalyabin@bitbucket.org/kalyabin/yii2-git-view-testing.git"
HGPATH="$BASEPATH/testdata/hg"
HGREPO="$HGPATH/.hg"
HGURL="https://kalyabin@bitbucket.org/kalyabin/yii2-hg-view-testing"

## Installing GIT Repository
if [ ! -d "$GITPATH" ]; then
	echo "Create git repository path: $GITPATH\n"
	mkdir $GITPATH
fi

if [ ! -d "$GITREPO" ]; then
	echo "Clone git repository to: $GITPATH\n"
	git clone $GITURL $GITPATH
fi

## Installing Mercurial Repository
if [ ! -d "$HGPATH" ]; then
	echo "Create Mercurial repository path: $HGPATH\n"
	mkdir $HGPATH
fi

if [ ! -d "$HGREPO" ]; then
	echo "Clone Mercurial repository to: $HGPATH\n"
	hg clone $HGURL $HGPATH
fi
