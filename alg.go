package main

import (
	"bufio"
	"crypto/md5"
	"crypto/rc4" //nolint:staticcheck // RC4 需与已有加密文件保持兼容
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"
)

const (
	chunkSize = 1024 * 1024
	encSuffix = ".fuckbaiduyun"
	doneName  = "fuckbaiduyunDONE"
	logName   = "fuck.log"
)

type Progress struct {
	jobsTotal atomic.Int32
	jobsDone  atomic.Int32
	fileTotal atomic.Int64
	fileDone  atomic.Int64
}

func (p *Progress) Jobs() (done, total int32) {
	return p.jobsDone.Load(), p.jobsTotal.Load()
}

func (p *Progress) File() (done, total int64) {
	return p.fileDone.Load(), p.fileTotal.Load()
}

func initLog() {
	path := logName
	if exe, err := os.Executable(); err == nil {
		path = filepath.Join(filepath.Dir(exe), logName)
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.SetFlags(log.LstdFlags)
		return
	}
	log.SetOutput(io.MultiWriter(os.Stdout, f))
	log.SetFlags(log.LstdFlags)
}

func getallfiles(root string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func pathsOverlap(a, b string) bool {
	a, err1 := filepath.Abs(a)
	b, err2 := filepath.Abs(b)
	if err1 != nil || err2 != nil {
		return strings.Contains(a, b) || strings.Contains(b, a)
	}
	a = filepath.Clean(a)
	b = filepath.Clean(b)
	if strings.EqualFold(a, b) {
		return true
	}
	sep := string(os.PathSeparator)
	al := strings.ToLower(a + sep)
	bl := strings.ToLower(b + sep)
	return strings.HasPrefix(al, bl) || strings.HasPrefix(bl, al)
}

func skipName(name string) bool {
	base := filepath.Base(name)
	return strings.HasPrefix(base, "fuckbaiduyun")
}

func isEncryptedMain(name string) bool {
	base := filepath.Base(name)
	return strings.HasSuffix(base, encSuffix)
}

func isEncryptedSplit(name string) bool {
	matched, _ := filepath.Match("*"+encSuffix+".*", filepath.Base(name))
	return matched
}

func dojob(p *Progress, input, output string, encrypting bool, key string, split int) error {
	input = filepath.Clean(input)
	output = filepath.Clean(output)

	log.Printf("start input=%s output=%s encrypt=%v split=%d", input, output, encrypting, split)

	if pathsOverlap(input, output) {
		return errors.New("输入输出目录有重叠")
	}

	log.Printf("get all file begin %s", input)
	files, err := getallfiles(input)
	if err != nil {
		return err
	}
	log.Printf("get all file done %d", len(files))

	done := make(map[string]struct{})
	if err := loadDone(output, done); err != nil {
		return err
	}
	log.Printf("loadDone %d", len(done))

	var jobtotal int32
	for _, ss := range files {
		if skipName(ss) {
			continue
		}
		if isEncryptedMain(ss) {
			if !encrypting {
				jobtotal++
			}
		} else if !isEncryptedSplit(ss) && encrypting {
			jobtotal++
		}
	}
	p.jobsTotal.Store(jobtotal)
	log.Printf("all job file jobtotal %d", jobtotal)

	for _, ss := range files {
		if skipName(ss) {
			continue
		}
		if isEncryptedMain(ss) {
			if !encrypting {
				if err := defuck(p, key, ss, done, input, output); err != nil {
					return err
				}
			}
			continue
		}
		if isEncryptedSplit(ss) {
			continue
		}
		if encrypting {
			if err := fuck(p, key, split, ss, done, input, output); err != nil {
				return err
			}
		}
	}

	return delDone(output)
}

func defuck(p *Progress, key, ss string, done map[string]struct{}, input, output string) error {
	ss = filepath.Clean(ss)
	log.Printf("start back : %s", ss)

	if _, ok := done[ss]; ok {
		n := p.jobsDone.Add(1)
		log.Printf("end back skip done : %d/%d %s", n, p.jobsTotal.Load(), ss)
		return nil
	}

	rel, err := filepath.Rel(input, ss)
	if err != nil {
		return err
	}
	outputss := filepath.Join(output, strings.TrimSuffix(rel, encSuffix))
	if outputss == ss {
		return errors.New("filename is same " + ss)
	}
	if err := os.MkdirAll(filepath.Dir(outputss), os.ModePerm); err != nil {
		return err
	}

	ofile, err := os.OpenFile(outputss, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		return err
	}
	defer ofile.Close()
	bufferedWriter := bufio.NewWriter(ofile)

	if err := processParts(ss, key, p, func(plain []byte) error {
		n, err := bufferedWriter.Write(plain)
		if err != nil {
			return err
		}
		if n != len(plain) {
			return fmt.Errorf("diff size %d %d", len(plain), n)
		}
		return bufferedWriter.Flush()
	}); err != nil {
		return err
	}

	if err := bufferedWriter.Flush(); err != nil {
		return err
	}
	if err := ofile.Close(); err != nil {
		return err
	}

	n := p.jobsDone.Add(1)
	done[ss] = struct{}{}
	if err := saveDone(output, ss); err != nil {
		return err
	}
	log.Printf("end back : %d/%d %s", n, p.jobsTotal.Load(), ss)
	return nil
}

func saveDone(output, ss string) error {
	name := filepath.Join(output, doneName)
	file, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.WriteString(filepath.Clean(ss) + "\n"); err != nil {
		return err
	}
	log.Printf("save Done %s", ss)
	return nil
}

func delDone(output string) error {
	err := os.Remove(filepath.Join(output, doneName))
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func loadDone(output string, done map[string]struct{}) error {
	if err := os.MkdirAll(output, os.ModePerm); err != nil {
		return err
	}

	name := filepath.Join(output, doneName)
	file, err := os.Open(name)
	if err != nil {
		if os.IsNotExist(err) {
			f, cerr := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
			if cerr != nil {
				return cerr
			}
			return f.Close()
		}
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		item := filepath.Clean(strings.TrimSpace(scanner.Text()))
		if item == "" {
			continue
		}
		done[item] = struct{}{}
		log.Printf("load Done %s", item)
	}
	return scanner.Err()
}

func processParts(mainPath, key string, p *Progress, handle func([]byte) error) error {
	path := mainPath
	part := 0
	for {
		ifile, err := os.Open(path)
		if err != nil {
			if part == 0 || !os.IsNotExist(err) {
				return err
			}
			return nil
		}

		fi, err := ifile.Stat()
		if err != nil {
			ifile.Close()
			return err
		}
		p.fileDone.Store(0)
		p.fileTotal.Store(fi.Size())
		log.Printf("start part : %s", path)

		buf := make([]byte, chunkSize)
		for {
			n, rerr := io.ReadFull(ifile, buf)
			if n > 0 {
				p.fileDone.Add(int64(n))
				plain := decrypt(buf[:n], key)
				if err := handle(plain); err != nil {
					ifile.Close()
					return err
				}
			}
			if rerr == nil {
				continue
			}
			if errors.Is(rerr, io.EOF) || errors.Is(rerr, io.ErrUnexpectedEOF) {
				break
			}
			ifile.Close()
			return rerr
		}
		ifile.Close()

		path = mainPath + "." + strconv.Itoa(part)
		part++
	}
}

func fuckverify(key, ss string, p *Progress, md5str string) error {
	ss = filepath.Clean(ss)
	log.Printf("start fuckverify : %s", ss)

	h := md5.New()
	if err := processParts(ss, key, p, func(plain []byte) error {
		_, err := h.Write(plain)
		return err
	}); err != nil {
		return err
	}

	newmd5str := fmt.Sprintf("%x", h.Sum(nil))
	if newmd5str != md5str {
		return errors.New("fuckverify fail " + ss)
	}
	log.Printf("fuckverify ok: %s", ss)
	return nil
}

func fuck(p *Progress, key string, split int, ss string, done map[string]struct{}, input, output string) error {
	ss = filepath.Clean(ss)
	log.Printf("start fuck : %s", ss)

	if _, ok := done[ss]; ok {
		n := p.jobsDone.Add(1)
		log.Printf("end fuck skip done : %d/%d %s", n, p.jobsTotal.Load(), ss)
		return nil
	}

	ifile, err := os.Open(ss)
	if err != nil {
		return err
	}
	defer ifile.Close()

	rel, err := filepath.Rel(input, ss)
	if err != nil {
		return err
	}
	outputss := filepath.Join(output, rel)
	if outputss == ss {
		return errors.New("filename is same " + ss)
	}
	if err := os.MkdirAll(filepath.Dir(outputss), os.ModePerm); err != nil {
		return err
	}

	outPath := outputss + encSuffix
	ofile, err := os.OpenFile(outPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	bufferedWriter := bufio.NewWriter(ofile)

	closeCurrent := func() error {
		if err := bufferedWriter.Flush(); err != nil {
			_ = ofile.Close()
			return err
		}
		return ofile.Close()
	}

	fi, err := ifile.Stat()
	if err != nil {
		_ = closeCurrent()
		return err
	}
	p.fileDone.Store(0)
	p.fileTotal.Store(fi.Size())

	h := md5.New()
	buf := make([]byte, chunkSize)
	var cur, post int

	for {
		n, rerr := io.ReadFull(ifile, buf)
		if n > 0 {
			p.fileDone.Add(int64(n))
			if _, err := h.Write(buf[:n]); err != nil {
				_ = closeCurrent()
				return err
			}
			d := encrypt(buf[:n], key)
			cur += n
			if split > 0 && cur > split {
				if err := closeCurrent(); err != nil {
					return err
				}
				next := outputss + encSuffix + "." + strconv.Itoa(post)
				ofile, err = os.OpenFile(next, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
				if err != nil {
					return err
				}
				bufferedWriter = bufio.NewWriter(ofile)
				post++
				cur -= split
				log.Printf("start fuck : %s", next)
			}
			wn, err := bufferedWriter.Write(d)
			if err != nil {
				_ = closeCurrent()
				return err
			}
			if wn != len(d) {
				_ = closeCurrent()
				return fmt.Errorf("diff size %d %d", len(d), wn)
			}
			if err := bufferedWriter.Flush(); err != nil {
				_ = closeCurrent()
				return err
			}
		}
		if rerr == nil {
			continue
		}
		if errors.Is(rerr, io.EOF) || errors.Is(rerr, io.ErrUnexpectedEOF) {
			break
		}
		_ = closeCurrent()
		return rerr
	}

	if err := closeCurrent(); err != nil {
		return err
	}

	md5str := fmt.Sprintf("%x", h.Sum(nil))
	if err := fuckverify(key, outPath, p, md5str); err != nil {
		return err
	}

	n := p.jobsDone.Add(1)
	done[ss] = struct{}{}
	if err := saveDone(output, ss); err != nil {
		return err
	}
	log.Printf("end fuck : %d/%d %s", n, p.jobsTotal.Load(), ss)
	return nil
}

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func crypt(data []byte, passphrase string) []byte {
	c, err := rc4.NewCipher([]byte(createHash(passphrase)))
	if err != nil {
		panic(err)
	}
	dst := make([]byte, len(data))
	c.XORKeyStream(dst, data)
	return dst
}

func encrypt(data []byte, passphrase string) []byte {
	return crypt(data, passphrase)
}

func decrypt(data []byte, passphrase string) []byte {
	return crypt(data, passphrase)
}
