package configure

import (
	"fmt"
	"testing"
)

var testConfig = "../test/test_conf.json"

func TestReadConfigure(t *testing.T) {
	conf := NewConfigure(testConfig)
	func() {
		fmt.Println("test conf struct all.")
		if conf == nil {
			t.Fatalf("conf is nil.")
		}
	}()

	func() {
		fmt.Println("port setting test.")
		if conf.Port() != 8080 {
			t.Fatalf("port is not equal 8080, value is \"%v\"", conf.Port())
		}
	}()

	func() {
		fmt.Println("max size test")
		w, h := conf.MaxSize()
		if w != 5000 {
			t.Fatalf("value is %v", w)
		}
		if h != 5000 {
			t.Fatalf("value is %v", h)
		}
	}()

	func() {
		fmt.Println("use_server_timing test")
		ok := conf.UseServerTiming()
		if !ok {
			t.Fatalf("value is %v", ok)
		}
	}()

	func() {
		fmt.Println("enable_metrics_endpoint test")
		ok := conf.EnableMetricsEndpoint()
		if !ok {
			t.Fatalf("value is %v", ok)
		}
	}()

	func() {
		fmt.Println("max_clients test")
		n := conf.MaxClients()
		if n != 50 {
			t.Fatalf("value is %d", n)
		}
	}()
}
