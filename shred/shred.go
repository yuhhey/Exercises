package main

import (
	"crypto/rand"
	"flag"
	"io"
	"fmt"
	"log"
	"os"
)

func shredFile(path string, passes int, zero bool, remove bool, verbose bool) error {
	f, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		if os.IsPermission(err) {
			if err := os.Chmod(path, 0666); err != nil {
				return err
			}
			f, err = os.OpenFile(path, os.O_RDWR, 0666)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() {
    		return fmt.Errorf("%s is not a regular file", path)
	}
	/*if info.Mode().IsDir() {
			return fmt.Errorf("%s is a director.", path)
	}*/
	size := info.Size()

	for i := 0; i < passes; i++ {
		log.Printf("pass %d/%d", i+1, passes)
		if verbose {
			log.Printf("Pass %d/%d", i+1, passes)
		}
		if _, err := f.Seek(0, 0); err != nil {
			return err
		}
		log.Printf("Seek done")
		if i == passes-1 && zero {
			zeros := make([]byte, size)
			if _, err := f.Write(zeros); (err != nil) && (err != io.ErrShortWrite) {
				return err
			}
			log.Printf("Zero write  done")
		} else {
			buf := make([]byte, size)
			if _, err := rand.Read(buf); err != nil {
				return err
			}
			log.Printf("Read done")
			if _, err := f.Write(buf); (err != nil) && (err != io.ErrShortWrite) {
				log.Fatal(err)
				return err
			}
			log.Printf("Write done")
		}
		f.Sync()
		if verbose {
        		log.Printf("FINISHED pass %d/%d", i+1, passes)
    		}
	}

	if remove {
		return os.Remove(path)
	}
	return nil
}

func main() {
	n := flag.Int("n", 3, "number of overwrite passes")
	z := flag.Bool("z", false, "add final zero overwrite")
	u := flag.Bool("u", false, "remove file after overwriting")
	v := flag.Bool("v", false, "verbose output")
	flag.Parse()

	if flag.NArg() != 1 {
		log.Fatal("usage: shred [options] file")
	}

	if err := shredFile(flag.Arg(0), *n, *z, *u, *v); err != nil {
		log.Fatal(err)
	}
}

