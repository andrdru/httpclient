package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/andrdru/httpclient"
)

func main() {
	var service = os.Getenv("SERVICE_NAME")
	if service == "" {
		log.Fatalf("no SERVICE_NAME env variable")
	}
	var host = os.Getenv("SERVICE_HOST")
	if host == "" {
		log.Fatalf("no SERVICE_HOST env variable")
	}

	http.HandleFunc("/", handler)

	_ = http.ListenAndServe(host, nil)
}

func handler(writer http.ResponseWriter, request *http.Request) {
	var service = os.Getenv("SERVICE_NAME")

	var xServiceReplace = request.Header.Get("X-Service-Replace")
	var replaces = strings.Split(xServiceReplace, ";")
	var serviceMap = make(map[string]string)
	for _, r := range replaces {
		var parts = strings.Split(r, "=")
		if len(parts) < 2 {
			continue
		}

		serviceMap[parts[0]] = parts[1]
	}

	var h = http.Header{}
	h.Set("X-Service-Replace", xServiceReplace)

	var ctx = context.Background()

	var req []byte
	if service == "aaa" {
		if serviceMap["bbb"] != "" {
			ctx = httpclient.HostReplaceCtx(ctx, serviceMap["bbb"])
		}

		var c = httpclient.NewHttpClient(
			":8081",
			httpclient.Scheme("http"),
			httpclient.Log(log.Default()),
			httpclient.LogLevel(httpclient.LoggerLevelInfo),
		)
		var err error
		_, req, err = c.Request(ctx, http.MethodGet, "/", []byte{}, nil)
		if err != nil {
			log.Printf("error request %s", err)
		}
	}

	if service == "bbb" {
		if serviceMap["ccc"] != "" {
			ctx = httpclient.HostReplaceCtx(ctx, serviceMap["ccc"])
		}

		var c = httpclient.NewHttpClient(
			":8082",
			httpclient.Scheme("http"),
			httpclient.Log(log.Default()),
			httpclient.LogLevel(httpclient.LoggerLevelInfo),
		)
		var err error
		_, req, err = c.Request(ctx, http.MethodGet, "/", []byte{}, h)
		if err != nil {
			log.Printf("error request %s", err)
		}
	}

	_, _ = writer.Write([]byte("inject local service " + service))

	if req != nil {
		_, _ = writer.Write([]byte(":"))
		_, _ = writer.Write(req)
	}
}
