#!/usr/bin/bash

rm -rf test
rm -rf dst

mkdir test
mkdir dst
echo test >dst/a
touch test/a
touch test/b
echo This is a test >test/b
touch test/c.txt
mkdir test/test1/
touch test/test1/d

