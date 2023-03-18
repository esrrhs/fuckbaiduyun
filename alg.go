package main

import (
	"bufio"
	"crypto/md5"
	"crypto/rc4"
	"encoding/hex"
	"fmt"
	"github.com/esrrhs/gohome/loggo"
	"github.com/esrrhs/gohome/rbuffergo"
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
			fullName = filepath.FromSlash(fullName)
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

func dojob(jobtotal *int32, jobdone *int32, filetotal *int64, filedone *int64, input string, output string, doen bool,
	key string, split int) {

	input = input + "/"
	output = output + "/"
	input = filepath.FromSlash(input)
	output = filepath.FromSlash(output)
	input = filepath.Clean(input)
	output = filepath.Clean(output)

	var done map[string]int
	done = make(map[string]int)
	var workResultLock sync.WaitGroup

	loggo.Ini(loggo.Config{Level: loggo.LEVEL_DEBUG, Prefix: "fuck", MaxDay: 2})

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

	loggo.Info("all job file jobtotal %v", *jobtotal)

	for _, ss := range s {

		if strings.HasPrefix(filepath.Base(ss), "fuckbaiduyun") {
			continue
		}
		if strings.HasSuffix(filepath.Base(ss), "fuckbaiduyun") {
			if !doen {
				defuck(workResultLock, &num, key, split, ss, false,
					jobdone, jobtotal, done, input, output, filetotal, filedone)
			}
		} else {
			if doen {
				fuck(workResultLock, &num, key, split, ss, false,
					jobdone, jobtotal, done, input, output, filetotal, filedone)
			}
		}
	}

	workResultLock.Wait()
	delDone(output)
}

func defuck(workResultLock sync.WaitGroup, num *int32, key string, split int, ss string, flag bool,
	jobdone *int32, jobtotal *int32, done map[string]int, input string, output string,
	filetotal *int64, filedone *int64) {

	ss = filepath.FromSlash(ss)
	ss = filepath.Clean(ss)
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

	if outputss == ss {
		showErrorStr("filename is same " + ss)
	}

	inputfolderPath := filepath.Dir(ss)

	var son []string

	rd, err := ioutil.ReadDir(inputfolderPath)
	if err != nil {
		showError(err)
	}
	for _, fi := range rd {
		if !fi.IsDir() {
			m, _ := filepath.Match("*.fuckbaiduyun.*", fi.Name())
			if m {
				loggo.Info("back add split: %v %v", ss, fi.Name())
				name := inputfolderPath + "/" + fi.Name()
				son = append(son, filepath.FromSlash(name))
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
	fi, err := ifile.Stat()
	if err != nil {
		showError(err)
	}
	*filedone = 0
	*filetotal = fi.Size()

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
					fi, err := ifile.Stat()
					if err != nil {
						showError(err)
					}
					*filedone = 0
					*filetotal = fi.Size()
					loggo.Info("start back : %v", ss+"."+strconv.Itoa(post))
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
			*filedone += int64(numBytesRead)
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

	atomic.AddInt32(jobdone, 1)

	done[ss] = 1
	saveDone(output, ss)

	loggo.Info("end back : %v/%v %v", *jobdone, *jobtotal, ss)
}

func saveDone(output string, ss string) {

	name := filepath.FromSlash(output + "/fuckbaiduyunDONE")
	file, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		showError(err)
	}

	wd := bufio.NewWriter(file)
	wd.WriteString(filepath.FromSlash(ss) + "\n")
	wd.Flush()

	loggo.Info("save Done %v", filepath.FromSlash(ss))

	file.Close()
}

func delDone(output string) {
	name := filepath.FromSlash(output + "/fuckbaiduyunDONE")
	os.Remove(name)
}

func loadDone(output string, done map[string]int) {

	os.MkdirAll(output, os.ModePerm)

	name := filepath.FromSlash(output + "/fuckbaiduyunDONE")
	file, err := os.Open(name)
	if err != nil {
		file, err = os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
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

		name := filepath.FromSlash(line)
		name = strings.TrimSpace(name)
		done[name] = 1
		loggo.Info("load Done %v", name)
	}

	file.Close()
}

func fuckverify(key string, ss string, filetotal *int64, filedone *int64, md5str string) {

	ss = filepath.FromSlash(ss)
	ss = filepath.Clean(ss)
	loggo.Info("start fuckverify : %v", ss)

	inputfolderPath := filepath.Dir(ss)

	var son []string

	rd, err := ioutil.ReadDir(inputfolderPath)
	if err != nil {
		showError(err)
	}
	for _, fi := range rd {
		if !fi.IsDir() {
			m, _ := filepath.Match("*.fuckbaiduyun.*", fi.Name())
			if m {
				loggo.Info("back add split: %v %v", ss, fi.Name())
				name := inputfolderPath + "/" + fi.Name()
				son = append(son, filepath.FromSlash(name))
			}
		}
	}

	ifile, err := os.Open(ss)
	if err != nil {
		showError(err)
	}

	bufferedReader := bufio.NewReader(ifile)
	fi, err := ifile.Stat()
	if err != nil {
		showError(err)
	}
	*filedone = 0
	*filetotal = fi.Size()

	h := md5.New()

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
					fi, err := ifile.Stat()
					if err != nil {
						showError(err)
					}
					*filedone = 0
					*filetotal = fi.Size()
					loggo.Info("start back : %v", ss+"."+strconv.Itoa(post))
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
			*filedone += int64(numBytesRead)
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

			h.Write(d)
		}
	}

	ifile.Close()

	newmd5str := fmt.Sprintf("%x", h.Sum(nil))

	if newmd5str != md5str {
		showErrorStr("fuckverify fail " + ss)
	}

	loggo.Info("fuckverify ok: %v", ss)
}

func fuck(workResultLock sync.WaitGroup, num *int32, key string, split int, ss string, flag bool,
	jobdone *int32, jobtotal *int32, done map[string]int, input string, output string,
	filetotal *int64, filedone *int64) {

	ss = filepath.FromSlash(ss)
	ss = filepath.Clean(ss)
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

	if outputss == ss {
		showErrorStr("filename is same " + ss)
	}

	// Open file for writing
	ofile, err := os.OpenFile(outputss+".fuckbaiduyun",
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		showError(err)
	}

	bufferedReader := bufio.NewReader(ifile)
	fi, err := ifile.Stat()
	if err != nil {
		showError(err)
	}
	*filedone = 0
	*filetotal = fi.Size()

	bufferedWriter := bufio.NewWriter(ofile)

	byteSlice := make([]byte, 4*1024*1024)

	rb := rbuffergo.New(16*1024*1024, false)

	h := md5.New()

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
			*filedone += int64(numBytesRead)

			h.Write(byteSlice[:numBytesRead])
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
				loggo.Info("start fuck : %v", outputss+".fuckbaiduyun"+"."+strconv.Itoa(post))
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

	md5str := fmt.Sprintf("%x", h.Sum(nil))
	fuckverify(key, outputss+".fuckbaiduyun", filetotal, filedone, md5str)

	atomic.AddInt32(jobdone, 1)

	done[ss] = 1
	saveDone(output, ss)

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
	return dst
}

func decrypt(data []byte, passphrase string) []byte {
	c, err := rc4.NewCipher([]byte(createHash(passphrase)))
	if err != nil {
		showError(err)
	}
	dst := make([]byte, len(data))
	c.XORKeyStream(dst, data)
	return dst
}
