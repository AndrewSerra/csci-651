# Description

The file outputs all the connections that are made which
simulate the `netstat` command. It will check the following
socket types:
 - inet4
 - inet6
 - tcp4
 - tcp6
 - udp4
 - udp6

Each socket type is printed separately. For each connection, 
the following is printed:
 - local ip (v4 and v6)
 - local port
 - remote ip
 - remote port
 - process id
 - connection status
 - file descriptor

# Python Version

The python version used when developing is 3.11. It is not tested, but there are no packages that should prevent it from running on any version of python >=3.7

# How to run

The only required package to run is `psutil`. To install the package run:
```
pip3 install psutil
```

Go to the project folder. Run the program, run:
```
sudo python3 route.py
```

There are no options available. It will only display connections.

# Images

Images are located in the project folder under the images directory.
