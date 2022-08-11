package dict

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDict(t *testing.T) {
	Convey("test public apis", t, func() {
		Convey("test new dict", func() {
			d := NewDict()
			So(d.Len(), ShouldEqual, 0)
			So(d.Cap(), ShouldEqual, 0)
			So(d.rehashIdx, ShouldEqual, -1)
			So(len(d.hashTables), ShouldEqual, 2)
		})

		Convey("test store", func() {
			d := NewDict()
			d.Store("hello", "world")
			So(d.Len(), ShouldEqual, 1)
			So(d.Cap(), ShouldEqual, _initialHashtableSize)

			d.Store(1, 2)
			d.Store(3, 4)
			d.Store(4, 5)
			d.Store(5, 6)

			// should expand size now
			So(d.Cap(), ShouldEqual, _initialHashtableSize*2)
			So(d.isRehashing(), ShouldBeTrue)
		})

		Convey("test load", func() {
			d := NewDict()
			d.Store(1, 2)
			d.Store(3, 4)

			v, ok := d.Load(1)
			So(v, ShouldEqual, 2)
			So(ok, ShouldBeTrue)

			v, ok = d.Load(100)
			So(v, ShouldBeNil)
			So(ok, ShouldBeFalse)
		})

		Convey("test load or store", func() {
			d := NewDict()
			v, loaded := d.LoadOrStore(1, 1)
			So(v, ShouldEqual, 1)
			So(loaded, ShouldBeFalse)

			v, loaded = d.LoadOrStore(1, 1)
			So(v, ShouldEqual, 1)
			So(loaded, ShouldBeTrue)
		})

		Convey("test delete", func() {
			d := NewDict()
			d.Store(1, 2)
			So(d.Len(), ShouldEqual, 1)

			d.Delete(1)
			So(d.Len(), ShouldEqual, 0)
		})

		Convey("test resize", func() {
			d := NewDict()
			for i := 0; i < 100; i++ {
				d.Store(i, i)
			}
			So(d.Len(), ShouldEqual, 100)
			So(d.Cap(), ShouldEqual, 128)

			// delete half keys
			for i := 0; i < 100; i += 2 {
				d.Delete(i)
			}

			So(d.Len(), ShouldEqual, 50)
			So(d.Cap(), ShouldEqual, 128)

			// resize it
			So(d.Resize(), ShouldBeNil)
			So(d.Len(), ShouldEqual, 50)
			So(d.Cap(), ShouldEqual, 64)
		})

		Convey("test rehash for a while", func() {
			d := NewDict()
			for i := 0; i < 100; i++ {
				d.Store(i, i)
			}

			So(d.Len(), ShouldEqual, 100)
			So(d.Cap(), ShouldEqual, 128)

			So(d.hashTables[0].used, ShouldBeGreaterThan, 0)
			// some of those buckets has been moved to the ht[1]
			So(d.hashTables[1].used, ShouldBeGreaterThan, 0)

			// make sure rehashing is not finished
			So(d.isRehashing(), ShouldBeTrue)

			d.RehashForAWhile(1 * time.Microsecond)

			// rehashing finished
			So(d.isRehashing(), ShouldBeFalse)
			So(d.hashTables[0].used, ShouldEqual, 100)
			So(d.hashTables[1].used, ShouldEqual, 0)
		})
	})
}
