package nfs

import (
	"fmt"
	"io"
	"os"

	"github.com/go-nfs/nfsv3/nfs"
	"github.com/go-nfs/nfsv3/nfs/rpc"
)

type NFSClient struct {
	target *nfs.Target
}

func NewRootAuth() rpc.Auth {
	return rpc.NewAuthUnix("localhost", 0, 0).Auth()
}

func CustomAuth(hostname string, uid, gid int) rpc.Auth {
	return rpc.NewAuthUnix(hostname, uint32(uid), uint32(gid)).Auth()
}

func NewNFSClient(host, basepath string, auth rpc.Auth) (*NFSClient, error) {
	client, err := nfs.DialMount(host, true)
	if err != nil {
		return nil, fmt.Errorf("NewNFSClient failed: %v", err)
	}
	cc, err := client.Mount(basepath, auth)
	return &NFSClient{
		cc,
	}, err
}

func (c *NFSClient) Close() error {
	return c.target.Close()
}

func (c *NFSClient) ListDir(path string) ([]*nfs.EntryPlus, error) {
	return c.target.ReadDirPlus(path)
}

func (c *NFSClient) Stat(fname string) (os.FileInfo, error) {
	info, _, err := c.target.Lookup(fname)
	if err != nil {
		return nil, fmt.Errorf("Stat failed: %v", err)
	}
	return info, nil
}

func (c *NFSClient) Cat(fname string) ([]byte, error) {
	buf, err := c.target.Open(fname)
	if err != nil {
		return nil, fmt.Errorf("Cat failed: %v", err)
	}
	defer buf.Close()
	return io.ReadAll(buf)
}

func (c *NFSClient) Mkdir(fname string) error {
	if _, err := c.target.Mkdir(fname, 0777); err != nil {
		return fmt.Errorf("Mkdir failed: %v", err)
	}
	return nil
}

func (c *NFSClient) Upload(localpath, remotepath string) error {
	f, err := os.Open(localpath)
	if err != nil {
		return fmt.Errorf("Upload failed: %v", err)
	}
	defer f.Close()

	fileinfo, err := f.Stat()
	if err != nil {
		return fmt.Errorf("read fail: %s", err.Error())
	}
	filesize := fileinfo.Size()
	wr, err := c.target.OpenFile(remotepath, 0777)
	if err != nil {
		return fmt.Errorf("write fail: %s", err.Error())
	}

	// Copy filesize
	n, err := io.CopyN(wr, f, int64(filesize))
	if err != nil {
		return fmt.Errorf("error copying: n=%d, %s", n, err.Error())
	}

	if err = wr.Close(); err != nil {
		return fmt.Errorf("error committing: %s", err.Error())
	}
	return nil
}

func (c *NFSClient) Download(remotepath, localpath string) error {
	f, err := os.Create(localpath)
	if err != nil {
		return fmt.Errorf("Download failed: %v", err)
	}
	defer f.Close()

	rdr, err := c.target.Open(remotepath)
	if err != nil {
		return fmt.Errorf("read error: %v", err)
	}
	defer rdr.Close()

	_, err = io.Copy(f, rdr)
	if err != nil {
		return fmt.Errorf("error copying: %s", err.Error())
	}
	return nil
}

func (c *NFSClient) Remove(fname string) error {
	if err := c.target.Remove(fname); err != nil {
		return fmt.Errorf("Remove failed: %v", err)
	}
	return nil
}

func (c *NFSClient) RemoveAll(fname string) error {
	if err := c.target.RemoveAll(fname); err != nil {
		return fmt.Errorf("RemoveAll failed: %v", err)
	}
	return nil
}

func (c *NFSClient) Rename(oldname, newname string) error {
	if err := c.target.Rename(oldname, newname); err != nil {
		return fmt.Errorf("Rename failed: %v", err)
	}
	return nil
}

func (c *NFSClient) RemoveDir(fname string) error {
	if err := c.target.RmDir(fname); err != nil {
		return fmt.Errorf("RemoveDir failed: %v", err)
	}
	return nil
}
