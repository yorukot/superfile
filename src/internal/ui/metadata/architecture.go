package metadata

import (
	"debug/elf"
	"debug/macho"
	"debug/pe"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
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

// binaryFormat is a cheap classification based on a file's leading magic bytes.
type binaryFormat int

const (
	formatUnknown binaryFormat = iota
	formatELF
	formatPE
	formatMachO
)

// detectBinaryFormat reads at most the first 4 bytes of the file and decides
// which (if any) binary parser is worth invoking.
//
// This gate matters: debug/pe accepts files WITHOUT an "MZ" header as raw COFF
// objects, so calling pe.Open on an arbitrary file (e.g. a large video)
// misinterprets its leading bytes as section/symbol counts and can read and
// allocate memory proportional to the file's size before failing (issue #1550).
// The debug/* parsers are documented as unsuitable for adversarial/arbitrary
// inputs, so only hand them files whose magic bytes plausibly match.
func detectBinaryFormat(filePath string) binaryFormat {
	f, err := os.Open(filePath)
	if err != nil {
		return formatUnknown
	}
	defer f.Close()

	var magic [4]byte
	if _, err := io.ReadFull(f, magic[:]); err != nil {
		return formatUnknown
	}

	switch {
	case magic == [4]byte{0x7f, 'E', 'L', 'F'}:
		return formatELF
	case magic[0] == 'M' && magic[1] == 'Z':
		return formatPE
	case isCOFFMachine(binary.LittleEndian.Uint16(magic[:2])):
		// Raw COFF objects (no MZ stub) begin directly with a known machine type.
		return formatPE
	case isMachOMagic(binary.BigEndian.Uint32(magic[:])):
		return formatMachO
	default:
		return formatUnknown
	}
}

// isCOFFMachine reports whether m is a COFF machine type superfile can display.
func isCOFFMachine(m uint16) bool {
	switch m {
	case pe.IMAGE_FILE_MACHINE_I386, pe.IMAGE_FILE_MACHINE_AMD64,
		pe.IMAGE_FILE_MACHINE_ARM, pe.IMAGE_FILE_MACHINE_ARMNT, pe.IMAGE_FILE_MACHINE_ARM64,
		pe.IMAGE_FILE_MACHINE_RISCV32, pe.IMAGE_FILE_MACHINE_RISCV64, pe.IMAGE_FILE_MACHINE_RISCV128:
		return true
	default:
		return false
	}
}

// Byte-swapped forms of the Mach-O magics, produced when the file was written
// for the opposite byte order to the reader. Apple's headers call these the
// CIGAM values ("magic" reversed); the standard library only names the native
// forms, so they are spelled out here rather than left as bare literals.
const (
	machoCigam32  uint32 = 0xcefaedfe // macho.Magic32, bytes reversed
	machoCigam64  uint32 = 0xcffaedfe // macho.Magic64, bytes reversed
	machoCigamFat uint32 = 0xbebafeca // macho.MagicFat, bytes reversed
)

// isMachOMagic reports whether m is a Mach-O (thin or universal) magic number,
// in either byte order.
func isMachOMagic(m uint32) bool {
	switch m {
	case macho.Magic32, macho.Magic64, macho.MagicFat,
		machoCigam32, machoCigam64, machoCigamFat:
		return true
	default:
		return false
	}
}

func GetBinaryArchitecture(filePath string) (string, error) {
	switch detectBinaryFormat(filePath) {
	case formatELF:
		if arch, err := getELFArchitecture(filePath); err == nil {
			return arch, nil
		}
	case formatPE:
		if arch, err := getPEArchitecture(filePath); err == nil {
			return arch, nil
		}
	case formatMachO:
		if arch, err := getMachOArchitecture(filePath); err == nil {
			return arch, nil
		}
	case formatUnknown:
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
