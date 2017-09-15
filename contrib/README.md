# Installation Notes

## Directory

```
mkdir /opt/cda
chown -R www-data:www-data /opt/cda
mv cda /opt/cda/
chmod +x /opt/cda/cda
```

## Systemd unit
```
mv cda.unit /etc/systemd/system/
systemd daemon-reload
systemd enable cda
systemd start cda
```
