package playbook

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/golang/glog"

	"gitee.com/wisecloud/wise-deploy/database"
)

const (
	clusterTemplate = "cluster-template"
)

var (
	ansibleFile = []string{
		"file",
		"template",
		"scripts",
		"yat",
		"ansible.cfg",
		"clean.ansible",
		"install.ansible",
	}
)

func InstantiateCluster(wd string, cluster *database.Cluster) error {
	oldwd, err := os.Getwd()
	if err != nil {
		return err
	}

	if err := os.Chdir(wd); err != nil {
		return err
	}
	defer os.Chdir(oldwd)

	newFolder := fmt.Sprintf("cluster-%s-%s", cluster.Name, cluster.ID)

	if err := os.RemoveAll(newFolder); err != nil {
		return err
	}

	if err := os.Mkdir(newFolder, 0755); err != nil {
		return err
	}

	for _, com := range cluster.Components {
		if err := os.Mkdir(
			path.Join(newFolder, com.Name+PlaybookSuffix),
			0755,
		); err != nil {
			return err
		}

		for _, f := range ansibleFile {
			p := path.Join(
				clusterTemplate,
				com.Name+PlaybookSuffix,
				f,
			)

			if _, err := os.Stat(p); os.IsNotExist(err) {
				glog.V(4).Infof("%s's %s folder is not exist", com.Name, f)
				continue
			} else if err != nil {
				return fmt.Errorf("check %s's %s folder error: %v", com.Name, f, err)
			}

			if err := os.Symlink(
				path.Join(wd, p),
				path.Join(newFolder, com.Name+PlaybookSuffix, f),
			); err != nil {
				return err
			}
		}

		if err := MkGroupVars(
			path.Join(wd,
				clusterTemplate,
				com.Name+PlaybookSuffix,
			),
			path.Join(wd,
				newFolder,
				com.Name+PlaybookSuffix,
			),
		); err != nil {
			return err
		}
	}

	return nil
}

func MkGroupVars(oldFolder, newFolder string) error {
	oldFolder = path.Join(oldFolder, "group_vars")
	newFolder = path.Join(newFolder, "group_vars")

	if err := os.Mkdir(newFolder, 0755); err != nil {
		return err
	}

	files, err := ioutil.ReadDir(oldFolder)
	if err != nil {
		return err
	}

	for _, f := range files {
		if err := CopyFile(
			path.Join(oldFolder, f.Name()),
			path.Join(newFolder, f.Name()),
		); err != nil {
			return err
		}
	}

	return nil
}

// CopyFile copies a file from src to dst. If src and dst files exist, and are
// the same, then return success. Otherise, attempt to create a hard link
// between the two files. If that fail, copy the file contents from src to dst.
func CopyFile(src, dst string) (err error) {
	sfi, err := os.Stat(src)
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}
	err = copyFileContents(src, dst)

	return
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}
