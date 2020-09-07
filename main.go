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
	w.SetSize(500, 170, webview.HintNone)
	w.Bind("save", func(serverUrl, password string) {
		conf.Save(serverUrl, password)
	})
	w.Bind("exit", func() {
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
						background-color: #f5f7fa;
						color: #656667;
						font-family: Ubuntu, "times new roman", times, roman, serif;
						font-size: 14px;
					}
					.container {
						width: 480px;
						clear: both;
					}
					.container input {
						color: #656667;
						width: 350px;
						clear: both;
						float: right;
						padding: 5px 5px;
						text-decoration: none;
						display: inline-block;
						font-size: 14px;
					}
					.container input:focus  {
						border: 1px solid orange;
					}
					form {
						margin: 20px 0px 0px 0px;
					}
					fieldset {
						border: none;
					}
					label {
						display: inline-block;
						width: 80px;
						padding: 10px 5px;
						vertical-align: sub;
					}
					#buttons {
						float: right;
					}
					button {						
						color: #656667;
						background-image: linear-gradient(#f5f7fa, #eef0f2);
						border: 1px solid #b3bac0;
						padding: 5px 28px;
						text-align: center;
						text-decoration: none;
						display: inline-block;
						font-size: 14px;
						margin: 1px 1px;
						border-radius: 5px;
					}
					button:hover {
						border: 1px solid #5d5f62;	
						margin: 1px 1px;
					}
					button:active {
						background-image: none;
						background-color: #e9ebee;
					}
				</style>
			</head>
			<body>
				<div class="container">
					<form>
						<fieldset>
							<legend>Deluge configuration</legend>
							<label for="serverUrl">Server URL:</label>
							<input type="text" name="serverUrl" id="serverUrl">
							<label for="password">Password:</label>
							<input type="password" name="password" id="password">
						</fieldset>
						<div id="buttons">
							<button type="button" onclick="exit();">Exit</button>
							<button type="button" onclick="save(document.getElementById('serverUrl').value, document.getElementById('password').value);">Save</button>
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
