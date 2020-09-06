package main

import (
	"fmt"
	"github.com/adelolmo/delauncher/config"
	"github.com/adelolmo/delauncher/crypt"
	"github.com/adelolmo/delauncher/magnet"
	"github.com/adelolmo/delauncher/notifications"
	"github.com/webview/webview"
	"os"
)

var secretKey = []byte{11, 22, 33, 44, 55, 66, 77, 88, 99, 00, 11, 22, 33, 44, 55, 66}
var key = crypt.Key{
	Value: secretKey,
}

var conf = config.NewConfig(key)

func main() {
	switch len(os.Args) {
	case 1:
		configure()
	case 2:
		link, err := magnet.NewLink(os.Args[1])
		if err != nil {
			notifications.Message(err.Error())
		}
		addMagnet(link)
	default:
		fmt.Print("usage: delauncher [MAGNET_LINK]")
		os.Exit(1)
	}
}

func configure() {
	configProperties, err := conf.Get()
	if err != nil {
		notifications.Message(fmt.Sprintf("Unable to read configuration. Error %s", err.Error()))
		os.Exit(1)
	}
	w := webview.New(true)
	defer w.Destroy()
	w.SetTitle("Delauncher")
	w.SetSize(500, 160, webview.HintNone)
	w.Bind("save", func(serverUrl, password string) {
		fmt.Printf("ServerUrl: %s  password: %s\n", serverUrl, password)
		conf.Save(serverUrl, password)
	})
	w.Bind("quit", func() {
		w.Terminate()
	})
	type Config struct {
		ServerUrl, Password string
	}
	w.Bind("showConfig", func() Config {
		return Config{
			ServerUrl: configProperties.ServerUrl,
			Password:  configProperties.Password,
		}
	})
	w.Navigate(`data:text/html,
		<!doctype html>
		<html>
			<head>
				<style>
					body {
						background-color:#e8ebee;
					}
					.container {
						width: 480px;
						clear: both;
					}
					.container input {
						width: 280px;
						clear: both;
						float: right;
					}
					fieldset {
						background-color:#b0d7fa;
					}
					label {
						display: inline-block;
						width: 160px;
						text-align: right;
						vertical-align: sub;
					}
					#buttons {
						float: right;
					}
					button {
						background-color:#bdcfdf;
					}
				</style>
			</head>
			<body>
				<div class="container">
					<form>
						<fieldset>
							<legend>Configuration</legend>
							<div class="block">
								<label for="serverUrl">Deluge server URL:</label>
								<input type="text" name="serverUrl" id="serverUrl">
							</div>
							<br/>
							<div class="block">
								<label for="password">Password:</label>
								<input type="password" name="password" id="password">
							</div>
							<br/>
						</fieldset>
						<div id="buttons">
							<button type="button" onclick="save(document.getElementById('serverUrl').value, document.getElementById('password').value);">Save</button>
							<button type="button" onclick="quit();">Quit</button>
						</div>
					</form>
				</div>
			</body>
			<script>
				window.onload = function() {
					showConfig().then(function(config) {
						document.getElementById('serverUrl').value = config.ServerUrl
						document.getElementById('password').value = config.Password
					});
				};
			</script>
		</html>
	`)
	w.Run()
}

func addMagnet(magnetLink magnet.Link) {
	configValues, err := conf.Get()
	if err != nil {
		notifications.Message(fmt.Sprintf("Unable to read configuration. Error %s", err.Error()))
		os.Exit(1)
	}
	var server = magnet.Server{Url: configValues.ServerUrl, Password: configValues.Password}
	if err := server.Add(magnetLink); err != nil {
		fmt.Printf(err.Error())
		notifications.Message(err.Error())
		os.Exit(2)
	}
	notifications.Message(fmt.Sprintf("Link added:\n%s", magnetLink.Name))
}
