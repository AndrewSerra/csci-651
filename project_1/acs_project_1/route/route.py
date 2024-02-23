__author__ = "acs8929"

import psutil

def get_connections():

    for kind in ["inet4", "inet6", "tcp4", "tcp6", "udp4", "udp6"]:
        nc = psutil.net_connections(kind=kind)

        print(kind.capitalize() + " Connections:")

        for conn in nc:
            if isinstance(conn.laddr, tuple) and len(conn.laddr):
                print(f"\tLocal IP: {conn.laddr[0] if conn.laddr is not None else None}")
                print(f"\tLocal Port: {conn.laddr[1] if conn.laddr is not None else None}")
            else:
                print(f"\tLocal Address: {conn.laddr}")

            if isinstance(conn.raddr, tuple) and len(conn.raddr):
                print(f"\tRemote IP: {conn.raddr[0] if conn.raddr is not None else None}")
                print(f"\tRemote Port: {conn.raddr[1] if conn.raddr is not None else None}")
            else:
                print(f"\tRemote Address: {conn.raddr}")

            print(f"\tProcess ID: {conn.pid}")
            print(f"\tStatus: {conn.status}")
            print(f"\tFile Descriptor: {conn.fd}")
            print()

if __name__ == "__main__":
    get_connections()
