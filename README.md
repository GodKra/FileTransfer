# ꜰɪʟᴇ ᴛʀᴀɴsꜰᴇʀ

A simple program used for sending files and directories. give argument "-help" for availabe flags. This program consists of two things 
which does different jobs. They are "sender" and "reciever". Reciever starts a new server which will accept any data sent to
it and saves it. Sender sends data to the Reciever. you can change the IP which sender uses using the ```ip``` flag (0 for automatic ip finding). 
You should set a  ```type``` flag to switch between sending and recieving a file. Default state of the flag is "reciever" 
(because you need to create a server before connecting to it). You have to give either "reciever" or "sender" to the ```type``` flag.
The ```ip``` flag's default is 0. If ip flag is in the default state, it will attempt to search for a reciever server automatically.
When the server responds, it will read from the server waiting for a response. If it recieves "\`", then it will decide it is a 
reciever server and continue sending. Giving an IP flag will still make sure it is a downloader server.

**USES HIGH DISK SPACE**

### How it works:
* Sender `Zip`s any file selected (```path``` flag).
* Sender sends the zipped file size then file to the Reciever.
* Reciever reads the size and makes a progress bar.
* Reciever saves the read zipped file.

## Installation
  You can install this program from the '[Releases](https://github.com/GodKra/FileTransfer/releases/latest "Latest Release")' tab
## CommandLine Flags
```
Usage: filetransfer <flags>
  
  Available Flags:
    --filePath [value]: Path of the file you want to transfer. Must for Sender
    --ip [value]:       The IP of the downloader you want to send the file to. 0 for automatic. Optional for Sender
    --fileName [value]: Name to use when saving the recieved files. Optional for Downloader
    --type [value]:     The type of filetransfer. 'sender' to send files. 'reciever' to recieve files
    --help:             Prints this.
```
