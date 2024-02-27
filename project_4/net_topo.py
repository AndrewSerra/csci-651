__author__ = "Andrew Serra"
'''

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
        # Add 3 Routers
        r1 = self.addHost('r1', cls=LinuxRouter, ip='20.10.100.1/24')
        r2 = self.addHost('r2', cls=LinuxRouter, ip='20.10.100.2/24')
        # r3 = self.addHost('r3', cls=LinuxRouter, ip='20.10.100.3/24')

        # Add 3 switches
        s1 = self.addSwitch('s1')
        s2 = self.addSwitch('s2')
        # s3 = self.addSwitch('s3')

        # Link Routers to each other 
        self.addLink(r1, r2, 
                     intfName1='r1-eth1', params1={ 'ip': '20.10.100.4/24' },
                     intfName2='r2-eth1', params2={ 'ip': '20.10.100.5/24' })
        # self.addLink(r2, r3,
        #              intfName1='r2-eth2', params1={ 'ip': '20.10.100.6/24' },
        #              intfName2='r3-eth1', params2={ 'ip': '20.10.100.7/24' })
        # self.addLink(r3, r1,
        #              intfName1='r3-eth2', params1={ 'ip': '20.10.100.8/24' },
        #              intfName2='r1-eth2', params2={ 'ip': '20.10.100.9/24' })


        # connect routers to switches
        self.addLink(r1, s1, intfName1='r1-eth3', params1={ 'ip': '20.10.172.1/26' })
        self.addLink(r2, s2, intfName1='r2-eth3', params1={ 'ip': '20.10.172.65/25' })
        # self.addLink(r3, s3, intfName1='r3-eth3', params1={ 'ip': '20.10.172.193/27' })

        # Configure LAN A
        h1 = self.addHost('h1', ip='20.10.172.2/26', defaultRoute='20.10.172.1/26')
        h2 = self.addHost('h2', ip='20.10.172.3/26', defaultRoute='20.10.172.1/26')

        self.addLink(s1, h1)
        self.addLink(s1, h2)

        # Configure LAN B
        h3 = self.addHost('h3', ip='20.10.172.66/25', defaultRoute='20.10.172.65/25')
        h4 = self.addHost('h4', ip='20.10.172.67/25', defaultRoute='20.10.172.65/25')

        self.addLink(s2, h3)
        self.addLink(s2, h4)

        # Configure LAN C
        # h5 = self.addHost('h5', ip='20.10.172.194/27', defaultRoute='20.10.172.193/27')
        # h6 = self.addHost('h6', ip='20.10.172.195/27', defaultRoute='20.10.172.193/27')

        # self.addLink(s3, h5)
        # self.addLink(s3, h6)


if __name__ == "__main__":
    setLogLevel('info')

    topo = NetworkTopography()
    net = Mininet(topo=topo)

    info(net['h1'].cmd("ip route add default via 20.10.172.1"))
    info(net['h2'].cmd("ip route add default via 20.10.172.1"))
    info(net['h3'].cmd("ip route add default via 20.10.172.65"))
    info(net['h4'].cmd("ip route add default via 20.10.172.65"))
    # net['h5'].cmd("ip route add default via 20.10.172.193")
    # net['h6'].cmd("ip route add default via 20.10.172.193")
    

    # info(net['r1'].cmd("ip route add 20.10.172.128/25 via 20.10.101.2"))
    # info(net['r1'].cmd("ip route add 20.10.172.224/27 via 20.10.101.1"))

    # info(net['r2'].cmd("ip route add 20.10.172.192/26 via 20.10.101.1"))
    # # info(net['r2'].cmd("ip route add 20.10.172.224/27 via 20.10.101.2"))

    # net['r3'].cmd("ip route add 20.10.172.128/25 via 20.10.100.7")
    # net['r3'].cmd("ip route add 20.10.172.192/26 via 20.10.100.8")

    net.start()
    CLI(net)
    net.stop()
