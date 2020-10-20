#### The Global pcap-savefile Header
##### What’s the magic number? What does it tell you about the byte ordering in the pcap-specific aspects of the file?
d4c3 b2a1. Per the pcap-savefile docs, this means my machine is opposite byte order as the machine that produced this packet. My machine uses little endian, so this was big endian. It also means I'll need to reverse byte ordering upon parsing the file.
##### What are the major and minor versions? Don’t forget about the byte ordering!
0200 0400. If our machines had the same byte ordering, this would mean we're on version 512.1024, but since we reverse these byte orderings, this is 2.4.
##### Are the values that ought to be zero in fact zero?
Yes - the docs say the timezone offset and timezone accuracy are always zero - both are zero.
##### What is the snapshot length?
The packet value is `ea05 0000` but after reversing the byte order we have `0000 05ea`, or 1514 in decimal. 
##### What is the link layer header type?
This value is `0100 0000`, which after reversing the byte order is `0000 0001`, and per `man 7 pcap-linktype` and http://www.tcpdump.org/linktypes.html, is link layer header type ethernet.

#### Per-packet Headers
##### What is the size of the first packet?
`4e00 0000` which is 78 in decimal.
##### Was any data truncated?
No. The un-truncated length of the packet is also `4e00 0000`

