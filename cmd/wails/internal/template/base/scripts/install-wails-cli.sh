#! /bin/bash

echo -e "Start running the script..."
cd ../

echo -e "Current Go version: \c"
go version

echo -e "Install the Wails command line tool..."
go install github.com/unix-world/wails-app/cmd/wails@latest

echo -e "Successful installation!"

echo -e "End running the script!"
