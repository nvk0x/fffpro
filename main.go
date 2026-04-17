package main

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	workers     = flag.Int("w", 50, "number of workers")
	timeout     = flag.Int("timeout", 25, "request timeout (seconds)")
	retries     = flag.Int("retries", 3, "retry count")
	outputDir   = flag.String("o", "out", "output directory")
	saveAll     = flag.Bool("S", true, "save all responses")
	ignoreEmpty = flag.Bool("ignore-empty", true, "ignore empty responses")
)

var client *http.Client

func main() {
	flag.Parse()

	client = createHTTPClient(time.Duration(*timeout) * time.Second)

	os.MkdirAll(*outputDir, 0755)

	jobs := make(chan string, *workers*2)
	var wg sync.WaitGroup

	for i := 0; i < *workers; i++ {
		wg.Add(1)
		go worker(jobs, &wg)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		u := strings.TrimSpace(scanner.Text())
		if u != "" {
			jobs <- u
		}
	}

	close(jobs)
	wg.Wait()
}

func worker(jobs <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for rawURL := range jobs {
		processURL(rawURL)
	}
}

func processURL(rawURL string) {

	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return
	}

	if !isHostAlive(req.URL.Host) {
		return
	}

	var resp *http.Response

	for i := 0; i < *retries; i++ {
		resp, err = client.Do(req)
		if err == nil {
			break
		}
		time.Sleep(time.Duration(i+1) * time.Second)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "[-] %s\n", rawURL)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if *ignoreEmpty && len(bytes.TrimSpace(body)) == 0 {
		return
	}

	fmt.Printf("[+] %d %s\n", resp.StatusCode, rawURL)

	if *saveAll {
		saveResponse(req, resp, body)
	}
}

func saveResponse(req *http.Request, resp *http.Response, body []byte) {

	hash := sha1.Sum([]byte(req.Method + req.URL.String()))
	name := fmt.Sprintf("%x", hash)

	dir := filepath.Join(*outputDir, req.URL.Hostname())
	os.MkdirAll(dir, 0755)

	bodyPath := filepath.Join(dir, name+".body")
	headerPath := filepath.Join(dir, name+".headers")

	os.WriteFile(bodyPath, body, 0644)

	f, err := os.Create(headerPath)
	if err != nil {
		return
	}
	defer f.Close()

	var b strings.Builder

	b.WriteString(fmt.Sprintf("%s %s\n\n", req.Method, req.URL.String()))
	b.WriteString(fmt.Sprintf("HTTP %d\n", resp.StatusCode))

	for k, v := range resp.Header {
		b.WriteString(fmt.Sprintf("%s: %s\n", k, strings.Join(v, ",")))
	}

	f.WriteString(b.String())
}

func createHTTPClient(timeout time.Duration) *http.Client {

	tr := &http.Transport{
		MaxIdleConns:        500,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
		DisableKeepAlives:   false,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
	}

	return &http.Client{
		Transport: tr,
		Timeout:   timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

func isHostAlive(host string) bool {
	conn, err := net.DialTimeout("tcp", host+":80", 3*time.Second)
	if err != nil {
		conn, err = net.DialTimeout("tcp", host+":443", 3*time.Second)
		if err != nil {
			return false
		}
	}
	conn.Close()
	return true
}