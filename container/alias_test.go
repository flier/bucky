package container

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSimpleAliasRegistry(t *testing.T) {
	Convey("create SimpleAliasRegistry", t, func() {
		r := NewSimpleAliasRegistry()

		So(r, ShouldNotBeNil)
		So(r.aliases, ShouldNotBeNil)
		So(r.AllowAliasOverriding, ShouldBeTrue)

		So(r.GetAliases("test"), ShouldBeNil)

		Convey("register alias with name", func() {
			So(r.IsAlias("alias"), ShouldBeFalse)
			So(r.RegisterAlias("name", "alias"), ShouldBeNil)
			So(r.aliases, ShouldResemble, map[string]string{"alias": "name"})
			So(r.IsAlias("alias"), ShouldBeTrue)
			So(r.IsAlias("name"), ShouldBeFalse)
			So(r.GetAliases("name"), ShouldResemble, []string{"alias"})
		})

		Convey("register invalid alias or name", func() {
			So(r.RegisterAlias("", "alias"), ShouldNotBeNil)
			So(r.RegisterAlias("name", ""), ShouldNotBeNil)
		})

		Convey("register exists alias with other name", func() {
			So(r.RegisterAlias("name", "alias"), ShouldBeNil)
			So(r.RegisterAlias("name2", "alias"), ShouldBeNil)
			So(r.GetAliases("name"), ShouldBeNil)
			So(r.GetAliases("name2"), ShouldResemble, []string{"alias"})
			So(r.aliases, ShouldResemble, map[string]string{"alias": "name2"})

			r.AllowAliasOverriding = false

			So(r.RegisterAlias("name3", "alias"), ShouldNotBeNil)
		})

		Convey("register alias same to name", func() {
			So(r.RegisterAlias("name", "alias"), ShouldBeNil)
			So(r.IsAlias("alias"), ShouldBeTrue)
			So(r.aliases, ShouldResemble, map[string]string{"alias": "name"})
			So(r.RegisterAlias("alias", "alias"), ShouldBeNil)
			So(r.IsAlias("alias"), ShouldBeFalse)
			So(r.aliases, ShouldResemble, map[string]string{})
		})

		Convey("register indirect alias to name", func() {
			So(r.RegisterAlias("name", "direct"), ShouldBeNil)
			So(r.RegisterAlias("direct", "indirect"), ShouldBeNil)
			So(r.GetAliases("name"), ShouldResemble, []string{"direct", "indirect"})
		})

		Convey("register alias circle reference to self", func() {
			So(r.RegisterAlias("a", "b"), ShouldBeNil)
			So(r.RegisterAlias("b", "c"), ShouldBeNil)
			So(r.RegisterAlias("c", "a"), ShouldNotBeNil)
		})

		Convey("register alias and remove it", func() {
			So(r.RegisterAlias("name", "alias"), ShouldBeNil)
			So(r.IsAlias("alias"), ShouldBeTrue)

			So(r.RemoveAlias("alias"), ShouldBeNil)
			So(r.RemoveAlias("alias"), ShouldNotBeNil)
		})
	})
}
