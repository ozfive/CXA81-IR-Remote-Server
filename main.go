package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net"
	// "net/http"
	"os"
	"os/exec"
)

/*
	CXA81 IRSend Commands:

	0000000000000401 SOURCE_A4
	0000000000000403 SOURCE_A3
	0000000000000404 SOURCE_A1
	0000000000000405 SOURCE_A2
	0000000000000408 SOURCE_CYCLE
	0000000000000410 KEY_VOLUMEUP
	0000000000000411 KEY_VOLUMEDOWN
	000000000000040c KEY_SLEEP
	000000000000040d KEY_MUTE_UNMUTE
	000000000000040e POWER_ON
	000000000000040f POWER_OFF
	0000000000000414 DEST_AB
	000000000000041c DEST_A
	000000000000041d DEST_B
	000000000000041e DEST_B1
	0000000000000423 DEST_A1
	0000000000000427 DEST_B2
	0000000000000432 MUTE
	0000000000000433 UNMUTE
	0000000000000434 DISP_ON
	0000000000000435 DISP_OFF
	000000000000064c POWER_ONOFF
	000000000000064e POWERON
	000000000000064f POWEROFF

*/

func main() {

	fmt.Println("CXA81-IR-Remote-Server")

	router := gin.Default()

	// SYSTEM ROUTERS

	// /poweronoff POWER_ONOFF Turns the power on OR off.
	router.GET("/poweronoff", func(c *gin.Context) {
		cmd := "irsend"
		args := []string{"SEND_ONCE", "cambridge", "POWER_ONOFF"}
		if err := exec.Command(cmd, args...).Run(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	})

	// /poweron POWERON Turns the power on.
	router.GET("/poweron", func(c *gin.Context) {
		cmd := "irsend"
		args := []string{"SEND_ONCE", "cambridge", "POWERON"}
		if err := exec.Command(cmd, args...).Run(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	})

	// /poweroff POWEROFF Turns the power on.
	router.GET("/poweroff", func(c *gin.Context) {
		cmd := "irsend"
		args := []string{"SEND_ONCE", "cambridge", "POWEROFF"}
		if err := exec.Command(cmd, args...).Run(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	})

	// /volumeup KEY_VOLUMEUP Turns the volume up + 1.
	router.GET("/volumeup", func(c *gin.Context) {
		cmd := "irsend"
		args := []string{"SEND_ONCE", "cambridge", "KEY_VOLUMEUP"}
		if err := exec.Command(cmd, args...).Run(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	})
	// /volumedown KEY_VOLUMEDOWN  Turns the volume down - 1.
	router.GET("/volumedown", func(c *gin.Context) {
		cmd := "irsend"
		args := []string{"SEND_ONCE", "cambridge", "KEY_VOLUMEDOWN"}
		if err := exec.Command(cmd, args...).Run(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	})

	// SOURCE ROUTERS

	router.GET("/sourceA1", func(c *gin.Context) {
		cmd := "irsend"
		args := []string{"SEND_ONCE", "cambridge", "SOURCE_A1"}
		if err := exec.Command(cmd, args...).Run(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	})

	router.GET("/sourceA2", func(c *gin.Context) {
		cmd := "irsend"
		args := []string{"SEND_ONCE", "cambridge", "SOURCE_A2"}
		if err := exec.Command(cmd, args...).Run(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	})

	router.GET("/sourceA3", func(c *gin.Context) {
		cmd := "irsend"
		args := []string{"SEND_ONCE", "cambridge", "SOURCE_A3"}
		if err := exec.Command(cmd, args...).Run(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	})

	router.GET("/sourceA4", func(c *gin.Context) {
		cmd := "irsend"
		args := []string{"SEND_ONCE", "cambridge", "SOURCE_A4"}
		if err := exec.Command(cmd, args...).Run(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	})

	router.GET("/sourcecycle", func(c *gin.Context) {
		cmd := "irsend"
		args := []string{"SEND_ONCE", "cambridge", "SOURCE_CYCLE"}
		if err := exec.Command(cmd, args...).Run(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	})

	localIP, err := GetOutboundIP()
	if err != nil {
		os.Exit(1)
	}

	localPort := "8081"

	router.Run(localIP + ":" + localPort)
}

// Get preferred outbound ip of this machine
func GetOutboundIP() (string, error) {

	connection, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer connection.Close()

	localAddr := connection.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String(), nil
}
