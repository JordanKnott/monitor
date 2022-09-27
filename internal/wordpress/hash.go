package wordpress

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"os"
	"path/filepath"
)

func makeHash(f string) []byte {
	file, err := os.Open(f)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	buf := make([]byte, 30*1024)
	sha256 := sha256.New()
	for {
		n, err := file.Read(buf)
		if n > 0 {
			_, err := sha256.Write(buf[:n])
			if err != nil {
				log.Fatal(err)
			}
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Printf("Read %d bytes: %v", n, err)
			break
		}
	}

	return sha256.Sum(nil)
}

func MakeHashTree(filepaths []string) map[string]string {
	fileHashes := make(map[string]string)
	for _, path := range filepaths {

		err := filepath.Walk(path,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() {
					hash := makeHash(path)
					fileHashes[path] = hex.EncodeToString(hash)
				}
				return nil
			})
		if err != nil {
			log.Println(err)
		}
	}
	return fileHashes
}
