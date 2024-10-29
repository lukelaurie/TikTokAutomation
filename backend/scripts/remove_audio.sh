#!/bin/bash 

BASE_DIR="../assets/video"

if [ "$#" -lt 1 ]; then 
    echo "Please pass the name of the folder you want to remove the audio from"
    exit 1
fi

FOLDER_NAME="$1"
VIDEO_DIR="$BASE_DIR/$FOLDER_NAME"

cd "$VIDEO_DIR" || { echo "Directory '$VIDEO_DIR' was not found"; exit 1; }

for file in *.mp4; do 
    ffmpeg -i "$file" -c:v copy -an -y "temp.mp4" > /dev/null 2>&1
    rm "$file"
    mv "temp.mp4" "$file"

    echo "Removed audio for '$file'"
done