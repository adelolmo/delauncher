# delauncher
Send magnet links to a remote deluge server with a click

# How to Install

## Setup apt repository

```
sudo apt-get install apt-transport-https
```

```
wget -O - http://adelolmo.github.io/andoni.delolmo@gmail.com.gpg.key | sudo apt-key add -
echo "deb http://adelolmo.github.io/xenial xenial main" | sudo tee /etc/apt/sources.list.d/adelolmo.list
sudo apt-get update
```

## Install package
```
sudo apt-get install delauncher
```