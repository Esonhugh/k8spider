//go:build ignore

// Written by Rob Fuller / mubix
// Source: https://github.com/vmware/go-nfs-client
// Using: https://github.com/go-nfs/nfsv3

package nfs

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/go-nfs/nfsv3/nfs"
	"github.com/go-nfs/nfsv3/nfs/rpc"
	log "github.com/sirupsen/logrus"
)

func ls(v *nfs.Target, path string) ([]*nfs.EntryPlus, error) {
	dirs, err := v.ReadDirPlus(path)
	if err != nil {
		return nil, fmt.Errorf("readdir error: %s", err.Error())
	}
	log.Debugf("Filename\tUID\tGID\tMode\tSize")
	for _, d := range dirs {
		log.Debugf("%s\t%s\t%s\t%s\t%d", d.FileName, d.Attr.Attr.UID, d.Attr.Attr.GID, d.Attr.Attr.Mode, d.Attr.Attr.Size())
	}
	return dirs, nil
}

func uploadFile(v *nfs.Target, path string, dest string) error {
	f, err := os.Open(path)
	log.Tracef("Opening %s", path)
	if err != nil {
		log.Errorf("Could not open local file: %s", err.Error())
		return err
	}
	fileinfo, err := f.Stat()
	if err != nil {
		log.Errorf("read fail: %s", err.Error())
		return err
	}
	filesize := fileinfo.Size()
	wr, err := v.OpenFile(dest, 0777)
	if err != nil {
		log.Errorf("write fail: %s", err.Error())
		return err
	}

	// calculate the sha
	h := sha256.New()
	t := io.TeeReader(f, h)

	// Copy filesize
	n, err := io.CopyN(wr, t, int64(filesize))
	if err != nil {
		log.Errorf("error copying: n=%d, %s", n, err.Error())
		return err
	}
	expectedSum := h.Sum(nil)

	if err = wr.Close(); err != nil {
		log.Errorf("error committing: %s", err.Error())
		return err
	}

	//
	// get the file we wrote and calc the sum
	rdr, err := v.Open(dest)
	if err != nil {
		log.Errorf("read error: %v", err)
		return err
	}

	h = sha256.New()
	t = io.TeeReader(rdr, h)

	_, err = ioutil.ReadAll(t)
	if err != nil {
		log.Errorf("readall error: %v", err)
		return err
	}
	actualSum := h.Sum(nil)

	if bytes.Compare(actualSum, expectedSum) != 0 {
		log.Errorf("sums didn't match. actual=%x expected=%s", actualSum, expectedSum) //  Got=0%x expected=0%x", string(buf), testdata)
		return errors.New("sums didn't match")
	}

	log.Debugf("Sums match %x %x", actualSum, expectedSum)
	return nil
}

func catFile(v *nfs.Target, path string) error {
	wr, err := v.Open(path)
	if err != nil {
		log.Errorf("read fail: %s", err.Error())
		return err
	}
	defer wr.Close()

	pathsplit := strings.Split(path, "/")
	filename := pathsplit[len(pathsplit)-1]

	fileInfo, _, err := wr.Lookup(path)
	if err != nil {
		log.Errorf("look up path %v fail: %s", path, err.Error())
		return err
	}
	filesize := fileInfo.Size()
	log.Debugf("File size: %d", filesize)

	if err != nil {
		log.Errorf("read fail: %s", err.Error())
		return err
	}
	log.Infof("Name: %s", filename)
	f := os.Stdout

	// calculate the sha
	h := sha256.New()
	t := io.TeeReader(wr, h)

	// Copy filesize
	n, err := io.CopyN(f, t, int64(filesize))
	if err != nil {
		log.Errorf("error copying: n=%d, %s", n, err.Error())
		return err
	}
	expectedSum := h.Sum(nil)

	if err = wr.Close(); err != nil {
		log.Errorf("error committing: %s", err.Error())
		return err
	}
	f.Close()

	//
	// get the file we wrote and calc the sum
	rdr, err := os.Open(filename)
	if err != nil {
		log.Errorf("read error: %v", err)
		return err
	}

	h = sha256.New()
	t = io.TeeReader(rdr, h)

	_, err = ioutil.ReadAll(t)
	if err != nil {
		log.Errorf("readall error: %v", err)
		return err
	}
	actualSum := h.Sum(nil)

	if bytes.Compare(actualSum, expectedSum) != 0 {
		log.Fatalf("sums didn't match. actual=%x expected=%s", actualSum, expectedSum) //  Got=0%x expected=0%x", string(buf), testdata)
	}

	log.Debugf("Sums match %x %x", actualSum, expectedSum)
	return err
}

func downloadFile(v *nfs.Target, path string) error {
	wr, err := v.Open(path)
	if err != nil {
		log.Errorf("read fail: %s", err.Error())
		return err
	}
	pathsplit := strings.Split(path, "/")
	filename := pathsplit[len(pathsplit)-1]

	fileInfo, _, err := wr.Lookup(path)
	if err != nil {
		log.Errorf("look up path %v fail: %s", path, err.Error())
		return err
	}
	filesize := fileInfo.Size()
	log.Debugf("File size: %d", filesize)

	log.Debugf("Name: %s", filename)
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0777)

	if err != nil {
		log.Errorf("read fail: %s", err.Error())
		return err
	}

	// calculate the sha
	h := sha256.New()
	t := io.TeeReader(wr, h)

	// Copy filesize
	n, err := io.CopyN(f, t, int64(filesize))
	if err != nil {
		log.Errorf("error copying: n=%d, %s", n, err.Error())
		return err
	}
	expectedSum := h.Sum(nil)

	if err = wr.Close(); err != nil {
		log.Errorf("error committing: %s", err.Error())
		return err
	}
	f.Close()

	//
	// get the file we wrote and calc the sum
	rdr, err := os.Open(filename)
	if err != nil {
		log.Errorf("read error: %v", err)
		return err
	}

	h = sha256.New()
	t = io.TeeReader(rdr, h)

	_, err = ioutil.ReadAll(t)
	if err != nil {
		log.Errorf("readall error: %v", err)
		return err
	}
	actualSum := h.Sum(nil)

	if bytes.Compare(actualSum, expectedSum) != 0 {
		log.Errorf("sums didn't match. actual=%x expected=%s", actualSum, expectedSum) //  Got=0%x expected=0%x", string(buf), testdata)
		return errors.New("sums didn't match")
	}

	log.Debugf("Sums match %x %x", actualSum, expectedSum)
	wr.Close()
	f.Close()
	return nil
}

func removeFile(v *nfs.Target, path string) {
	err := v.Remove(path)
	if err != nil {
		log.Fatalf("rm of %s err: %s", path, err.Error())
	}
}

func makeDirectory(v *nfs.Target, path string) {
	if _, err := v.Mkdir(path, 0775); err != nil {
		log.Fatalf("mkdir error: %v", err)
	}
}

func removeDirectory(v *nfs.Target, path string) {
	if err := v.RmDir(path); err != nil {
		log.Fatalf("mkdir error: %v", err)
	}
}

func main() {
	log.SetLevel(log.DebugLevel)
	log.Infof(strconv.Itoa(len(os.Args)))
	if len(os.Args) <= 3 {
		log.Infof("%s <host>:<target path> <access level root:0:0> <command ls/up/down/rm/mkdir/rmdir> <path if required> <dest if upload>", os.Args[0])
		os.Exit(-1)
	}
	fmt.Println(os.Args)
	b := strings.Split(os.Args[1], ":")
	c := strings.Split(os.Args[2], ":")

	host := b[0]
	target := b[1]
	user := c[0]
	uidstring := c[1]
	gidstring := c[2]
	cmd := os.Args[3]
	path := "."
	dest := ""
	if len(os.Args) > 4 {
		path = os.Args[4]
	}
	if len(os.Args) > 5 {
		dest = os.Args[5]
	}

	// convert uid to int
	uid64, err := strconv.Atoi(uidstring)
	uid := uint32(uid64)
	if err != nil {
		log.Fatalf("UID needs to be an integer - %v", err)
	}

	// convert uid to int
	gid64, err := strconv.ParseUint(gidstring, 10, 32)
	gid := uint32(gid64)
	if err != nil {
		log.Fatalf("GID needs to be an integer - %v", err)
	}

	log.Infof("host=%s target=%s command=%s\n", host, target, cmd)

	// connect
	mount, err := nfs.DialMount(host, true)
	if err != nil {
		log.Fatalf("unable to dial MOUNT service: %v", err)
	}
	defer mount.Close()

	// auth
	auth := rpc.NewAuthUnix(user, uid, gid)

	v, err := mount.Mount(target, auth.Auth())
	if err != nil {
		log.Fatalf("unable to mount volume: %v", err)
	}
	defer v.Close()

	switch cmd {
	case "ls":
		ls(v, path)
	case "down":
		downloadFile(v, path)
	case "up":
		uploadFile(v, path, dest)
	case "rm":
		removeFile(v, path)
	case "mkdir":
		makeDirectory(v, path)
	case "rmdir":
		removeDirectory(v, path)
	default:
		log.Infof("No command given, just running an directory listing")
	}

	if err = mount.Unmount(); err != nil {
		log.Fatalf("unable to umount target: %v", err)
	}
	mount.Close()
	log.Infof("Completed tests")
}
