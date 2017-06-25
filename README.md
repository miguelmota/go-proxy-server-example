# go proxy server example

> golang proxy server example with TLS client authentication

## Usage

start server

```bash
go build main.go
./main
```

Docker run

```bash
docker build -t go_basic .
docker run -p 8080:8000 go_basic --follow
```

available on `<ip>:8080/proxy`

find docker ip using `docker-machine ip`

---

Creating a new CA

1. create the CA key

```bash
openssl genrsa -out ca.key 1024
```

2. create a certificate signing request

```bash
openssl req -new -key ca.key -out ca.csr
```

3. self-sign the request for the creation of the certificate

```bash
openssl x509 -req -in ca.csr -out ca.crt -signkey ca.key
```

4. check the cert

```bash
openssl x509 -in ca.crt -text
```

Generate a new certificate

First, create private key

```bash
openssl genrsa -out example.com.key 1024
```

Second, create a new certificate signing request with private key

```bash
openssl req -new -key example.com.key -out example.com.csr
```

Lastly, sign certificate signing request with certificate authority private key and cert

```bash
openssl ca -in example.com.csr -cert ca.crt -keyfile ca.key -out example.com.crt
```

Check contents of certificate

```bash
openssl x509 -in example.com.crt -text
```

---

# License

MIT