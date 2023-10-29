This program requires you to run ssl on your PC. 
I am running on a windows machine using this package for from the chocolatey package manager for it: 
https://community.chocolatey.org/packages/openssl

generate key with 

`openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -subj "/CN=localhost" -addext "subjectAltName = DNS:localhost"`

Then run 

`openssl rsa -in key.pem -out key.pem`

Run server from root with go 
Then open the three clients in seperate terminals. 
The IDs are very important. 
Just use whatever port for server, just make sure to tell the clients which one. 

``` bash
go run server/server.go -port 5120
go run client/client.go -id 1 -sport 5120
go run client/client.go -id 2 -sport 5120
go run client/client.go -id 3 -sport 5120
```

