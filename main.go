package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/mailgun/groupcache"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

func getPeerIPs(cname string) []string {
	addr, _ := net.LookupIP(cname)
	peers := make([]string, len(addr))
	for _, a := range addr {
		peers = append(peers, fmt.Sprintf("http://%s:8000", a.String()))
	}
	log.Printf("Peers: %v", peers)
	return peers
}

func main() {
	// NewHTTPPool registers /_groupcache/ with http.DefaultServeMux.
	pool := groupcache.NewHTTPPool(fmt.Sprintf("http://%s:8000", os.Getenv("MY_POD_IP")))

	cfg, _ := rest.InClusterConfig()
	clientset, _ := kubernetes.NewForConfig(cfg)

	selector, _ := fields.ParseSelector("metadata.name=sample")

	watchlist := cache.NewListWatchFromClient(clientset.CoreV1().RESTClient(), "endpoints", "default", selector)
	_, controller := cache.NewInformer(
		watchlist,
		&v1.Endpoints{},
		time.Second*0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				pool.Set(getPeerIPs("sample-headless")...)
			},
			DeleteFunc: func(obj interface{}) {
				pool.Set(getPeerIPs("sample-headless")...)
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				pool.Set(getPeerIPs("sample-headless")...)
			},
		},
	)

	go controller.Run(nil)

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
