# Description

The program reads in a pcap file generated from Wireshark, then prints
the contents of each packet's header to std out. If there is a count limit
set using the '-c' flag, it will be limited. Filter expressions follow 
[pcap-filter](https://www.tcpdump.org/manpages/pcap-filter.7.html).

Author: Andrew Serra

# Python Version

The python version used when developing is 3.11. It is not tested, but there are no packages that should prevent it from running on any version of python >=3.7

# How to run

The only required package to run is `scapy`. To install the package run:
```
pip3 install scapy
```

Go to the project folder. Run the program, run:
```
python3 pktsniffer.py -r <file-absolute-path> [-c COUNT] [FILTER EXPRESSION]
```

The filter expressions follow [pcap-filter](https://www.tcpdump.org/manpages/pcap-filter.7.html).
