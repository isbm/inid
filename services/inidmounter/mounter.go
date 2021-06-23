package inidmounter

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"
)

type InidMountpoint struct {
	Type   string
	Device string
	Path   string
	Opts   string
}

type InidPremounter struct {
	conf        map[string]map[string]interface{}
	createdDirs []string
}

func NewInidPremounter(conf map[string]map[string]interface{}) *InidPremounter {
	ipm := new(InidPremounter)
	ipm.conf = conf
	ipm.createdDirs = make([]string, 0)
	return ipm
}

func (ipm *InidPremounter) parseOpts(opts string) (uintptr, []string) {
	options := []string{}
	var flags uintptr
	for _, opt := range strings.Split(opts, ",") {
		opt = strings.TrimSpace(opt)
		if opt == "" {
			continue
		}
		if strings.Contains(opt, "=") {
			options = append(options, opt)
		} else {
			switch opt {
			case "nosuid":
				flags |= syscall.MS_NOSUID
			case "noexec":
				flags |= syscall.MS_NOEXEC
			case "nodev":
				flags |= syscall.MS_NODEV
			case "remount":
				flags |= syscall.MS_REMOUNT
			case "ro":
				flags |= syscall.MS_RDONLY
			default:
				log.Printf("Unsupported mounting flag: %s\n", opt)
			}
		}
	}

	return flags | syscall.MS_NOATIME | syscall.MS_SILENT, options
}

// Mount a filesystem
func (ipm *InidPremounter) Mount(device string, fstype string, point string, opts string) error {
	if _, err := os.Stat(point); err != nil {
		err := os.MkdirAll(point, 0755)
		if err != nil {
			return fmt.Errorf("Error creating directory %s: %s", point, err.Error())
		}
		ipm.createdDirs = append(ipm.createdDirs, point)
	}
	flags, options := ipm.parseOpts(opts)

	// Remount?
	for _, o := range options {
		if o == "remount" {
			fstype = ""
			options = []string{}
			device = point
		}
	}

	return syscall.Mount(device, point, fstype, flags, strings.Join(options, ","))
}

// Unmount a filesystem
func (ipm *InidPremounter) Umount(point string) error {
	if err := syscall.Unmount(point, 0); err != nil {
		return err
	}
	// Check if the target is indeed empty
	// Remove the custom directory, pre-created by mount
	// To get those, iterate over the collected
	for _, cd := range ipm.createdDirs {
		log.Printf("Remove dir %s\n", cd)
	}
	return nil
}

func (ipm *InidPremounter) bulk(mount bool) []error {
	errors := []error{}
	for device, conf := range ipm.conf {
		var fstype string
		var mpath string
		var opts string

		// filesystem type, or device (for devpts, cgroup etc)
		if f := conf["type"]; f != nil {
			if fp, ok := f.(string); ok {
				fstype = fp
			}
		}
		if fstype == "" {
			fstype = device
		}

		// mount point
		if p := conf["path"]; p != nil {
			if fp, ok := p.(string); ok {
				mpath = fp
			}
		}
		if mpath == "" {
			errors = append(errors, fmt.Errorf("Unable to find a mount path for device %s", device))
			continue
		}

		// Get options
		if o := conf["opts"]; o != nil {
			if fp, ok := o.(string); ok {
				opts = fp
			}
		}

		if mount {
			if err := ipm.Mount(device, fstype, mpath, opts); err != nil {
				log.Printf("Error mounting %s on %s as %s with %s: %s", device, mpath, fstype, opts, err.Error())
			}
		} else {
			log.Printf("Unmounting path %s", mpath)
		}
	}
	return errors
}

func (ipm *InidPremounter) Start() error {
	errors := ipm.bulk(true)
	if len(errors) > 0 {
		var buff bytes.Buffer
		for idx, err := range errors {
			buff.WriteString(fmt.Sprintf("%d. %s\n", idx+1, err.Error()))
		}
		return fmt.Errorf("Error mounting resources (see below):\n%s", buff.String())
	}
	return nil
}

func (ipm *InidPremounter) Stop() error {
	ipm.bulk(false)
	return nil
}
