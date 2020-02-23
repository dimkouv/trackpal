# Deploying Trackpal API

> Instructions for deploying trackpal API to your server.


We are going to use nginx and systemd for deployment but first we need a postgres instance.

### Postgres and migrations
TBD

### Systemd

First add a new systemd service for trackpal by creating 
`/etc/systemd/system/trackpal.service`
with the following content:

```service
[Unit]
Description=Trackpal API
After=network.target

[Service]
User=my-username
Group=my-group
WorkingDirectory=/path/to/trackpal
ExecStart=/path/to/trackpal/trackpal
Environment='TRACKPAL_SIGN_KEY=a-key-to-sign-jwt'
Environment='TRACKPAL_MAIL_SENDER=bot@trackpal.xyz'
Environment='TRACKPAL_POSTGRES_DSN=user=dbuser password=dbpw dbname=mydb sslmode=disable'
Environment='TRACKPAL_SERVER_ADDR=127.0.0.1:8000'
Environment='TRACKPAL_SMTP_PORT=587'
Environment='TRACKPAL_SMTP_HOST=my-smtp.server.com'
Environment='TRACKPAL_SMTP_USER=smtp-server-user'
Environment='TRACKPAL_SMTP_PASSWORD=smtp-server-pw'

[Install]
WantedBy=multi-user.target
```

Now you can use the following commands to manage the server
```bash
sudo systemctl start|stop|enable|status trackpal 
```

### Nginx

Next step is to setup Nginx as a reverse proxy. Create 
`/etc/nginx/sites-available/trackpal` with the following content

```nginx
server {
    listen 443 ssl;
    server_tokens off;
    client_max_body_size 2m;

    ssl on;
    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    server_name api.trackpal.xyz;

    location /v1/ {
        proxy_set_header X-Real-IP $remote_addr;
        proxy_pass http://127.0.0.1:8000/;
    }
}
```

```
# validate nginx syntax
sudo nginx -t

# enable server
sudo ln -s /etc/nginx/sites-available/trackpal /etc/nginx/sites-enabled/trackpal

# reload nginx
sudo systemctl reload nginx 
```
