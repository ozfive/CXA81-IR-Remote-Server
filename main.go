package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
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
	buttonsExcludeRegex := viper.GetString("buttons_exclude_regex")

	// get hostname from OS
	hostname, _ := os.Hostname()

	// load LIRC remote commands
	var commands []string
	if buttonsExcludeRegex != "" {
		commands = getIrCommands(remoteName, &buttonsExcludeRegex)
	} else {
		commands = getIrCommands(remoteName, nil)
	}

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

func getIrCommands(remoteName string, buttonsExcludeRegex *string) []string {
	fmt.Println("[EXEC] irsend LIST " + remoteName + " \"\"")
	cmd := exec.Command("irsend", "LIST", remoteName, "")
	output, err := cmd.Output()
	if err != nil {
		panic(fmt.Errorf("Error listing IR commands: %s %s", strings.Replace(string(output), "\n\n", "\n", -1), strings.Replace(err.Error(), "\n\n", "", -1)))
	}

	fmt.Printf("Output:\n%s\n", output)

	var excludeRegex *regexp.Regexp
	if buttonsExcludeRegex != nil && *buttonsExcludeRegex != "" {
		fmt.Printf("Compiling regex: %s\n", *buttonsExcludeRegex)
		excludeRegex, err = regexp.Compile(*buttonsExcludeRegex)
		if err != nil {
			panic(fmt.Errorf("Error compiling regex: %s", err))
		}
	}

	lines := strings.Split(string(output), "\n")
	var commands []string
	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) == 2 {
			if excludeRegex == nil || !excludeRegex.MatchString(parts[1]) {
				commands = append(commands, parts[1])
			}
		}
	}
	return commands
}

func irsendHandler(remoteName, command string) gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("[EXEC] irsend SEND_ONCE " + remoteName + " " + command)
		cmd := exec.Command("irsend", "SEND_ONCE", remoteName, command)
		err := cmd.Run()
		if err != nil {
			c.JSON(500, gin.H{"status": "error", "message": fmt.Sprintf("Error executing command %s: %s", command, err)})
			return
		}
		c.JSON(200, gin.H{"status": "success", "command": command})
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
            cursor: pointer;
        }
        .button:hover {
            background-color: #0056b3;
        }
        #toast {
            visibility: hidden;
            min-width: 250px;
            margin-left: -125px;
            background-color: #333;
            color: #fff;
            text-align: center;
            border-radius: 2px;
            padding: 16px;
            position: fixed;
            z-index: 1;
            left: 50%;
            bottom: 30px;
            font-size: 17px;
        }
        #toast.show {
            visibility: visible;
            -webkit-animation: fadein 0.5s, fadeout 0.5s 2.5s;
            animation: fadein 0.5s, fadeout 0.5s 2.5s;
        }
        @-webkit-keyframes fadein {
            from {bottom: 0; opacity: 0;}
            to {bottom: 30px; opacity: 1;}
        }
        @keyframes fadein {
            from {bottom: 0; opacity: 0;}
            to {bottom: 30px; opacity: 1;}
        }
        @-webkit-keyframes fadeout {
            from {bottom: 30px; opacity: 1;}
            to {bottom: 0; opacity: 0;}
        }
        @keyframes fadeout {
            from {bottom: 30px; opacity: 1;}
            to {bottom: 0; opacity: 0;}
        }
    </style>
    <script>
        function showToast(message) {
            var x = document.getElementById("toast");
            x.className = "show";
            x.textContent = message;
            setTimeout(function(){ x.className = x.className.replace("show", ""); }, 3000);
        }

        function sendCommand(command) {
            fetch('/' + command.toLowerCase(), {
                method: 'GET'
            })
            .then(response => response.json())
            .then(data => {
                if (data.status === "success") {
                    showToast("Command " + command + " executed successfully.");
                } else {
                    showToast("Error executing command " + command);
                }
            })
            .catch(error => showToast("Error: " + error));
        }
    </script>
</head>
<body>
    <h1>IR Commands</h1>
    <div class="container">
        {{range $index, $element := .}}
        {{if eq (mod $index 3) 0}}
        </div><div class="container">
        {{end}}
        <div class="button" onclick="sendCommand('{{$element}}')">{{$element}}</div>
        {{end}}
    </div>
    <div id="toast"></div>
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
