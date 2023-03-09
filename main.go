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
package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}

	router := gin.Default()

	router.GET("/poweronoff", irsendHandler(viper.GetString("remote_name"), viper.GetString("poweronoff_command")))
	router.GET("/poweron", irsendHandler(viper.GetString("remote_name"), viper.GetString("poweron_command")))
	router.GET("/poweroff", irsendHandler(viper.GetString("remote_name"), viper.GetString("poweroff_command")))
	router.GET("/volumeup", irsendHandler(viper.GetString("remote_name"), viper.GetString("volumeup_command")))
	router.GET("/volumedown", irsendHandler(viper.GetString("remote_name"), viper.GetString("volumedown_command")))
	router.GET("/sleep", irsendHandler(viper.GetString("remote_name"), viper.GetString("sleep_command")))
	router.GET("/mute", irsendHandler(viper.GetString("remote_name"), viper.GetString("mute_command")))
	router.GET("/sourceA1", irsendHandler(viper.GetString("remote_name"), viper.GetString("sourceA1_command")))
	router.GET("/sourceA2", irsendHandler(viper.GetString("remote_name"), viper.GetString("sourceA2_command")))
	router.GET("/sourceA3", irsendHandler(viper.GetString("remote_name"), viper.GetString("sourceA3_command")))
	router.GET("/sourceA4", irsendHandler(viper.GetString("remote_name"), viper.GetString("sourceA4_command")))
	router.GET("/sourcecycle", irsendHandler(viper.GetString("remote_name"), viper.GetString("sourcecycle_command")))
	router.GET("/destAB", irsendHandler(viper.GetString("remote_name"), viper.GetString("destAB_command")))
	router.GET("/destA", irsendHandler(viper.GetString("remote_name"), viper.GetString("destA_command")))
	router.GET("/destA1", irsendHandler(viper.GetString("remote_name"), viper.GetString("destA1_command")))
	router.GET("/destB", irsendHandler(viper.GetString("remote_name"), viper.GetString("destB_command")))
	router.GET("/destB1", irsendHandler(viper.GetString("remote_name"), viper.GetString("destB1_command")))
	router.GET("/destB2", irsendHandler(viper.GetString("remote_name"), viper.GetString("destB2_command")))

	localIP, err := getOutboundIP()
	if err != nil {
		panic(fmt.Errorf("Fatal error getting IP address: %s", err))
	}

	localPort := viper.GetString("port")

	server := &http.Server{
		Addr:    localIP + ":" + localPort,
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			panic(fmt.Errorf("Fatal error starting server: %s", err))
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig

	if err := server.Shutdown(context.Background()); err != nil {
		fmt.Println("Error during server shutdown:", err)
	}
}

func irsendHandler(remoteName, command string) gin.HandlerFunc {
	return func(c *gin.Context) {
		cmd := exec.Command("irsend", "SEND_ONCE", remoteName, command)
		if err := cmd.Run(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send IR command"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "IR command sent successfully"})
	}
}

func getOutboundIP() (string, error) {
	connection, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer connection.Close()
	localAddr := connection.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String(), nil
}
