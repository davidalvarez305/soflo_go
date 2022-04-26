#!/bin/bash

cd
cd soflo_go
git add .

echo What changes did you make?
read CHANGES

git commit -m "$CHANGES"

git push -u origin main

echo Finished pushing files.