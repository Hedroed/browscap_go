package browscap_go

import (
	"io/ioutil"
	"strings"
	"testing"
)

const (
	TEST_INI_FILE     = "./test-data/full_php_browscap.ini"
	TEST_USER_AGENT   = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/37.0.2062.120 Safari/537.36"
	TEST_IPHONE_AGENT = "Mozilla/5.0 (iPhone; U; CPU iPhone OS 4_3_2 like Mac OS X; en-us) AppleWebKit/533.17.9 (KHTML, like Gecko) Version/5.0.2 Mobile/8H7 Safari/6533.18.5"
)

func TestInitBrowsCap(t *testing.T) {
	if err := InitBrowsCap(TEST_INI_FILE, false); err != nil {
		t.Fatalf("%v", err)
	}
}

func BenchmarkInit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		InitBrowsCap(TEST_INI_FILE, true)
	}
}

func BenchmarkGetBrowser(b *testing.B) {
	data, err := ioutil.ReadFile("test-data/user_agents_sample.txt")
	if err != nil {
		b.Error(err)
	}

	uas := strings.Split(strings.TrimSpace(string(data)), "\n")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		idx := i % len(uas)

		_, ok := GetBrowserData(uas[idx])
		if !ok {
			b.Errorf("User agent not recognized: %s", uas[idx])
		}
	}
}
