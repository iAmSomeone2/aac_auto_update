// Package update contains functionality for grabbing the most recent version of
// the crowdfunding patrons file, comparing the new copy against the old one, and
// producing a slice of Patrons.
package update

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

const (
	AppDir       string = "adopt-a-cell"
	BaseFileName string = "patrons_raw.txt"
	OldFileName  string = "patrons_raw.old.txt"
	searchStr    string = "var data"
)

func GetCacheDir() string {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		log.Printf("WARN: Cache directory not found! Using working directory instead.")
		cacheDir = "./"
	}

	return cacheDir
}

// CheckForUpdate grabs the latest version of the patrons file from the crowdfunding
// page and determines if anything in the file has changed. If a change is detected,
// then a string pointing to the resulting file is returned. Otherwise, "" is returned.
// Additionally, if the file cannot be downloaded or the downloaded file is
// identical to the original, "" is returned.
func CheckForUpdate(url string) string {
	// Create the file for the contents to be read into.
	cacheDir := GetCacheDir()
	cacheDir = path.Join(cacheDir, AppDir)

	err := os.MkdirAll(cacheDir, os.ModePerm)
	if err != nil {
		log.Panic(err)
	}

	fullBase := path.Join(cacheDir, BaseFileName)
	fullOldBase := path.Join(cacheDir, OldFileName)

	err = os.Rename(fullBase, fullOldBase)

	out, err := os.Create(fullBase)
	if err != nil {
		log.Println(err)
		return ""
	}
	defer out.Close()

	// Grab the file from online
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return ""
	}
	defer resp.Body.Close()

	// Write the contents of the GET method to the file
	n, err := io.Copy(out, resp.Body)
	if err != nil {
		log.Println(err)
		return ""
	}

	log.Printf("\nBytes copied to %s: %d\n", fullBase, n)

	err = cleanFile(fullBase)
	if err != nil {
		log.Println(err)
	}

	/*
		Compare the two files. If there is no difference, return an empty
		string to indicate that nothing further needs to be done.
	*/
	different := compareFiles(fullBase, fullOldBase)

	if different {
		return out.Name()
	}
	return ""
}

// md5Hash computes and returns the MD5 hash of the file the filePath string
// specifies.
func md5Hash(filePath string) (string, error) {
	// Init the return string in case things break
	var md5Str string

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return md5Str, err
	}
	defer file.Close()

	hash := md5.New()

	// Copy contents of file into the hash interface
	if _, err := io.Copy(hash, file); err != nil {
		return md5Str, err
	}

	// Get the 16 byte hash
	byteHash := hash.Sum(nil)[:16]

	// Convert to string
	md5Str = hex.EncodeToString(byteHash)
	return md5Str, nil
}

// compareFiles returns true if the files have different MD5 hashes.
func compareFiles(filePath0, filePath1 string) bool {
	// Start by getting the hashes
	hash0, err := md5Hash(filePath0)
	if err != nil {
		log.Println(err)
	}
	hash1, err := md5Hash(filePath1)
	if err != nil {
		log.Println(err)
	}

	log.Printf("Hash0: %s\n", hash0)
	log.Printf("Hash1: %s\n", hash1)

	return hash0 != hash1
}

// cleanFile gets rid of everything in the given file except for the line:
// 'var data = ...'.
// This will make parsing easier and allow the MD5 check to work correctly.
func cleanFile(filePath string) error {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	content := string(file)

	// Read through the file and only grab the "var data=..." line
	scanner := bufio.NewScanner(strings.NewReader(content))
	var data string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, searchStr) {
			data = line
		}
	}

	// Overwrite the file with the new content
	err = ioutil.WriteFile(filePath, []byte(data), 0644)
	if err != nil {
		return err
	}

	return nil
}
