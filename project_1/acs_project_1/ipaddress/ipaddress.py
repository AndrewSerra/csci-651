__author__ = "acs8929"

import psutil

def list_interfaces():
    addrs = psutil.net_if_addrs()
    ifaces = psutil.net_if_stats()

    for iface, data in addrs.items():
        print("Interface")
        print(f"\tName: {iface}")

        for addr in data:
            print(f"\tFamily: {addr.family}")
            print(f"\tMAC Address: {addr.address}")
            print(f"\tNetMask: {addr.netmask}")
            print(f"\t[FLAG] Point To Point: {addr.ptp}")
            print(f"\t[FLAG] Broadcast: {addr.broadcast}")
            print(f"\tMaximum Transmission Unit: {ifaces.get(iface).mtu}")
            print()


if __name__ == "__main__":
    list_interfaces()
