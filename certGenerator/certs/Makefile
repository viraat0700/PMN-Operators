cert:
	sudo openssl req -nodes -new -newkey rsa:4096 -keyout ca.key -x509 -sha256 -days 365 -out ca.crt -config ca.cnf

	sudo openssl genrsa -out server.key 4096
	sudo openssl req -new -key server.key -out server.csr -config server.cnf
	sudo openssl req -noout -text -in server.csr
	sudo openssl x509 -req -days 9999 -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -extfile server.cnf -extensions v3_ext
	
	sudo openssl genrsa -out client.key 4096 
	sudo openssl req -new -key client.key -out client.csr -config client.cnf 
	sudo openssl req -noout -text -in client.csr
	sudo openssl x509 -req -days 9999 -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out client.crt -extfile client.cnf -extensions v3_ext

	
	sudo chmod 777 server.crt server.key ca.crt client.crt client.key
