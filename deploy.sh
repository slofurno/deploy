#!/bin/bash

echo 'tag: $TRAVIS_TAG' | envsubst | cat

