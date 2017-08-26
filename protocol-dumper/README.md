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
