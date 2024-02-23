__author__ = "Andrew Serra"
"""
pktsniffer is a cli that reads in a pcap file
from wireshark and print header contents of the
packets and their layers

Date: 02/03/2024
"""
import sys
from scapy.utils import rdpcap
from scapy.sendrecv import sniff
from scapy.layers.l2 import Ether
from scapy.layers.inet import PacketList, Packet, IP, TCP, UDP, ICMP
from argparse import ArgumentParser
from pathlib import Path

def does_file_exists(filename: str) -> bool:
    return Path(filename).exists()

def print_ether_header(header: Ether):
    prefix = "ETHER:"
    packet_size = len(header)
    dst, src = header.dst, header.src
    t = header.type

    print_txt = f"""
    {prefix} ---- Ether Header ----
    {prefix} 
    {prefix} Packet size = {packet_size} bytes
    {prefix} Destination = {dst}
    {prefix} Source      = {src}
    {prefix} Ethertype   = {t}
    {prefix}"""
    print(print_txt)

def print_ip_header(header: IP):
    prefix = "IP:"
    
    print_txt = f"""
    {prefix} ---- IP Header ----
    {prefix} 
    {prefix} Version          = {header.version}
    {prefix} Header Length    = {header.ihl * 4} bytes
    {prefix} Type of Service  = {hex(header.tos)}
    {prefix} Total Length     = {header.len} bytes
    {prefix} Identification   = {header.id}
    {prefix} Flags            = {header.flags}
    {prefix}    .{int(header.flags == "DF")}.. = do not fragment
    {prefix}    ..{int(header.flags == "MF")}. = more fragments
    {prefix} Fragment Offsets = {header.frag}
    {prefix} Time to live     = {header.ttl} seconds/hops
    {prefix} Protocol         = {header.proto}
    {prefix} Header Checksum  = {hex(header.chksum)}
    {prefix} Source Address   = {header.src}
    {prefix} Dest Address     = {header.dst}
    {prefix} Options          = {", ".join(header.options) if len(header.options) else "no options"}
    {prefix}"""
    print(print_txt)

def print_tcp_header(header: TCP):
    prefix = "TCP:"
    options = ""

    if len(header.options) == 0:
        options = "no options"
    else:
        for option in header.options:
            options += f"\n    {prefix}     {option}"

    print_txt = f"""
    {prefix} ---- TCP Header ----
    {prefix} 
    {prefix} Source Port      = {header.sport}
    {prefix} Destination Port = {header.dport}
    {prefix} Sequence Number  = {header.seq}
    {prefix} Ack Number       = {header.ack}
    {prefix} Data Offset      = {header.dataofs}
    {prefix} Flags            = {header.flags}
    {prefix}    ..{int("U" in header.flags)}. .... = urgent pointer
    {prefix}    ...{int("A" in header.flags)} .... = ack
    {prefix}    .... {int("P" in header.flags)}... = push
    {prefix}    .... .{int("R" in header.flags)}.. = reset
    {prefix}    .... ..{int("S" in header.flags)}. = syn
    {prefix}    .... ...{int("F" in header.flags)} = fin
    {prefix} Window           = {header.window}
    {prefix} Checksum         = {hex(header.chksum)}
    {prefix} Urgent Pointer   = {header.urgptr}
    {prefix} Options          = {options}
    {prefix}"""
    print(print_txt)

def print_udp_header(header: UDP):
    prefix = "UDP:"
    
    print_txt = f"""
    {prefix} ---- UDP Header ----
    {prefix} 
    {prefix} Source Port      = {header.sport}
    {prefix} Destination Port = {header.dport}
    {prefix} Length           = {header.len} bytes
    {prefix} Header Checksum  = {hex(header.chksum)}
    {prefix}"""
    print(print_txt)

def print_icmp_header(header: ICMP):
    prefix = "ICMP:"

    print_txt = f"""
    {prefix} ---- ICMP Header ----
    {prefix} 
    {prefix} Type             = {header.type}
    {prefix} Code             = {header.code}
    {prefix} Header Checksum  = {hex(header.chksum)}
    {prefix}"""
    print(print_txt)

def filter_packets(packets: PacketList, filter_expr: list[str], count: int = 0) -> PacketList:
    return sniff(offline=packets, filter=" ".join(filter_expr), count=count)

def print_pcap_file(filename: str, filter_expr: list[str], count: int = 0) -> None:
    
    if not does_file_exists(filename):
        print(f"File {filename} does not exist.")
        sys.exit(1)

    packets: PacketList = filter_packets(rdpcap(filename), filter_expr, count)

    packet_count = 1
    packet: Packet
    for packet in packets:
        print(f"\nPACKET NUMBER -- {packet_count}")
        packet_count += 1

        if packet.haslayer("Ether"):
            print_ether_header(packet[Ether])

        if packet.haslayer("IP"):
            print_ip_header(packet[IP])

        if packet.haslayer("TCP"):
            print_tcp_header(packet[TCP])

        if packet.haslayer("UDP"):
            print_udp_header(packet[UDP])

        if packet.haslayer("ICMP"):
            print_icmp_header(packet[ICMP])
        

if __name__ == "__main__":
    
    parser = ArgumentParser("Reads a pcap file and prints the packets")

    parser.add_argument(
        "-r",
        required=True,
        dest="read_file",
        help="pcap or pcapng file to read from")
    
    _filter_group = parser.add_argument_group("filter", "pcap file reading filter options")
    _filter_group.add_argument("-c", default=0, type=int)
    _filter_group.add_argument("filter_statement", nargs='*')

    args = parser.parse_args()

    print_pcap_file(args.read_file, args.filter_statement, args.c)
