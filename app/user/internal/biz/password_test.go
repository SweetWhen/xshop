package biz

import "testing"

func TestPwEncode(t *testing.T) {
	impl := NewPWEncode()
	orgPw := "admin123"
	// encoded := impl.Encode(orgPw)
	encoded := "CzBHBtFl0T5nY2k2$df2262b7425360283312ed9d84763b497bf919611edd83fac5e7aae0c81e05db"
	t.Logf("len(encoded):%d,encoded:%s \n", len(encoded), encoded)
	err := impl.Decode(orgPw, encoded)
	if err != nil {
		t.Fatalf("Decode err:%s, orgPw:%s, encoded:%s\n", err.Error(), orgPw, encoded)
	}
}
