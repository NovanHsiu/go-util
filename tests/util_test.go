package test_tests

import (
	"fmt"
	"testing"
	"time"
	"util"
)

func TestUtil(t *testing.T) {
	sqlnowtime := util.SQLTimeFormatToString(time.Now())
	t.Errorf("SQLTimeFormatToString test failed! SQL Now Time: %s", sqlnowtime)
}

func TestErrorHandler(t *testing.T) {
	fmt.Println(util.CreateResponse("200.1"))
}
