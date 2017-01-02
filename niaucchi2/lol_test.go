package niaucchi2

import (
	"io"
	"net"
	"testing"
	"time"
)

func TestLol(t *testing.T) {
	go func() {
		lsnr, _ := net.Listen("tcp", "127.0.0.1:13371")
		for {
			zzz, _ := lsnr.Accept()
			go func() {
				defer zzz.Close()
				cont := NewServerCtx()
				err := cont.Absorb(zzz)
				if err != nil {
					panic(err.Error())
				}
				for {
					clnt, err := cont.Accept()
					if err != nil {
						panic(err.Error())
					}
					go func() {
						defer clnt.Close()
						io.Copy(clnt, clnt)
					}()
				}
			}()
		}
	}()
	time.Sleep(time.Millisecond * 200)
	wire, err := net.Dial("tcp", "127.0.0.1:13371")
	if err != nil {
		panic(err.Error())
	}
	cont := NewClientCtx()
	err = cont.Absorb(wire)
	if err != nil {
		panic(err.Error())
	}
	lsnr, _ := net.Listen("tcp", "127.0.0.1:13370")
	for {
		clnt, _ := lsnr.Accept()
		go func() {
			rmt, err := cont.Tunnel()
			if err != nil {
				panic(err.Error())
			}
			go func() {
				defer rmt.Close()
				defer clnt.Close()
				io.Copy(clnt, rmt)
			}()
			defer rmt.Close()
			defer clnt.Close()
			io.Copy(rmt, clnt)
		}()
	}
}
