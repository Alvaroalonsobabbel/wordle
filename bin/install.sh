#!/bin/bash

url=$(curl -s https://api.github.com/repos/Alvaroalonsobabbel/wordle/releases/latest | grep "browser_download_url" | grep "wordle" | sed 's/.*"browser_download_url": "\(.*\)"/\1/')

sudo curl -sSL -o /usr/local/bin/wordle $url
sudo chmod +x /usr/local/bin/wordle
