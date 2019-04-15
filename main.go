package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func getallfiles(pathname string, s []string) ([]string, error) {
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		fmt.Println("read dir fail:", err)
		return s, err
	}
	for _, fi := range rd {
		if fi.IsDir() {
			fullDir := pathname + "/" + fi.Name()
			s, err = getallfiles(fullDir, s)
			if err != nil {
				fmt.Println("read dir fail:", err)
				return s, err
			}
		} else {
			fullName := pathname + "/" + fi.Name()
			s = append(s, fullName)
		}
	}
	return s, nil
}


var workResultLock sync.WaitGroup

func main() {
	log.Print("start")

	var s []string

	s, err := getallfiles("./", s)
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Print(s)

	total := 0
	for _, ss := range s {

		if strings.HasPrefix(filepath.Base(ss), "fuckbaiduyun") {
			continue
		}

		total++
	}

	en := 0
	for _, ss := range s{
		if strings.HasSuffix(ss, "fuckbaiduyun") {
			en++
		}
	}

	doen := false
	if en > total / 2 {
		doen = true
	}
	if en == total {
		doen = false
	}
	if en == 0 {
		doen = true
	}


	for _, ss := range s{

		if strings.HasPrefix(filepath.Base(ss), "fuckbaiduyun") {
			continue
		}
		if strings.HasSuffix(filepath.Base(ss), "fuckbaiduyun") {
			if !doen {
				workResultLock.Add(1)
				go defuck(ss)
			}
		} else {
			if doen {
				workResultLock.Add(1)
				go fuck(ss)
			}
		}
	}
	workResultLock.Wait()
}

func defuck(ss string) {
	log.Print("start back : ", ss)

	ifile, err := os.Open(ss)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer ifile.Close()


	// Open file for writing
	ofile, err := os.OpenFile(strings.TrimSuffix(ss, ".fuckbaiduyun"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer ofile.Close()


	bufferedReader := bufio.NewReader(ifile)

	bufferedWriter := bufio.NewWriter(ofile)

	for i := 0; i < len("fuckbaiduyun"); i++{
		bufferedReader.ReadByte()
	}

	byteSlice := make([]byte, 1024*1024)

	for {
		numBytesRead, err := bufferedReader.Read(byteSlice)

		if numBytesRead == 0 {
			break
		}

		if err != nil {
			log.Fatal(err)
			return
		}

		numBytesWrite, err := bufferedWriter.Write(byteSlice[:numBytesRead])
		if err != nil || numBytesRead != numBytesWrite {
			log.Fatal(err)
			log.Fatal(numBytesRead,numBytesWrite)
			return
		}

		bufferedWriter.Flush()
	}

	log.Print("end back : ", ss)

	err = os.Remove(ss)
	if err != nil {
		log.Fatal(err)
	}

	workResultLock.Done()
}

func fuck(ss string) {
	log.Print("start fuck : ", ss)

	ifile, err := os.Open(ss)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer ifile.Close()


	// Open file for writing
	ofile, err := os.OpenFile(ss+".fuckbaiduyun", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer ofile.Close()


	bufferedReader := bufio.NewReader(ifile)

	bufferedWriter := bufio.NewWriter(ofile)

	bufferedWriter.WriteString("fuckbaiduyun")
	bufferedWriter.Flush()

	byteSlice := make([]byte, 1024*1024)

	for {
		numBytesRead, err := bufferedReader.Read(byteSlice)

		if numBytesRead == 0 {
			break
		}

		if err != nil {
			log.Fatal(err)
			return
		}

		numBytesWrite, err := bufferedWriter.Write(byteSlice[:numBytesRead])
		if err != nil || numBytesRead != numBytesWrite {
			log.Fatal(err)
			log.Fatal(numBytesRead,numBytesWrite)
			return
		}

		bufferedWriter.Flush()
	}

	log.Print("end fuck : ", ss)

	err = os.Remove(ss)
	if err != nil {
		log.Fatal(err)
	}

	workResultLock.Done()
}
