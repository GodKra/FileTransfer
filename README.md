# ꜰɪʟᴇ ᴛʀᴀɴsꜰᴇʀ

A simple client and server used for sending files. give argument "-help" for availabe flags. Now you must set a 
```type``` flag to send a file. You have to give either "--type" "downloader" or "sender". Downloader starts a
new server which will accept any data sent to it and saves it. Sender sends data to the downloader. you can change the IP which sender uses using the ```ip``` flag.

**USES HIGH DISK SPACE**

### How it works:
* Sender `Zip`s any file selected using ```path``` flag.
* Sender sends the zipped file size then file to the Downloader.
* Downloader reads the size and makes a progress bar.
* Downloader saves the read zipped file.
