#!/bin/bash

FILE_NAME=`uuidgen`.txt
MAX_COMMIT_COUNT=100
DATES=( {{ .Dates }} )

touch $FILE_NAME
for DATE in "${DATES[@]}"; do
    for ((i = 0; i < $MAX_COMMIT_COUNT; i++)); do
        echo `date` >> $FILE_NAME
        git add $FILE_NAME
        git commit -m "🤪 bump commit 😬" --date $DATE
    done
done
