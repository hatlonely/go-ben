package monitor

import (
	"fmt"
	"testing"
	"time"

	"github.com/hatlonely/go-kit/strx"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGoPsutil(t *testing.T) {
	Convey("TestGoPsutil", t, func() {
		Convey("cpu", func() {
			for i := 0; i < 1; i++ {
				vs, err := cpu.Percent(time.Second, false)
				So(err, ShouldBeNil)
				fmt.Println(strx.JsonMarshalIndent(vs[0]))
			}
		})

		Convey("mem", func() {
			vm, err := mem.VirtualMemory()
			So(err, ShouldBeNil)
			fmt.Println(strx.JsonMarshalIndent(vm))
		})

		Convey("disk", func() {
			du, err := disk.Usage("/")
			So(err, ShouldBeNil)
			fmt.Println(du.Used)
		})

		//Convey("disk io", func() {
		//	io, err := disk.IOCounters()
		//	So(err, ShouldBeNil)
		//	fmt.Println(strx.JsonMarshalIndent(io))
		//})

		Convey("net io", func() {
			io, err := net.IOCounters(true)
			So(err, ShouldBeNil)
			fmt.Println(strx.JsonMarshalIndent(io))
		})
	})
}
