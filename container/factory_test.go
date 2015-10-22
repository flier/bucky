package container

import (
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type factoryInstance struct {
	factoryMethod ObjectFactoryMethod
	objectType    reflect.Type
	singleton     bool
}

func (f *factoryInstance) GetInstance() (interface{}, error) { return f.factoryMethod() }

func (f *factoryInstance) ObjectType() reflect.Type { return f.objectType }

func (f *factoryInstance) IsSingleton() bool { return f.singleton }

func TestFactoryInstanceRegistry(t *testing.T) {
	Convey("create FactoryInstanceRegistry", t, func() {
		r := NewFactoryInstanceRegistry()

		So(r, ShouldNotBeNil)
		So(r.factoryInstanceCache, ShouldNotBeNil)

		obj, exists := r.getCachedObjectForFactoryInstance("name")

		So(obj, ShouldBeNil)
		So(exists, ShouldBeFalse)

		Convey("get prototype object", func() {
			obj, err := r.getObjectFromFactoryInstance(&factoryInstance{
				factoryMethod: func() (interface{}, error) { return 123, nil },
			}, "name")

			So(obj, ShouldEqual, 123)
			So(err, ShouldBeNil)

			obj, err = r.getObjectFromFactoryInstance(&factoryInstance{
				factoryMethod: func() (interface{}, error) { return 456, nil },
			}, "name")

			So(obj, ShouldEqual, 456)
			So(err, ShouldBeNil)
		})

		Convey("get singleton object", func() {
			r.addSingleton("name", nil)

			obj, err := r.getObjectFromFactoryInstance(&factoryInstance{
				factoryMethod: func() (interface{}, error) { return 123, nil },
				singleton:     true,
			}, "name")

			So(obj, ShouldEqual, 123)
			So(err, ShouldBeNil)

			Convey("singleon object should be cached", func() {
				obj, err = r.getObjectFromFactoryInstance(&factoryInstance{
					factoryMethod: func() (interface{}, error) { return 456, nil },
					singleton:     true,
				}, "name")

				So(obj, ShouldEqual, 123)
				So(err, ShouldBeNil)
			})
		})
	})
}
