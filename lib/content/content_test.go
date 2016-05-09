package content

/*
var (
	GetSourceURL = getSourceURL
)

var urlPath = "/example/main_l.jpg"
var urlPath2 = "/test/test/main_l.jpg"
var urlPath3 = "/rnc/img/1617/o/0013076053.jpg"
var urlPath4 = "/abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890~`!@#$%^&*()-_+={}[]|\\;:'\"`,.<>?//test.jpg"

func isNil(i interface{}) bool {
	return i == nil
}

func TestGetSourceURL0(t *testing.T) {
	t.Log("何も考えずにURLを置き換える時のテスト")
	conf := configure.NewConfigure("../test/test_conf.json")
	if GetSourceURL(urlPath, conf) != "http://example.com/main_l.jpg" {
		t.Fatalf("source URL is \"%v\".", GetSourceURL(urlPath, conf))
	}
	if GetSourceURL("", conf) != "" {
		t.Fatalf("source URL is \"%v\".", GetSourceURL("", conf))
	}
	if GetSourceURL(urlPath, nil) != "" {
		t.Fatalf("source URL is \"%v\".", GetSourceURL(urlPath, nil))
	}
	if GetSourceURL("", nil) != "" {
		t.Fatalf("source URL is \"%v\".", GetSourceURL("", nil))
	}
}

func TestGetSourceURL1(t *testing.T) {
	t.Log("置き換えるURLに\"/\"が入っている時のテスト")
	conf := configure.NewConfigure("../test/test_conf.json")
	if GetSourceURL(urlPath2, conf) != "http://example.com/main_l.jpg" {
		t.Log("http://example.com/main_l.jpg")
		t.Fatalf("source URL is \"%v\".", GetSourceURL(urlPath, conf))
	}
	if GetSourceURL("", conf) != "" {
		t.Fatalf("source URL is \"%v\".", GetSourceURL("", conf))
	}
	if GetSourceURL(urlPath2, nil) != "" {
		t.Fatalf("source URL is \"%v\".", GetSourceURL(urlPath, nil))
	}
	if GetSourceURL("", nil) != "" {
		t.Fatalf("source URL is \"%v\".", GetSourceURL("", nil))
	}
}

func TestGetSourceURL2(t *testing.T) {
	t.Log("置き換えるURLに\"/\"が入っている時のテスト")
	t.Log("合致するformatがなかったパターン")
	conf := configure.NewConfigure("../test/test_conf3.json")
	if GetSourceURL(urlPath, conf) != "" {
		t.Fatalf("source URL is \"%v\".", GetSourceURL(urlPath, conf))
	}
	if GetSourceURL("", conf) != "" {
		t.Fatalf("source URL is \"%v\".", GetSourceURL("", conf))
	}
	if GetSourceURL(urlPath, nil) != "" {
		t.Fatalf("source URL is \"%v\".", GetSourceURL(urlPath, nil))
	}
	if GetSourceURL("", nil) != "" {
		t.Fatalf("source URL is \"%v\".", GetSourceURL("", nil))
	}
}

func TestGetSourceURL3(t *testing.T) {
	t.Log("置き換えるURLに\"/\"が入っている時のテスト")
	t.Log("合致するformatがなかったパターン その２")
	conf := configure.NewConfigure("../test/test_conf4.json")
	if GetSourceURL(urlPath2, conf) != "" {
		t.Fatalf("source URL is \"%v\".", GetSourceURL(urlPath2, conf))
	}
	if GetSourceURL("", conf) != "" {
		t.Fatalf("source URL is \"%v\".", GetSourceURL("", conf))
	}
	if GetSourceURL(urlPath2, nil) != "" {
		t.Fatalf("source URL is \"%v\".", GetSourceURL(urlPath2, nil))
	}
	if GetSourceURL("", nil) != "" {
		t.Fatalf("source URL is \"%v\".", GetSourceURL("", nil))
	}
}

func TestGetSourceURL4(t *testing.T) {
	t.Log("strings.Trimできないパターン")
	conf := configure.NewConfigure("../test/test_conf.json")
	if GetSourceURL(urlPath3, conf) != "" {
		t.Fatalf("source URL is \"%v\".", GetSourceURL(urlPath3, conf))
	}
}

func TestGetSourceURL5(t *testing.T) {
	t.Log("strings.Trimできないパターン")
	conf := configure.NewConfigure("../test/test_conf.json")
	if GetSourceURL(urlPath3, conf) != "" {
		t.Fatalf("source URL is \"%v\".", GetSourceURL(urlPath4, conf))
	}
}
*/
