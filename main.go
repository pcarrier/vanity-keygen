package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/pem"
	"flag"
	"fmt"
	"github.com/mikesmitty/edkey"
	"golang.org/x/crypto/ssh"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"io"
	"log"
	"os"
	"regexp"
	"runtime"
	"sync/atomic"
	"time"
)

const syncStep = 1000

var (
	threads = flag.Int("threads", runtime.NumCPU(), "number of threads to run")
)

func incr(arr *[]byte) {
	for i := len(*arr) - 1; i >= 0; i-- {
		(*arr)[i]++
		if (*arr)[i] != 0 {
			break
		}
	}
}

type result struct {
	privRepr []byte
	pubRepr  []byte
}

func main() {
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Fprintf(flag.CommandLine.Output(), "%s [FLAGS] REGEX, with flags:\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}
	researched := regexp.MustCompile(flag.Arg(0))

	log.Printf("Looking for a public key matching %v", researched)
	start := time.Now()
	attempts := int64(0)
	resultCh := make(chan result, 1)

	for i := 0; i < *threads; i++ {
		go func() {
			localAttempts := 0
			seed := make([]byte, ed25519.SeedSize)
			publicKey := make(ed25519.PublicKey, ed25519.PublicKeySize)

			if _, err := io.ReadFull(rand.Reader, seed); err != nil {
				log.Fatalf("Could not read randomness: %v", err)
			}

			for {
				privateKey := ed25519.NewKeyFromSeed(seed)

				copy(publicKey, privateKey[32:])
				pubKey, err := ssh.NewPublicKey(publicKey)
				if err != nil {
					log.Fatalf("Could not derive public key: %v", err)
				}
				pubRepr := ssh.MarshalAuthorizedKey(pubKey)
				// Remove trailing newline
				pubRepr = pubRepr[:len(pubRepr)-1]

				if researched.Match(pubRepr) {
					resultCh <- result{
						privRepr: pem.EncodeToMemory(&pem.Block{
							Type:  "OPENSSH PRIVATE KEY",
							Bytes: edkey.MarshalED25519PrivateKey(privateKey),
						}),
						pubRepr: pubRepr,
					}
				}

				localAttempts++
				if localAttempts == syncStep {
					localAttempts = 0
					atomic.AddInt64(&attempts, syncStep)
				}

				incr(&seed)
			}
		}()
	}

	ticker := time.NewTicker(time.Second)
	printer := message.NewPrinter(language.English)

	for {
		select {
		case found := <-resultCh:
			log.Printf("Public key:\n%s", found.pubRepr)
			log.Printf("Private key:\n%s", found.privRepr)
			os.Exit(0)
		case <-ticker.C:
			rate := int(float64(attempts) / time.Since(start).Seconds())
			log.Printf("Generated %s keypairs (%s Hz)", printer.Sprint(atomic.LoadInt64(&attempts)), printer.Sprint(rate))
		}
	}
}
