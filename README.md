# InstallSoftwareCounter
Collect the software installed on the Windows, count the collection and write to the csv files.


client:
-h argument is the server IP address, port is 8765.
read installed information (name and version) from the registry (means if there is no information in the registry, no software information will be collect).
Send all collected information to the server by using compressed json string, compress by zlib.

test on windows 7 and windows 10

server:
listen on the 8765 port, recive the compressed json string, get the information and count it. Write the hosts information into hostSoftwareList-xxx.csv, the count data into softwareCount-xxxx.csv.



PS:
This is in the alpha test.
