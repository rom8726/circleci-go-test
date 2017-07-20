package proba

import (
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestApplication_SomeFunc(t *testing.T) {
	convey.Convey("SomeFunc() should work correctly", t, func() {
		app := Application{}
		convey.So(app.SomeFunc(), convey.ShouldEqual, 3)
	})
}
