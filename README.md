# Wol-E

Simple Telegram bot for turning PCs with WoL.
Простой телеграм бот для включения компьютеров в сети с помощью Wake-on-LAN.

## Running
1. Configuration example (```config.yml```) | Образец конфигурации
```yaml
# list of devices, that will be available for the bot
# список устройств, к которым у бота будет доступ
devices:
  - name: 'main'
    ip: '192.168.0.*'
    mac: 'A0:B1:C2:D3:E4:F5'

# telegram bot settings
# настройки бота
bot:
  token: "token from https://t.me/botfather" # токен для работы
  admins: # list of unique telegram IDs of users, that can use the bot | список уникальных ID пользователей, которые смогут пользоваться ботом
    - 1234567890

# application settings
# настройки приложения
settings:
  debug: true # for development | для отладки
  timezone: "GMT+3" # for time in logs | для показа времени в логах
```
2. Command | Команда для запуска
```shell
./wol-e -config ./config.yml
```
3. OpenWRT service | OpenWRT сервис
```
#!/bin/sh /etc/rc.common
USE_PROCD=1  # Enable procd
START=95     # Start order (higher = later)
STOP=01      # Stop order (lower = earlier)

start_service() {
    procd_open_instance
    procd_set_param command /bin/wol-e -config /root/wol-e.yml  # Command to run
    procd_set_param respawn                      # Auto-respawn if crashed
    procd_set_param respawn_retry 5              # Retry 5 times before stopping
    procd_set_param stdout 1                     # Redirect stdout to log
    procd_set_param stderr 1                     # Redirect stderr to log
    procd_set_param user root                  # Run as user "nobody"
    procd_close_instance
}
```
