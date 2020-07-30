package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	location, cwd  string
	count, maxSize int
	wg             *sync.WaitGroup
)

func init() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	wg = &sync.WaitGroup{}
	flag.StringVar(&location, "location", cwd, "directory of the file to delete")
	flag.IntVar(&count, "count", 1, "count of files to create")
	flag.IntVar(&maxSize, "max-size", -1, "maximal size of the files in the MB, -1 - infinite")
}

func main() {
	flag.Parse()

	end := make(<-chan struct{})

	files, err := createFiles(count, location)
	if err != nil {
		log.Fatal(err)
	}

	SetupCloseHandler(files)

	block := createBlock(1024 * 1024) //1MB

	for _, f := range files {
		f := f
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := writeToFile(f, maxSize, *block)
			if err != nil {
				log.Printf("err: %v", err)
			}
		}()
	}
	err = deleteFiles(files)
	if err != nil {
		log.Fatalf("failed to remove at least on of the files: %s", err)
	}
	log.Println("Still writing to deleted files.")
	wg.Wait()
	log.Println("Done writing.")
	<-end
}

func SetupCloseHandler(files []*os.File) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("\r- Ctrl+C pressed in Terminal")
		for _, f := range files {
			err := f.Close()
			if err != nil {
				log.Printf("error closing fd: %s", err)
			}
		}
		os.Exit(0)
	}()
}

func createFiles(count int, location string) ([]*os.File, error) {
	var files []*os.File
	for i := 0; i < count; i++ {
		f, err := ioutil.TempFile(location, "keep-")
		if err != nil {
			return nil, err
		}
		files = append(files, f)
		log.Println("Created file: ", f.Name())
	}
	return files, nil
}

func deleteFiles(files []*os.File) (err error) {
	for _, f := range files {
		err = os.Remove(f.Name())
		log.Printf("deleted file: %s err: %v", f.Name(), err)
	}
	return err
}

func createBlock(size int) *[]byte {
	var data byte = 0xFF
	block := new([]byte)
	for i := 0; i < size; i++ {
		*block = append(*block, data)
	}
	return block
}

func writeToFile(fd *os.File, maxSize int, block []byte) error {
	switch maxSize {
	case -1:
		for {
			_, err := fd.Write(block)
			if err != nil {
				return err
			}
		}
	default:
		for i := 0; i < maxSize; i++ {
			_, err := fd.Write(block)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
