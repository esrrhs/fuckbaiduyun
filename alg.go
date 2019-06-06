package main

import (
	"bufio"
	"crypto/md5"
	"crypto/rc4"
	"encoding/hex"
	"fmt"
	"github.com/esrrhs/go-engine/src/loggo"
	"github.com/esrrhs/go-engine/src/rbuffergo"
	"io"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strconv"
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

func showError(err error) {
	loggo.Error("%v", err)
	panic(err.Error())
}

func showErrorStr(err string) {
	loggo.Error("%v", err)
	panic(err)
}

func dojob(jobtotal *int32, jobdone *int32, input string, output string, doen bool,
	key string, split int) {

	var done map[string]int
	done = make(map[string]int)
	var workResultLock sync.WaitGroup

	loggo.Ini(loggo.Config{loggo.LEVEL_DEBUG, "fuck", 2})

	loggo.Info("start")

	if strings.Contains(output, input) || strings.Contains(input, output) {
		showErrorStr("输入输出目录有重叠")
		return
	}

	var s []string

	loggo.Info("get all file begin %v", input)

	s, err := getallfiles(input+"/", s)
	if err != nil {
		showError(err)
	}

	loggo.Info("get all file done %v", len(s))

	loadDone(output, done)
	loggo.Info("loadDone %v", len(done))

	total := 0
	for _, ss := range s {

		if strings.HasPrefix(filepath.Base(ss), "fuckbaiduyun") {
			continue
		}

		total++
	}

	loggo.Info("all file total %v", total)

	var num int32

	num = 0

	for _, ss := range s {

		if strings.HasPrefix(filepath.Base(ss), "fuckbaiduyun") {
			continue
		}
		if strings.HasSuffix(filepath.Base(ss), "fuckbaiduyun") {
			if !doen {
				*jobtotal++
			}
		} else {
			if doen {
				*jobtotal++
			}
		}
	}

	loggo.Info("all job file jobtotal %v", jobtotal)

	for _, ss := range s {

		if strings.HasPrefix(filepath.Base(ss), "fuckbaiduyun") {
			continue
		}
		if strings.HasSuffix(filepath.Base(ss), "fuckbaiduyun") {
			if !doen {
				defuck(workResultLock, &num, key, split, ss, false,
					jobdone, jobtotal, done, input, output)
			}
		} else {
			if doen {
				fuck(workResultLock, &num, key, split, ss, false,
					jobdone, jobtotal, done, input, output)
			}
		}
	}

	workResultLock.Wait()
	delDone(output)
}

func defuck(workResultLock sync.WaitGroup, num *int32, key string, split int, ss string, flag bool,
	jobdone *int32, jobtotal *int32, done map[string]int, input string, output string) {

	loggo.Info("start back : %v", ss)

	if flag {
		defer workResultLock.Done()
		defer atomic.AddInt32(num, -1)
	}

	if done[ss] == 1 {
		atomic.AddInt32(jobdone, 1)
		loggo.Info("end back skip done : %v/%v %v", *jobdone, *jobtotal, ss)
		return
	}

	outputss := strings.Replace(strings.TrimSuffix(ss, ".fuckbaiduyun"), input, output, -1)
	folderPath := filepath.Dir(outputss)
	os.MkdirAll(folderPath, os.ModePerm)

	var son []string

	rd, err := ioutil.ReadDir(folderPath)
	if err != nil {
		showError(err)
	}
	for _, fi := range rd {
		if !fi.IsDir() {
			m, _ := filepath.Match("*.fuckbaiduyun.*", fi.Name())
			if m {
				loggo.Info("back add split: %v %v", ss, fi.Name())
				son = append(son, folderPath+"/"+fi.Name())
			}
		}
	}

	ifile, err := os.Open(ss)
	if err != nil {
		showError(err)
	}

	// Open file for writing
	ofile, err := os.OpenFile(outputss, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		showError(err)
	}

	bufferedReader := bufio.NewReader(ifile)

	bufferedWriter := bufio.NewWriter(ofile)

	byteSlice := make([]byte, 4*1024*1024)

	rb := rbuffergo.New(16*1024*1024, false)

	var post int

	fileend := false
	for !fileend || !rb.Empty() {
		for !fileend && rb.Size()+len(byteSlice) < rb.Capacity() {

			numBytesRead, err := bufferedReader.Read(byteSlice)

			if numBytesRead == 0 {
				find := false

				for _, sf := range son {
					if sf == ss+"."+strconv.Itoa(post) {
						find = true
					}
				}

				if find {
					ifile.Close()
					ifile, err = os.Open(ss + "." + strconv.Itoa(post))
					if err != nil {
						showError(err)
					}
					post++
					bufferedReader = bufio.NewReader(ifile)
					continue
				} else {
					fileend = true
					break
				}
			}

			if err != nil {
				showError(err)
			}

			rb.Write(byteSlice[:numBytesRead])
		}

		for !rb.Empty() {
			if rb.Size() < 1024*1024 {
				if !fileend {
					break
				}
			}

			numBytesRead := int(math.Min(float64(rb.Size()), 1024*1024))

			if !rb.Read(byteSlice[0:numBytesRead]) {
				showErrorStr("rbuffergo read fail " + ss)
			}

			d := decrypt(byteSlice[:numBytesRead], key)
			numBytesRead = len(d)

			numBytesWrite, err := bufferedWriter.Write(d)
			if err != nil {
				showError(err)
			}
			if numBytesRead != numBytesWrite {
				showErrorStr("diff size " + strconv.Itoa(numBytesRead) + " " + strconv.Itoa(numBytesWrite))
			}

			bufferedWriter.Flush()
		}
	}

	ifile.Close()
	ofile.Close()

	err = os.Remove(ss)
	if err != nil {
		showError(err)
	}

	atomic.AddInt32(jobdone, 1)

	done[ss] = 1
	saveDone(output)

	loggo.Info("end back : %v/%v %v", *jobdone, *jobtotal, ss)
}

func saveDone(output string) {

	file, err := os.Open(output + "/fuckbaiduyunDONE")
	if err != nil {
		file, err = os.OpenFile(output+"/fuckbaiduyunDONE", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			showError(err)
		}
	}

	wd := bufio.NewWriter(file)
	wd.WriteString(output + "\n")

	file.Close()
}

func delDone(output string) {
	os.Remove(output + "/fuckbaiduyunDONE")
}

func loadDone(output string, done map[string]int) {

	os.MkdirAll(output, os.ModePerm)

	file, err := os.Open(output + "/fuckbaiduyunDONE")
	if err != nil {
		file, err = os.OpenFile(output+"/fuckbaiduyunDONE", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			showError(err)
		}
	}

	rd := bufio.NewReader(file)
	for {
		line, err := rd.ReadString('\n') //以'\n'为结束符读入一行

		if err != nil || io.EOF == err {
			break
		}

		done[line] = 1
	}

	file.Close()
}

func fuck(workResultLock sync.WaitGroup, num *int32, key string, split int, ss string, flag bool,
	jobdone *int32, jobtotal *int32, done map[string]int, input string, output string) {

	loggo.Info("start fuck : %v", ss)

	if flag {
		defer workResultLock.Done()
		defer atomic.AddInt32(num, -1)
	}

	if done[ss] == 1 {
		atomic.AddInt32(jobdone, 1)
		loggo.Info("end fuck skip done : %v/%v %v", *jobdone, *jobtotal, ss)
		return
	}

	m, _ := filepath.Match("*.fuckbaiduyun.*", filepath.Base(ss))
	if m {
		atomic.AddInt32(jobdone, 1)
		loggo.Info("end fuck skip split: %v/%v %v", *jobdone, *jobtotal, ss)
		return
	}

	ifile, err := os.Open(ss)
	if err != nil {
		showError(err)
	}

	outputss := strings.Replace(ss, input, output, -1)
	folderPath := filepath.Dir(outputss)
	os.MkdirAll(folderPath, os.ModePerm)

	// Open file for writing
	ofile, err := os.OpenFile(outputss+".fuckbaiduyun",
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		showError(err)
	}

	bufferedReader := bufio.NewReader(ifile)

	bufferedWriter := bufio.NewWriter(ofile)

	byteSlice := make([]byte, 4*1024*1024)

	rb := rbuffergo.New(16*1024*1024, false)

	var cur int
	var post int

	fileend := false
	for !fileend || !rb.Empty() {
		for !fileend && rb.Size()+len(byteSlice) < rb.Capacity() {

			numBytesRead, err := bufferedReader.Read(byteSlice)

			if numBytesRead == 0 {
				fileend = true
				break
			}

			if err != nil {
				showError(err)
			}

			rb.Write(byteSlice[:numBytesRead])
		}

		for !rb.Empty() {
			if rb.Size() < 1024*1024 {
				if !fileend {
					break
				}
			}

			numBytesRead := int(math.Min(float64(rb.Size()), 1024*1024))

			if !rb.Read(byteSlice[0:numBytesRead]) {
				showErrorStr("rbuffergo read fail " + ss)
			}

			d := encrypt(byteSlice[:numBytesRead], key)

			cur += numBytesRead

			if cur > split {
				ofile.Close()
				bufferedWriter.Flush()
				ofile, err = os.OpenFile(outputss+".fuckbaiduyun"+"."+strconv.Itoa(post),
					os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
				if err != nil {
					showError(err)
				}
				bufferedWriter = bufio.NewWriter(ofile)
				post++
				cur -= split
			}

			numBytesWrite, err := bufferedWriter.Write(d)
			if err != nil {
				showError(err)
			}
			if numBytesRead != numBytesWrite {
				showErrorStr("diff size " + strconv.Itoa(numBytesRead) + " " + strconv.Itoa(numBytesWrite))
			}

			bufferedWriter.Flush()
		}
	}

	ifile.Close()
	ofile.Close()

	atomic.AddInt32(jobdone, 1)

	done[ss] = 1
	saveDone(output)

	loggo.Info("end fuck : %v/%v %v", *jobdone, *jobtotal, ss)
}

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func encrypt(data []byte, passphrase string) []byte {
	c, err := rc4.NewCipher([]byte(createHash(passphrase)))
	if err != nil {
		showError(err)
	}
	dst := make([]byte, len(data))
	c.XORKeyStream(dst, data)
	return data
}

func decrypt(data []byte, passphrase string) []byte {
	c, err := rc4.NewCipher([]byte(createHash(passphrase)))
	if err != nil {
		showError(err)
	}
	dst := make([]byte, len(data))
	c.XORKeyStream(dst, data)
	return data
}
