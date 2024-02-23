__author__ = "Andrew Serra"
'''
    Author: Andrew Serra
    Date: 02/14/2024
'''
import logging
from mininet.cli import CLI
from mininet.log import setLogLevel
from mininet.link import TCLink
from mininet.net import Mininet, Host

def _add_host(net: Mininet, hostname: str, mac: str) -> Host:
    return net.addHost(hostname, mac=mac)

def _add_link(net: Mininet, node1: Host, node2: Host) -> None:
    net.addLink(node1, node2)

def add_hosts(net: Mininet, h: list[tuple[str, str]]) -> list[Host]:
    hosts = []
    for hostname, mac, _ in h:
        host = _add_host(net, hostname, mac)
        hosts.append(host)
    return hosts

def add_links(net: Mininet, connections: list[tuple[Host]]) -> None:
    for n1, n2 in connections:
        _add_link(net, n1, n2)

if __name__ == "__main__":
    setLogLevel("info")
    global logger
    logger = logging.getLogger()

    net = Mininet(link=TCLink)
    hosts_data = [
        ("h1", "00:00:00:00:01:00", "50.0.0.1"),
        ("h2", "00:00:00:00:02:00", "50.0.0.2"),
        ("h3", "00:00:00:00:03:00", "50.0.0.3"),
        ("h4", "00:00:00:00:04:00", "50.0.0.4"),
        ("h5", "00:00:00:00:05:00", None), # bridge
    ]

    hosts = add_hosts(net, hosts_data)
    links = [(hosts[i], hosts[-1]) for i in range(len(hosts)-1)]
    add_links(net, links)

    net.build()

    h5 = hosts[-1] # get the bridge

    # network interfaces
    h5.cmd("ifconfig h5-eth0 0")
    h5.cmd("ifconfig h5-eth1 0")
    h5.cmd("ifconfig h5-eth2 0")
    h5.cmd("ifconfig h5-eth3 0")

    # set ip addresses
    for i in range(len(hosts)):
        hostname, _, ip =hosts_data[i]
        if ip is not None:
            hosts[i].setIP(ip)

    # create bridge
    h5.cmd("brctl addbr br0")

    # connect to bridge
    h5.cmd("brctl addif br0 h5-eth0")
    h5.cmd("brctl addif br0 h5-eth1")
    h5.cmd("brctl addif br0 h5-eth2")
    h5.cmd("brctl addif br0 h5-eth3")

    # start bridge
    h5.cmd("ifconfig br0 up")

    # test_hosts_pinging_each_other(hosts[:-1])

    CLI(net)
    net.stop()
