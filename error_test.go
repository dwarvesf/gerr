package gerr

import (
	"testing"
)

var errBadReqTest = New(400, "Bad request")
var errDetailFirstTest = New(10004, "first error")
var errDetailSecondTest = New(10005, "second error")

func errTestHandler() error {
	return doHandler()
}

func doHandler() error {
	err := doDetailHandlerLvl1()
	return errBadReqTest.Err(Trace(err), err)
}

func doDetailHandlerLvl1() error {
	err := callDB()
	err2 := callDB()
	return errDetailFirstTest.Err(err, err2)
}

func callDB() error {
	return errDetailSecondTest.Err()
}

func Test_handler(t *testing.T) {

	var wantErr error
	wantErr = E(400, "Bad request")

	if err := errTestHandler(); (err != nil) && err.Error() == wantErr.Error() {
		t.Errorf("handler() error = \n %+s", err)
	}

}
