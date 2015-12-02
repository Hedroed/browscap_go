package browscap_go

import (
	"bufio"
	"bytes"
	"io"
	"os"
)

var (
	// Ini
	sEmpty   = []byte{}     // empty signal
	nComment = []byte{'#'}  // number signal
	sComment = []byte{';'}  // semicolon signal
	sStart   = []byte{'['}  // section start signal
	sEnd     = []byte{']'}  // section end signal
	sEqual   = []byte{'='}  // equal signal
	sQuote1  = []byte{'"'}  // quote " signal
	sQuote2  = []byte{'\''} // quote ' signal

	// To reduce memory usage we will keep only next keys
	keepKeys = map[string]bool{
		// Required
		"Parent": true,

		// Used in Browser
		"Browser":          true,
		"Version":          true,
		"MajorVer":         true,
		"MinorVer":         true,
		"Browser_Type":     true,
		"Platform":         true,
		"Platform_Version": true,
		"Device_Type":      true,
		//"Device_Code_Name": true,
		//"Device_Brand_Name": true,

		"Comment":      true,
		"Browser_Bits": true,
		//"Browser_Maker": true,
		"Browser_Modus": true,
		//"Platform_Description": true,
		"Platform_Bits": true,
		//"Platform_Maker": true,
		"Alpha": true,
		"Beta":  true,
		//"Win16": true,
		//"Win32": true,
		//"Win64": true,
		"Frames":           true,
		"IFrames":          true,
		"Tables":           true,
		"Cookies":          true,
		"BackgroundSounds": true,
		"JavaScript":       true,
		"VBScript":         true,
		"JavaApplets":      true,
		"ActiveXControls":  true,
		//"isMobileDevice": true,
		//"isTablet": true,
		//"isSyndicationReader": true,
		"Crawler":     true,
		"CssVersion":  true,
		"AolVersion":  true,
		"Device_Name": true,
		//"Device_Maker": true,
		"RenderingEngine_Name":    true,
		"RenderingEngine_Version": true,
		//"RenderingEngine_Description": true,
		//"RenderingEngine_Maker": true,
	}
)

func loadFromIniFile(path string) (*dictionary, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	dict := newDictionary()

	buf := bufio.NewReader(file)
	sectionName := ""

	for {
		line, _, err := buf.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}

		// Empty line
		if bytes.Equal(sEmpty, line) {
			continue
		}

		// Trim
		line = bytes.TrimSpace(line)

		// Empty line
		if bytes.Equal(sEmpty, line) {
			continue
		}

		// Comment line
		if bytes.HasPrefix(line, nComment) || bytes.HasPrefix(line, sComment) {
			continue
		}

		// Section line
		if bytes.HasPrefix(line, sStart) && bytes.HasSuffix(line, sEnd) {
			sectionName = string(line[1 : len(line)-1])
			continue
		}

		// Create section
		if _, ok := dict.mapped[sectionName]; !ok {
			// Save mapped
			dict.mapped[sectionName] = make(section)

			dict.tree.Add(sectionName)

		}

		// Key => Value
		kv := bytes.SplitN(line, sEqual, 2)

		// Parse Key
		key := bytes.TrimSpace(kv[0])
		if ok, in := keepKeys[string(key)]; !ok || !in {
			continue
		}

		// Parse Value
		val := bytes.TrimSpace(kv[1])
		if bytes.HasPrefix(val, sQuote1) {
			val = bytes.Trim(val, `"`)
		}
		if bytes.HasPrefix(val, sQuote2) {
			val = bytes.Trim(val, `'`)
		}

		dict.mapped[sectionName][string(key)] = string(val)
	}

	dict.buildCompleteData()

	return dict, nil
}
