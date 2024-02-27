__author__ = "Andrew Serra"
'''
This is a network simulation containing
Three LAN networks that are connected on
20.10.100.0/24
Each LAN has a subnet with different sizes. Hosts are
connected via a switch.
'''
from mininet.log import setLogLevel, info
from mininet.cli import CLI
from mininet.net import Mininet
from mininet.topo import Topo
from mininet.node import Node

class LinuxRouter(Node):
    def config(self, **params):
        super(LinuxRouter, self).config(**params)
        self.cmd('sysctl net.ipv4.ip_forward=1')
        
    def terminate(self):
        self.cmd('sysctl net.ipv4.ip_forward=0')
        super(LinuxRouter, self).terminate()

class NetworkTopography(Topo):
    def build(self, *args, **params):
        r1 = self.addHost('r1', cls=LinuxRouter, ip='20.10.172.129/26')
        r2 = self.addHost('r2', cls=LinuxRouter, ip='20.10.172.1/25')
        r3 = self.addHost('r3', cls=LinuxRouter, ip='20.10.172.193/27')

        # Add 2 switches
        s1 = self.addSwitch('s1')
        s2 = self.addSwitch('s2')
        s3 = self.addSwitch('s3')

        # Add host-switch links in the same subnet
        self.addLink(s1,
                     r1,
                     intfName2='r1-eth1',
                     params2={'ip': '20.10.172.129/26'})

        self.addLink(s2,
                     r2,
                     intfName2='r2-eth1',
                     params2={'ip': '20.10.172.1/25'})
        
        self.addLink(s3,
                     r3,
                     intfName2='r3-eth1',
                     params2={'ip': '20.10.172.193/27'})

        # Add router-router link in a new subnet for the router-router connection
        self.addLink(r1,
                     r2,
                     intfName1='r1-eth2',
                     intfName2='r2-eth2',
                     params1={'ip': '20.10.100.1/24'},
                     params2={'ip': '20.10.100.2/24'})

        self.addLink(r2,
                     r3,
                     intfName1='r2-eth3',
                     intfName2='r3-eth2',
                     params1={'ip': '20.10.100.3/24'},
                     params2={'ip': '20.10.100.4/24'})
        
        self.addLink(r3,
                     r1,
                     intfName1='r3-eth3',
                     intfName2='r1-eth3',
                     params1={'ip': '20.10.100.5/24'},
                     params2={'ip': '20.10.100.6/24'})

        # Adding hosts specifying the default route
        h1 = self.addHost(name='h1',
                          ip='20.10.172.130/26',
                          defaultRoute='via 20.10.172.129')
        h2 = self.addHost(name='h2',
                          ip='20.10.172.131/26',
                          defaultRoute='via 20.10.172.129')
        h3 = self.addHost(name='h3',
                          ip='20.10.172.2/25',
                          defaultRoute='via 20.10.172.1')
        h4 = self.addHost(name='h4',
                          ip='20.10.172.3/25',
                          defaultRoute='via 20.10.172.1')
        h5 = self.addHost(name='h5',
                          ip='20.10.172.194/27',
                          defaultRoute='via 20.10.172.193')
        h6 = self.addHost(name='h6',
                          ip='20.10.172.195/27',
                          defaultRoute='via 20.10.172.193')

        # Add host-switch links
        self.addLink(h1, s1)
        self.addLink(h2, s1)
        self.addLink(h3, s2)
        self.addLink(h4, s2)
        self.addLink(h5, s3)
        self.addLink(h6, s3)
        

if __name__ == "__main__":
    setLogLevel('info')

    topo = NetworkTopography()
    net = Mininet(topo=topo)

    info(net['r1'].cmd("ip route add 20.10.172.0/25 via 20.10.100.2 dev r1-eth2"))   # B
    info(net['r1'].cmd("ip route add 20.10.172.192/27 via 20.10.100.5 dev r1-eth3")) # C

    info(net['r2'].cmd("ip route add 20.10.172.128/26 via 20.10.100.1 dev r2-eth2")) # A
    info(net['r2'].cmd("ip route add 20.10.172.192/27 via 20.10.100.4 dev r2-eth3")) # C

    info(net['r3'].cmd("ip route add 20.10.172.0/25 via 20.10.100.3 dev r3-eth2"))   # B
    info(net['r3'].cmd("ip route add 20.10.172.128/26 via 20.10.100.6 dev r3-eth3")) # A

    net.start()
    CLI(net)
    net.stop()
