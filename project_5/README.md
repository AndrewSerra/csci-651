# Project 5 - Ping and Traceroute

Name: Andrew Serra

## Instructions to Download GoLang and Execute a Project

1. **Download GoLang**: You can download GoLang from the official Go website. Choose the appropriate version for your operating system and follow the instructions provided.
2. **Install GoLang**: After downloading the GoLang installer, run the installer and follow the prompts. Make sure to set the Go path in your environment variables.
3. **Navigate to project directory**: Navigate to the directory of your Go project using the terminal. Once you're in the project directory, enter the directory for either acs8929_ping or acs8929_traceroute.
4. **Execute the Project**: Execute the command:
    
    `sudo go run acs8929_<program-name> [options] <destination>`
    

If there are missing dependencies, run `go install <package-name>` 

# Ping

The ping tool implemented is under the `acs8929_ping/` directory.

The code runs with the following command:

`sudo go run acs8929_ping [options] <destination>`  

It’s important to run with “sudo” to have priviledges for the raw socket.

There are four options:

- Count (-c) → Defaults to -1. If -1, the loop will run indefinitely. If not, the loop will only ping the destination the provided amount of times.
- Wait (-i) → Defaults to 1. This is the time waited before another ping is sent.
- Packet Size (-s) → Defaults to 8. Changes the packet size being delivered.
- Timeout (-t) → Defaults to 3. If this time limit is exceed after an ICMP Echo is sent, then it is considered a fail due to timeout.

Stats are displayed if there is a count option defined. Otherwise ping continues indefinitely.

There is one positional argument:

- Destination → The ip address or domain that is trying to be reached

## Functioning

For the ping operation, the code creates an ICMP message with a body of Echo. This message is sent through the socket. 

The response to the ICMP Echo is received either before the deadline or is considered a failure. The content of the message does not have an effect.

# Traceroute

The traceroute tool implemented is under the `acs8929_traceroute/` directory.

The code runs with the following command:

`sudo go run acs8929_traceroute [options] <destination>`  

It’s important to run with “sudo” to have priviledges for the raw socket.

There are three options:

- Display Non-symbolic Hop (-n) → Defaults to false. Instead of displaying domains, only use IP addresses.
- Number of probes per TTL test (-q) → Defaults to 3. For each TTL value, the provided amount of times will be how many pings are sent.
- Display stats for TTL probing (-S) → Defaults to false. After all probes are completed for a TTL value, displays how many succeded, failed, and the rate of loss.

There is one positional argument:

- Destination → The ip address or domain that is trying to be reached

## Functioning

Similar to the ping tool, the code creates an ICMP message with a body of Echo. This message is sent through the socket. 

The response to the ICMP Echo is received either before the fixed deadline, nobody is reached,  or is considered a failure. The content of the message does not have an effect.

Depending on the number of probes (-q flag), the same ttl value is sent the amount of times provided. The number of maximum hops is limited and is hardcoded in the tool. 

# Measurements and route changes

Afte recording the output of traceroute over multiple days and multiple different times in a day, there is certainly a change in the route that a packet takes. 

I do not believe the specific traffic to a site causes these changes. The general usage of a set of routers might result in other routes to be used. Either way the routes changed for all sites tested (amazon.com, google.com, netflix.com).
