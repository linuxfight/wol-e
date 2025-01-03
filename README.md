# Wol-E

Simple Telegram bot for turning PCs with WoL.

1. OpenRC:
```shell
#!/sbin/openrc-run

name="myapp"
description="My Go Application"
command="/path/to/myapp-amd64-linux"  # Path to your binary
command_background=true
pidfile="/run/${name}.pid"
stdout_log="/var/log/${name}.log"
stderr_log="/var/log/${name}-error.log"

depend() {
  need net
  before apache2  # Example: You can add any dependencies you need here
}

start() {
  ebegin "Starting ${name}"
  start-stop-daemon --start --quiet --background --pidfile ${pidfile} --make-pidfile --exec ${command}
  eend $?
}

stop() {
  ebegin "Stopping ${name}"
  start-stop-daemon --stop --quiet --pidfile ${pidfile}
  eend $?
}

restart() {
  stop
  sleep 1
  start
}
```

```shell
sudo chmod +x /etc/init.d/myapp
sudo rc-update add myapp default
sudo service myapp start
```

2. Systemd
```shell
[Unit]
Description=My Go Application
After=network.target

[Service]
ExecStart=/path/to/myapp-amd64-linux  # Path to your binary
WorkingDirectory=/path/to/working/directory  # Optional: if your app needs a specific working directory
Restart=always
User=your-user-name  # Make sure this matches the user under which you want the service to run
Group=your-group-name  # Optional, specify if necessary
Environment=ENV_VAR_NAME=value  # Optional: Set any environment variables your app may need
PIDFile=/run/myapp.pid
StandardOutput=append:/var/log/myapp.log
StandardError=append:/var/log/myapp-error.log

[Install]
WantedBy=multi-user.target
```

```shell
sudo systemctl daemon-reload
sudo systemctl enable myapp.service
sudo systemctl start myapp.service
```
