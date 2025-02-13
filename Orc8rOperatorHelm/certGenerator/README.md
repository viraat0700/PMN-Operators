Generating Cert And Verify the cert using Golang 

 step 1: Clone the repository ```https://github.com/shatru-patel/mTLS.git```

 step 2: Go to the cert folder and Run this command to generate certificate ```make cert```  
          (Example: ``ubunt@ubunt:~/mTLS/certs$ make cert``)                       
 Output will come like this
          ``ubunt@ubunt:~/mTLS/certs$ ls
          
          ca.cnf  ca.crt  ca.key  client.cnf  client.crt  client.csr  client.key  server.cnf  server.crt  server.csr  server.key
         
 
 step 3 : Now go to server folder and run ```go run main.go```    
         (Example: ``ubunt@ubunt:~/mTLS/server$ go run main.go``)

