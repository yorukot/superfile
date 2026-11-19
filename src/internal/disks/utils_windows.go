package disks

import (
	"fmt"
	"slices"

	"github.com/google/glazier/go/storage"
)

// Get disks without partitions, if any, and append them to disks.
// disks is a pointer to the caller's slice.
func getPartitionlessDisks(svc *storage.Service, disks *[]Disk) error {
	d, err := svc.GetDisks("WHERE NumberOfPartitions = 0")
	if err != nil {
		return err
	}
	defer d.Close()

	for i := range d.Disks {
		disk := &d.Disks[i]

		appendDisk(disks, Disk{
			Name:     disk.FriendlyName,
			Internal: disk.IsSystem || disk.IsBoot || disk.BootFromDisk,
		})
	}

	return nil
}

// Get volumes and append them to disks. disks is a pointer to the caller's slice.
func getVolumes(svc *storage.Service, disks *[]Disk) error {
	v, err := svc.GetVolumes("")
	if err != nil {
		return err
	}
	defer v.Close()

	p, err := svc.GetPartitions("")
	if err != nil {
		return err
	}
	defer p.Close()

	for i := range v.Volumes {
		volume := &v.Volumes[i]

		if shouldHideVolume(p.Partitions, *volume) {
			continue
		}

		name, mountPoint := getVolumeName(volume), getVolumeMountPoint(volume)

		appendDisk(disks, Disk{
			Name:       name,
			MountPoint: mountPoint,
			Mounted:    volume.DriveLetter != "" || volume.Path != "",
		})
	}

	return nil
}

func appendDisk(disks *[]Disk, disk Disk) {
	*disks = append(*disks, disk)
}

func getVolumeName(v *storage.Volume) string {
	name := v.FileSystemLabel

	if name != "" && v.DriveLetter != "" && v.DriveLetter != "\x00" {
		name = fmt.Sprintf("%s (%s)", name, v.DriveLetter)
	} else if name == "" && v.DriveLetter != "" && v.DriveLetter != "\x00" {
		name = v.DriveLetter + ":"
	} else if name == "" && v.Path != "" {
		name = v.Path
	}

	return name
}

func getVolumeMountPoint(v *storage.Volume) string {
	mountPoint := ""

	if v.DriveLetter != "" && v.DriveLetter != "\x00" {
		mountPoint = v.DriveLetter + `:\`
	} else if v.Path != "" {
		mountPoint = v.Path
	}

	return mountPoint
}

func shouldHideVolume(partitions []storage.Partition, volume storage.Volume) bool {
	if volume.Size == 0 {
		return true
	}

	for _, p := range partitions {
		if len(p.AccessPaths) == 0 {
			continue
		}
		if slices.Contains(p.AccessPaths, volume.Path) && p.IsHidden {
			return true
		}
	}

	return false
}
