package main

import (
	"bufio"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

var (
	fLogin    = flag.String("login", "", "login")
	fPassword = flag.String("password", "", "password")
)

func main() {
	flag.Parse()
	log.SetFlags(0)
	if flag.NArg() < 1 || *fLogin == "" || *fPassword == "" {
		flag.Usage()
		os.Exit(0)
	}
	contentLength, md5sum, sha256sum := fileInfo(flag.Arg(0))

	var tmpname [16]byte
	if _, err := io.ReadFull(rand.Reader, tmpname[:]); err != nil {
		log.Fatal(err)
	}
	// Go's http.Client doesn't want to set Content-Length header with
	// empty Body, so do low-level request.
	conn, err := tls.Dial("tcp", "webdav.yandex.ru:443", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	auth := base64.StdEncoding.EncodeToString([]byte(*fLogin + ":" + *fPassword))
	req := fmt.Sprintf("PUT /%x.tmp HTTP/1.1\r\nHost: webdav.yandex.ru\r\nAccept: */*\r\nAuthorization: Basic %s\r\nEtag: %x\r\nSha256: %x\r\nExpect: 100-continue\r\nContent-Type: application/binary\r\nContent-Length: %d\r\n\r\n", tmpname[:], auth, md5sum, sha256sum, contentLength)
	_, err = io.WriteString(conn, req)
	if err != nil {
		log.Fatal(err)
	}
	r := bufio.NewReader(conn)
	line, err := r.ReadString('\r')
	if err != nil {
		log.Fatal(err)
	}
	switch strings.TrimSpace(line) {
	case "HTTP/1.1 201 Created":
		log.Println("\n˙ ͜ʟ˙\nFILE EXISTS on Yandex.Disk!\n")
	case "HTTP/1.1 100 Continue":
		log.Println("\n(╯°□°)╯︵ ┻━┻\nFile does not exist on Yandex.Disk.\n")
	default:
		log.Printf("\nಠ_ಠ\nUnexpected response: %s\n\n", line)
	}
	// Delete temporary file.
	io.WriteString(conn, fmt.Sprintf("DELETE /%x.tmp HTTP/1.1\r\nHost: webdav.yandex.ru\r\nAccept: */*\r\nAuthorization: Basic %s\r\n\r\n", tmpname[:], auth))
}

func fileInfo(name string) (contentLength int64, md5sum, sha256sum []byte) {
	f, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	// Get file size.
	fi, err := f.Stat()
	if err != nil {
		log.Fatal(err)
	}
	contentLength = fi.Size()
	// Calculate hashes.
	mdw := md5.New()
	shaw := sha256.New()
	w := io.MultiWriter(mdw, shaw)
	if _, err := io.Copy(w, f); err != nil {
		log.Fatal(err)
	}
	md5sum = mdw.Sum(nil)
	sha256sum = shaw.Sum(nil)
	return
}
