package main

import (
  "os"
  "io"
  "io/ioutil"
  "net/http"
  "crypto/tls"
  "crypto/x509"
  "fmt"
  "flag"
  "time"
)

var appPath = os.Getenv("APP_PATH")

// cert locations
var (
  certFile = flag.String("cert", appPath + "certs/client.crt", "A PEM eoncoded certificate file.")
  keyFile = flag.String("key", appPath + "certs/client.key", "A PEM encoded private key file.")
  caFile = flag.String("CA", appPath + "certs/ca.crt", "A PEM eoncoded CA's certificate file.")
  httpClient *http.Client
)

const (
  MaxIdleConnections int = 100
  RequestTimeout int = 3600
)

func createHttpClient() (*http.Client, error) {
  flag.Parse()

  // Load client cert
  cert, err := tls.LoadX509KeyPair(*certFile, *keyFile)

  if err != nil {
    fmt.Println(err)
    return nil, err
  }

  // Load CA cert
  caCert, err := ioutil.ReadFile(*caFile)

  if err != nil {
    fmt.Println(err)
    return nil, err
  }

  caCertPool := x509.NewCertPool()
  caCertPool.AppendCertsFromPEM(caCert)

  // Setup HTTPS client
  tlsConfig := &tls.Config{
    Certificates: []tls.Certificate{cert},
    RootCAs: caCertPool,
    InsecureSkipVerify: true,
  }

  tlsConfig.BuildNameToCertificate()
  transport := &http.Transport{
      TLSClientConfig: tlsConfig,
      MaxIdleConnsPerHost: MaxIdleConnections,
  }
  client := &http.Client{
      Transport: transport,
      Timeout: time.Duration(RequestTimeout) * time.Second,
  }

  return client, nil
}

func PingHandler(w http.ResponseWriter, r *http.Request) {
  io.WriteString(w, "pong")
}

func ProxyHandler(w http.ResponseWriter, r *http.Request) {
  if r.Method == "POST" {

    host := "https://example.com"
    path := r.URL.String()

    url := host + path

    req, err := http.NewRequest("POST", url, r.Body)

    if err != nil {
      fmt.Println(err)
      http.Error(w, "", http.StatusInternalServerError)
      return
    }

    // copy headers to request
    for k, v := range r.Header {
      req.Header.Set(k, v[0])
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Access-Control", "application/json")
    req.Header.Set("Access-Control-Allow-Origin", req.Header.Get("Origin"))
    req.Header.Del("Host")
    req.Header.Del("Content-Length")

    resp, err := httpClient.Do(req)

    if err != nil {
      fmt.Println(err)
      http.Error(w, "", http.StatusInternalServerError)
      return
    }

    // re-use connection
    defer resp.Body.Close()

    // response body
    body, err := ioutil.ReadAll(resp.Body)

    if err != nil {
      fmt.Println(err)
      http.Error(w, "", http.StatusInternalServerError)
      return
    }

    for k, v := range resp.Header {
      w.Header().Set(k, v[0])
    }

    w.Header().Set("Access-Control-Allow-Credentials", "true")
    w.WriteHeader(200)
    w.Write(body)
  } else {
    http.Error(w, "404 Not Found", http.StatusNotFound)
    return
  }
}

func main() {
  _httpClient, err := createHttpClient()

  if err != nil {
    panic(err)
  }

  httpClient = _httpClient

  http.HandleFunc("/ping", PingHandler)
  http.HandleFunc("/proxy", ProxyHandler)

  port := os.Getenv("PORT")
  if port == "" {
    port = "8000"
  }

  host := ":" + port

  err = http.ListenAndServeTLS(host, "certs/server.crt", "certs/server.key", nil)

  if err != nil {
    fmt.Println(err)
  }

  fmt.Println("Listening on port " + port)
}
