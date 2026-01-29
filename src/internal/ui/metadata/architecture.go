package metadata

import (
	"debug/elf"
	"debug/macho"
	"debug/pe"
	"errors"
	"fmt"
	"strings"
)

const (
	archI386    = "i386"
	archX8664   = "x86-64"
	archARM     = "ARM"
	archARM64   = "ARM64"
	archPPC     = "PowerPC"
	archPPC64   = "PowerPC64"
	archRISCV   = "RISC-V"
	archS390x   = "s390x"
	archSPARC64 = "SPARC64"
	archMIPS    = "MIPS"
)

var errNotBinary = errors.New("not a recognized binary format")

func GetBinaryArchitecture(filePath string) (string, error) {
	if arch, err := getELFArchitecture(filePath); err == nil {
		return arch, nil
	}

	if arch, err := getPEArchitecture(filePath); err == nil {
		return arch, nil
	}

	if arch, err := getMachOArchitecture(filePath); err == nil {
		return arch, nil
	}

	return "", errNotBinary
}

func getELFArchitecture(filePath string) (string, error) {
	f, err := elf.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	arch := elfMachineToString(f.Machine)
	return fmt.Sprintf("ELF %s", arch), nil
}

func getPEArchitecture(filePath string) (string, error) {
	f, err := pe.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	arch := peArchitectureToString(f.Machine)
	return fmt.Sprintf("PE %s", arch), nil
}

func getMachOArchitecture(filePath string) (string, error) {
	f, err := macho.Open(filePath)
	if err == nil {
		defer f.Close()
		arch := machoCPUToString(f.Cpu)
		return fmt.Sprintf("Mach-O %s", arch), nil
	}

	fat, err := macho.OpenFat(filePath)
	if err != nil {
		return "", err
	}
	defer fat.Close()

	archs := make([]string, 0, len(fat.Arches))
	for _, arch := range fat.Arches {
		archs = append(archs, machoCPUToString(arch.Cpu))
	}

	if len(archs) == 1 {
		return fmt.Sprintf("Mach-O %s", archs[0]), nil
	}
	return fmt.Sprintf("Mach-O Universal (%s)", strings.Join(archs, ", ")), nil
}

//nolint:exhaustive // common architectures only
func elfMachineToString(machine elf.Machine) string {
	switch machine {
	case elf.EM_386:
		return archI386
	case elf.EM_X86_64:
		return archX8664
	case elf.EM_ARM:
		return archARM
	case elf.EM_AARCH64:
		return archARM64
	case elf.EM_MIPS:
		return archMIPS
	case elf.EM_PPC:
		return archPPC
	case elf.EM_PPC64:
		return archPPC64
	case elf.EM_RISCV:
		return archRISCV
	case elf.EM_S390:
		return archS390x
	case elf.EM_SPARCV9:
		return archSPARC64
	default:
		return machine.String()
	}
}

func peArchitectureToString(machine uint16) string {
	switch machine {
	case pe.IMAGE_FILE_MACHINE_I386:
		return archI386
	case pe.IMAGE_FILE_MACHINE_AMD64:
		return archX8664
	case pe.IMAGE_FILE_MACHINE_ARM:
		return archARM
	case pe.IMAGE_FILE_MACHINE_ARM64:
		return archARM64
	default:
		return fmt.Sprintf("Unknown (0x%x)", machine)
	}
}

func machoCPUToString(cpu macho.Cpu) string {
	switch cpu {
	case macho.Cpu386:
		return archI386
	case macho.CpuAmd64:
		return archX8664
	case macho.CpuArm:
		return archARM
	case macho.CpuArm64:
		return archARM64
	case macho.CpuPpc:
		return archPPC
	case macho.CpuPpc64:
		return archPPC64
	default:
		return cpu.String()
	}
}
