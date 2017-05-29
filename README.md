# ꜰɪʟᴇ ᴛʀᴀɴsꜰᴇʀ

A simple program used for sending files and directories. give argument "-help" for availabe flags. This program consists of two things 
which does different jobs. They are "sender" and "reciever". 
### Reciever 
  Starts a new server which will accept any data sent to it and saves it. 
### Sender 
  Sends data to the Reciever. you can change the IP which sender uses using the ```ip``` flag (0 for automatic ip finding). 
**USES HIGH DISK SPACE**

## Important Flags
  #### Type (Optional):
  Switches between sending and recieving a file. Default state of the type is "reciever" 
  (because you need to create a server before connecting to it). You have to give either "reciever" or "sender" to the `type` flag. 
  #### FilePath:
  Path of the file you want to transfer. If you have set the `filePath` flag, Then the program will decide you want to send a file. So the type will be changed to Sender. 
  #### IP:
  The IP of the Reciever you want to send the file to. The `ip` flag's default is 0. If ip flag is in the default state, it will attempt to search for a reciever     server automatically.
  When it connects to a server, it will wait for a response from the server. If it recieves "\`", then it will decide it is a Reciever server and continue sending. Setting the IP flag will also do the same. 
  #### Port: 
  The port of the Reciever server. `port` flag should be used if there are many people using the program in your localhost. If the port is the same in all programs,
  it could connect you to a wrong Reciever server. Port flag can also let you make run multiple Reciever servers in one PC. 
  Default value of `port` is '7084'. Which means 'FT'. 70 stands for 'F' and 84 stands for 'T'. FT = FileTransfer
  ##### Help for other flags


### How it works:
* Sender \`Zip\`s any file selected (`path` flag).
* Sender sends the zipped file size then the file to the Reciever.
* Reciever reads the size and makes a progress bar.
* Reciever saves the read zipped file.

## Installation
  You can install this program from the '[Releases](https://github.com/GodKra/FileTransfer/releases/latest "Latest Release")' tab

## CommandLine Flags
```
Usage: filetransfer <flags>

Available Flags:
        --filePath [value]: Path of the file you want to transfer. Must for Sender. If this exists, type will be considered as a sender
        --ip [value]:       The IP of the Reciever you want to send the file to. 0 for automatic. Optional for Sender
        --saveName [value]: Name to use when saving the recieved files. Optional for Reciever
        --type [value]:     Optional. The type of filetransfer. 'sender' to send files. 'reciever' to recieve files
        --help:             Prints this.
        --port [value]:     The port of the Reciever. Default is '7084'. Optional for both Sender and Reciever
```
