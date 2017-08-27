# ELP IP Camera protocol dumper

Use this as a man in the middle between the ELP closed source Windows-only
software and the camera to dump the communication for offline analysis.

This program will dump and decode the protocol used when the CMS software uses
Vendor = H264DVR.

The `Packet` struct is able to both encode and decode the H264DVR protocol.

Tested with: ELP-IP3100HR

## Usage

    go build
    ./protocol-dumper 192.168.1.10:34567 :34567

The first parameter is the camera IP and port.

The second parameter is the port to listen for connections on.

Configure CMS to connect to your IP and port. The dumper will then forward all
communications between the CMS and the camera, but will also dump it all in the
*dumps* directory.

## Dumps

The files are numbered starting at 1 for each run. Please move your files away
between runs because they'll otherwise be overwritten.

Each message is dumped in three formats: bin, txt, payload.

* Bin: The raw on-the-wire message
* Txt: A text representation of the message for easy inspection
* Payload: the message with the header stripped; i.e. just the data

## What's the purpose of this?

It started when I wanted to configure a camera I got and the only way to do it
was through a closed source Windows-only piece of software from the CD shipped
with the camera. Point of reference: I'm running Linux, I don't have a CD drive,
and even if I could attempt to run Wine or a Windows VM the software doesn't
seem to be available for download anywhere.

I'm hoping that one day somebody (if not me) will write some software that can
do some minimal configuration on these cameras.

Changing the network settings with the default password should be achievable.
To generate a password hash see the `dahua_hash` package in this repository.
