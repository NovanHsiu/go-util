package test_tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/NovanHsiu/goutil"
	"github.com/NovanHsiu/goutil/cipher"
	"github.com/NovanHsiu/goutil/network"
)

func TestUtil(t *testing.T) {
	sqlnowtime := goutil.SQLTimeFormatToString(time.Now())
	t.Errorf("SQLTimeFormatToString failed! SQL Now Time: %s", sqlnowtime)
}

func TestErrorHandler(t *testing.T) {
	fmt.Println(goutil.CreateResponse("200.1"))
}

func TestHttpRequest(t *testing.T) {
	fmt.Println(network.RequestTimeOutSecond)
}

func TestCipher(t *testing.T) {
	cp := cipher.DefaultCipher()
	passwd := cp.EncodePassword("abcdefg")
	t.Errorf("TestCipher failed! password: %s", passwd)
}
