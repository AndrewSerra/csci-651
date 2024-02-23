# Description

The file outputs all the connections that are made which
simulate the `ifconfig` command.

The following information is printed for each of the interfaces
identified from the networick interface controller:
 - interface name
 - family of the interface
 - hardware (MAC) address
 - netmask
 - flags (point to point, broadcast)
 - MTU (maximum transmission unit)

# Python Version

The python version used when developing is 3.11. It is not tested, but there are no packages that should prevent it from running on any version of python >=3.7

# How to run

The only required package to run is `psutil`. To install the package run:
```
pip3 install psutil
```

Go to the project folder. Run the program, run:
```
python3 ipaddress.py
```

There are no options available. It will only display interfaces and addresses.

# Images

Images are located in the project folder under the images directory.
