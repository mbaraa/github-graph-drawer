#!/bin/bash

FILE_NAME=`uuidgen`.txt
DATES=( {{ .Dates }} )

touch $FILE_NAME
for DATE in "${DATES[@]}"; do
    echo `date` >> $FILE_NAME
    git add $FILE_NAME
    git commit -m "🤪 bump commit 😬" --date $DATE
done
