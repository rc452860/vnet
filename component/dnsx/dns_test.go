package dnsx

import "testing"

func Test_Reslove(t *testing.T) {
	d := NewDNSWithPrefer("8.8.8.8:53", "114.114.114.114:53", true)
	ip, err := d.Reslove("google.com")
	if err != nil {
		t.Logf("%-v", err)
		t.FailNow()
	}
	t.Log(ip.String())
}

func Test_Reslove_IPv4_False(t *testing.T) {
	d := NewDNSWithPrefer("8.8.8.8:53", "114.114.114.114:53", false)
	ip, err := d.Reslove("www.baidu.com")
	if err != nil {
		t.Logf("%-v", err)
		t.FailNow()
	}
	t.Log(ip.String())
}

func Benchmark_Reslove(t *testing.B) {
	d := NewDNSWithPrefer("114.114.114.114:53", "223.5.5.5:53", true)
	d.Reslove("baidu.com")
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		_, err := d.Reslove("baidu.com")
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
	}
}
