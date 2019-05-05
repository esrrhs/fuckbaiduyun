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
	"sync/atomic"
)

func exists(path string) bool {
	_, err := os.Stat(path)
	//os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

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
var num int32
var jobtotal int32
var jobdone int32

func main() {
	log.Print("start")

	var s []string

	s, err := getallfiles("./", s)
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Print("get all file done ", len(s))

	for _, ss := range s {

		if strings.HasPrefix(filepath.Base(ss), "fuckbaiduyun") {
			continue
		}

		if strings.HasSuffix(filepath.Base(ss), "fuckbaiduyun") {
			if exists(strings.TrimSuffix(ss, ".fuckbaiduyun")) {
				err = os.Remove(ss)
				if err != nil {
					log.Fatal(err)
					return
				}
			}
		}
	}

	var s1 []string
	s, err = getallfiles("./", s1)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Print("fix all file done ", len(s))

	total := 0
	for _, ss := range s {

		if strings.HasPrefix(filepath.Base(ss), "fuckbaiduyun") {
			continue
		}

		total++
	}

	en := 0
	for _, ss := range s {
		if strings.HasSuffix(filepath.Base(ss), "fuckbaiduyun") {
			en++
		}
	}

	doen := false
	if en > total/2 {
		doen = true
	}
	if en == total {
		doen = false
	}
	if en == 0 {
		doen = true
	}

	num = 0

	jobdone = 0
	jobtotal = 0

	for _, ss := range s {

		if strings.HasPrefix(filepath.Base(ss), "fuckbaiduyun") {
			continue
		}
		if strings.HasSuffix(filepath.Base(ss), "fuckbaiduyun") {
			if !doen {
				jobtotal++
			}
		} else {
			if doen {
				jobtotal++
			}
		}
	}

	for _, ss := range s {

		if strings.HasPrefix(filepath.Base(ss), "fuckbaiduyun") {
			continue
		}
		if strings.HasSuffix(filepath.Base(ss), "fuckbaiduyun") {
			if !doen {
				if num < 100 {
					atomic.AddInt32(&num, 1)
					workResultLock.Add(1)
					go defuck(ss, true)
				} else {
					defuck(ss, false)
				}
			}
		} else {
			if doen {
				if num < 100 {
					atomic.AddInt32(&num, 1)
					workResultLock.Add(1)
					go fuck(ss, true)
				} else {
					fuck(ss, false)
				}
			}
		}
	}
	workResultLock.Wait()
}

func defuck(ss string, flag bool) {
	log.Print("start back : ", ss)

	if flag {
		defer workResultLock.Done()
		defer atomic.AddInt32(&num, -1)
	}

	ifile, err := os.Open(ss)
	if err != nil {
		log.Fatal(err)
		return
	}

	// Open file for writing
	ofile, err := os.OpenFile(strings.TrimSuffix(ss, ".fuckbaiduyun"), os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		log.Fatal(err)
		ifile.Close()
		return
	}

	bufferedReader := bufio.NewReader(ifile)

	bufferedWriter := bufio.NewWriter(ofile)

	for i := 0; i < len("fuckbaiduyun"); i++ {
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
			ifile.Close()
			ofile.Close()
			return
		}

		numBytesWrite, err := bufferedWriter.Write(byteSlice[:numBytesRead])
		if err != nil || numBytesRead != numBytesWrite {
			log.Fatal(err)
			log.Fatal(numBytesRead, numBytesWrite)
			ifile.Close()
			ofile.Close()
			return
		}

		bufferedWriter.Flush()
	}

	ifile.Close()
	ofile.Close()

	err = os.Remove(ss)
	if err != nil {
		log.Fatal(err)
		return
	}

	atomic.AddInt32(&jobdone, 1)

	log.Print("end back : ", jobdone, "/", jobtotal, " ", ss)
}

func fuck(ss string, flag bool) {
	log.Print("start fuck : ", ss)

	if flag {
		defer workResultLock.Done()
		defer atomic.AddInt32(&num, -1)
	}

	ifile, err := os.Open(ss)
	if err != nil {
		log.Fatal(err)
		return
	}

	// Open file for writing
	ofile, err := os.OpenFile(ss+".fuckbaiduyun", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal(err)
		ifile.Close()
		return
	}

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
			ifile.Close()
			ofile.Close()
			return
		}

		numBytesWrite, err := bufferedWriter.Write(byteSlice[:numBytesRead])
		if err != nil || numBytesRead != numBytesWrite {
			log.Fatal(err)
			log.Fatal(numBytesRead, numBytesWrite)
			ifile.Close()
			ofile.Close()
			return
		}

		bufferedWriter.Flush()
	}

	ifile.Close()
	ofile.Close()

	err = os.Remove(ss)
	if err != nil {
		log.Fatal(err)
		return
	}

	atomic.AddInt32(&jobdone, 1)

	log.Print("end fuck : ", jobdone, "/", jobtotal, " ", ss)
}
