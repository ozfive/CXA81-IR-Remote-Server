package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {

	// Load configuration
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}

	remoteName := viper.GetString("remote_name")
	localPort := viper.GetString("port")

	// get hostname from OS
	hostname, _ := os.Hostname()

	// load LIRC remote commands
	commands := getIrCommands(remoteName)

	fmt.Println("[Web-server] http://" + hostname + "/")

	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	for _, command := range commands {
		router.GET("/"+strings.ToLower(command), irsendHandler(remoteName, command))
	}

	router.GET("/", func(c *gin.Context) {
		htmlContent, err := generateHTML(commands)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error generating HTML: %s", err)
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(htmlContent))
	})

	server := &http.Server{
		// Addr:    localIP + ":" + localPort,
		Addr:    ":" + localPort,
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(fmt.Errorf("Fatal error starting server: %s", err))
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		fmt.Println("Error during server shutdown:", err)
	}
}

func getIrCommands(remoteName string) []string {
	fmt.Println("[EXEC] irsend LIST " + remoteName + " \"\"")
	cmd := exec.Command("irsend", "LIST", remoteName, "")
	output, err := cmd.Output()
	if err != nil {
		panic(fmt.Errorf("Error listing IR commands: %s %s", strings.Replace(string(output), "\n\n", "\n", -1), strings.Replace(err.Error(), "\n\n", "", -1)))
	}

	lines := strings.Split(string(output), "\n")
	var commands []string
	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) == 2 {
			commands = append(commands, parts[1])
		}
	}
	return commands
}

func irsendHandler(remoteName, command string) gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("[EXEC] irsend SEND_ONCE " + remoteName + " " + command)
		cmd := exec.Command("irsend", "SEND_ONCE", remoteName, command)
		if err := cmd.Run(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send IR command"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "IR command sent successfully"})
	}
}

func generateHTML(commands []string) (string, error) {
	const tpl = `
<!DOCTYPE html>
<html>
<head>
    <title>IR Commands</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 0;
            padding: 20px;
            background-color: #f4f4f4;
        }
        h1 {
            text-align: center;
        }
        .container {
            display: flex;
            flex-wrap: wrap;
            justify-content: center;
        }
        .button {
            display: block;
            width: 150px;
            margin: 10px;
            padding: 15px;
            text-align: center;
            background-color: #007BFF;
            color: white;
            text-decoration: none;
            border-radius: 5px;
            transition: background-color 0.3s;
        }
        .button:hover {
            background-color: #0056b3;
        }
    </style>
</head>
<body>
    <h1>IR Commands</h1>
    <div class="container">
        {{range $index, $element := .}}
        {{if eq (mod $index 3) 0}}
        </div><div class="container">
        {{end}}
        <a href="/{{$element | ToLower}}" class="button">{{$element}}</a>
        {{end}}
    </div>
</body>
</html>
`
	t, err := template.New("index").Funcs(template.FuncMap{"ToLower": strings.ToLower, "mod": func(i, j int) int { return i % j }}).Parse(tpl)
	if err != nil {
		return "", err
	}

	var htmlContent strings.Builder
	if err := t.Execute(&htmlContent, commands); err != nil {
		return "", err
	}

	return htmlContent.String(), nil
}
