package container

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSingletonRegistry(t *testing.T) {
	Convey("create SingletonRegistry", t, func() {
		r := NewDefaultSingletonRegistry()

		So(r, ShouldNotBeNil)
		So(r.singletonObjects, ShouldNotBeNil)
		So(r.singletonFactories, ShouldNotBeNil)
		So(r.earlySingletonObjects, ShouldNotBeNil)
		So(r.registeredSingletons, ShouldNotBeNil)
		So(r.singletonsCurrentlyInCreation, ShouldNotBeNil)
		So(r.inCreationCheckExclusions, ShouldNotBeNil)

		Convey("register instance with name", func() {
			So(r.ContainsSingleton("name"), ShouldBeFalse)
			So(r.RegisterSingleton("name", t), ShouldBeNil)
			So(r.ContainsSingleton("name"), ShouldBeTrue)
			So(r.SingletonCount(), ShouldEqual, 1)
			So(r.SingletonNames(), ShouldResemble, []string{"name"})

			i, err := r.GetSingleton("name")

			So(i == t, ShouldBeTrue)
			So(err, ShouldBeNil)
		})

		Convey("register instance with invalid name", func() {
			So(r.RegisterSingleton("", t), ShouldNotBeNil)
		})

		Convey("register instance with same name", func() {
			So(r.RegisterSingleton("name", t), ShouldBeNil)
			So(r.RegisterSingleton("name", t), ShouldNotBeNil)
		})

		Convey("add factory with name", func() {
			r.removeSingleton("name")
			So(r.addSingletonFactory("name", ObjectFactoryMethod(func() (interface{}, error) { return 123, nil })), ShouldBeNil)
			So(r.ContainsSingleton("name"), ShouldBeFalse)

			r.beforeSingletonCreation("name")
			So(r.isSingletonCurrentlyInCreation("name"), ShouldBeTrue)
			i, err := r.GetSingleton("name")
			r.afterSingletonCreation("name")
			So(r.isSingletonCurrentlyInCreation("name"), ShouldBeFalse)

			So(i, ShouldEqual, 123)
			So(err, ShouldBeNil)

			So(r.ContainsSingleton("name"), ShouldBeFalse)
		})

		Convey("add factory with invalid name", func() {
			So(r.addSingletonFactory("", ObjectFactoryMethod(func() (interface{}, error) { return t, nil })), ShouldNotBeNil)
			So(r.addSingletonFactory("name", nil), ShouldNotBeNil)
		})
	})
}
