package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/mailgun/groupcache"
)

func main() {
	// NewHTTPPool registers /_groupcache/ with http.DefaultServeMux.
	pool := groupcache.NewHTTPPool(fmt.Sprintf("http://%s:8000", os.Getenv("MY_POD_IP")))

	go func() {
		// Infinite loop every few seconds. In reality, should use k8s watch API
		// to get notified of when `v1.Endpoints` changes

		// Sample subscription to endpoints:
		// https://github.com/mailgun/gubernator/blob/1e6849ab820232acfd31440a33580496fb3d3f45/kubernetes.go

		// With k8s rbac to hit the control plane:
		// https://github.com/mailgun/gubernator/blob/1e6849ab820232acfd31440a33580496fb3d3f45/deploy/helm/templates/rbac.yaml
		for {
			// For now, just query the headless dns record over and over again
			time.Sleep(500 * time.Millisecond)

			addr, _ := net.LookupIP("sample-headless")
			peers := make([]string, len(addr))
			for _, a := range addr {
				peers = append(peers, fmt.Sprintf("http://%s:8000", a.String()))
			}

			pool.Set(peers...)
		}
	}()

	// 1024 * 1024 bytes
	group := groupcache.NewGroup("sample", 1024*1024, groupcache.GetterFunc(
		func(ctx groupcache.Context, key string, dest groupcache.Sink) error {
			oneMinuteFromNow := time.Now().Add(time.Minute)
			dest.SetBytes([]byte(fmt.Sprintf("The key \"%s\" was retrieved from peer \"%s\"", key, os.Getenv("MY_POD_IP"))), oneMinuteFromNow)
			return nil
		}))

	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			key := r.URL.Path[1:]
			bytes := []byte{}
			group.Get(r.Context(), key, groupcache.AllocatingByteSliceSink(&bytes))
			fmt.Fprintf(w, "%s\n%#v\n\n", bytes, group.Stats)
		}
	}))

	// Start a HTTP server to listen for peer requests from the groupcache
	log.Printf("Serving....\n")
	http.ListenAndServe(":8000", nil)
}
