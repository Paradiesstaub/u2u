# Ubuntu 2 USB

wip


# Setup on Ubuntu

### 1. Needed PPA’s and Packages

The Atom and Ubuntu-Make PPA are optional.

    sudo add-apt-repository ppa:ubuntu-desktop/ubuntu-make
    sudo add-apt-repository ppa:ubuntu-sdk-team/ppa
    sudo add-apt-repository ppa:webupd8team/atom

    sudo apt-get update

    sudo apt-get install git ubuntu-make ubuntu-sdk atom build-essential qtdeclarative5-dev qtbase5-private-dev qtdeclarative5-private-dev libqt5opengl5-dev qtdeclarative5-qtquick2-plugin pkg-config

### 2. Go installation

    umake go

Then create a golang source-code folder:

    mkdir ~/go

Next edit the .profile to set the mandatory $GOPATH variable:

       gedit ~/.profile

Add the following lines to the .profile file:

    if [ -d "$HOME/go" ] ; then
        export GOPATH="$HOME/go"
    fi

Logout of your current shell session and back again or reboot the computer. Finally check if go was successfully installed:

    go version # should return 1.5.2


### 3. Setup go-qml and get the ubuntu2usb code

    go get gopkg.in/qml.v1
    go get github.com/Paradiesstaub/u2u/tree/master

### 4. (optional) Atom Plugins & Configuration

The Atom editor can be used to write go code and Markdown documents.
Please go to Edit > Preferences and change the ‘Tab Length’ to 4.

    apm install go-plus goto script minimap highlight-selected minimap-highlight-selected
