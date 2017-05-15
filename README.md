# ꜰɪʟᴇ ᴛʀᴀɴsꜰᴇʀ

A simple program used for sending files and folders. give argument "-help" for availabe flags. This program consists of two things 
which does different jobs. They are "sender" and "downloader". Downloader starts a new server which will accept any data sent to
it and saves it. Sender sends data to the downloader. you can change the IP which sender uses using the ```ip``` flag. 
You should set a  ```type``` flag to send/download a file. Default state of the flag is "downloader" (because you need to create a 
server before connecting to it). You have to give either "--type" "downloader" or "sender".

**USES HIGH DISK SPACE**

### How it works:
* Sender `Zip`s any file selected using ```path``` flag.
* Sender sends the zipped file size then file to the Downloader.
* Downloader reads the size and makes a progress bar.
* Downloader saves the read zipped file.

## Installation
  You can install this program from the '[Releases](https://github.com/GodKra/FileTransfer/releases/latest "Latest Release")' tab
