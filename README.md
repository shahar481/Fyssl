# Fyssl

Welcome to Fyssl! The tool for all your protocol researching needs! 

## How does it work?


The tool listens on given addresses and forwards the traffic to the other. But you can change, and look at the traffic before you transfer it on. 


## Cool features

- Open an editor each time you get a packet, edit it, send it back edited

- Sends packets to a remote socket which you can then make python scripts which will edit them

- Drop wanted packets

- Hexdump wanted packets


## How can I use it?
Take a look at config/example/config.json
for an example config
```
./fyssl -c config.json
```

## Config Fields

### Connection

Tcp connection
```
"name": "example-connection", // Name of connection, used to be identified in logs 
"type": "tcp", // Type of connection (Currently theres only ssl/tcp)
"listen-address": "0.0.0.0:1234", // Address to listen on (will listen to multiple sockets)
"connect-address": "127.0.0.1:4567", // Address to forward connections to
"params": {}, // Not used in regular tcp connection (used in ssl as shown in the ssl example)
"modules": { // Modules are not implemented
  "catch": { // Catch module will use iptables to catch connections from specific processes and ports (not implemented)
    "process": ".",
    "ports": "1234"
  }
}
"actions":{} // Actions are not connection type specific and are explained further down in the readme
```

SSL connection
```
"name": "ssl-example", // Name of connection, used to be identified in logs 
"type":"ssl",// Type of connection (Currently theres only ssl/tcp)
"listen-address": "0.0.0.0:15689", // Address to listen on (will listen to multiple sockets)
"connect-address": "127.0.0.1:1876", // Address to forward connections to
"params": { // Used to specify key and cert path
  "key-path": "/tmp/key", // Key path
  "cert-path": "/tmp/cert" // Cert path
},
"modules": { // Not implemented
  "catch": {
    "process": ".",
    "ports": "100-200"
  }
}
```

### Actions
Actions are what processes the packets that are received by Fyssl. Connections create the sockets, then each packet
will be processed by actions until it matches one of them with regex. Then the action will do its designated target
(drop the packet for example). The actions placement in the list matters! The actions are processed from first to last.
The packet will try being matched by the first actions first, and the last actions last.
Currently there are 5 action types

#### Remote
This action sends the packets to a remote socket then waits for a reply with the changes to the packet.

It will try connecting to two remote sockets (unless the param is empty) and will forward packets to there.

It will wait for a reply with the formatted packet and will forward it.
Will only connect to the sockets once the action has been matched.

Packet structure is as so:
```
[2 bytes length][packet]
```
So to send an empty packet you have to set the length to 0

Config structure:
```
"name": "remote-processor", // Used for logging
"trigger": ".", // Regex to match action
"target": "remote", // What the action will do with the packet
"dump": false, // If it should dump the packet to a log file
"target-params": { // Parameters specific for each target
  "listener-socket-address": "127.0.0.1:1758", // Where to forward packets that have been sent from the listener
  "sender-socket-address": "" // Where to forward packets that have been sent from the sender
}
```

#### Editor
This will open the editor of your choice for each packet matched by the action. Then fyssl will wait for the file
to be saved. When the file is saved with the packet data it will forward it on.
```
"name": "editor-processor",
"trigger": ".",
"target": "editor",
"dump": false,
"target-params": {
  "editor": "/usr/bin/subl", // Editor to open the packet with
  "args": "{file}" // Args to give editor (must have {file} this arg will be replaced with the name of the temp file with the packet data)
}
```

#### Reply
This will send back static replies. Currently only asterisks are implemented which can split the packet and forward
only parts of it

For example:
```
REPLY-PARAMS - *0-4*

SENT:
ASDFASDFASDF

RECEIVED:
ASDFA
```

Structure:
```
"name": "zero-to-seven-replier",
"trigger": ".+abc.+",
"target": "reply",
"dump": true,
"target-params": {
  "message": "*0-7*" // Parameters for the reply, this will forward the first 8 bytes of the packet
}
```

#### Drop
This will just drop the packets and not forward them
```
"name": "asdf-replier",
"trigger": ".+abc.+",
"target": "drop",
"dump": true,
"target-params": {}
```

#### None
This will ignore the packets completely

Uses:
* To save processing time and not process all actions if you know that a specific regex does not need to be processed.
* To dump packets and not do any action with them
```
"name": "asdf-replier",
"trigger": ".+abc.+",
"target": "none",
"dump": true,
"target-params": {}
```