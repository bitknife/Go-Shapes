# Dev server

    $ ssh wtfserver@hannibal -p 9339


## Nginx WSS (secure WebSockets)

Commands, test and restart nginx

    $ nginx -t
    $ service nginx restart


### Skapa först self-signed på servern:

Just for testing

    sudo openssl req -addext "subjectAltName = DNS:wtf-dev-server.bitknife.se" -x509 -nodes -days 365 -newkey rsa:2048 -keyout /etc/ssl/private/nginx-selfsigned.key -out /etc/ssl/certs/nginx-selfsigned.crt

### Nginx config:
Create /etc/nginx/conf.d/wtf.conf, this one works:

    map $http_upgrade $connection_upgrade {
    default upgrade;
    '' close;
    }
    
    upstream wtfserver {
    server localhost:8888; # appserver_ip:ws_port
    }
    
    server {
    listen 888;
    
        ssl on;
        ssl_certificate /etc/ssl/certs/nginx-selfsigned.crt;
        ssl_certificate_key /etc/ssl/private/nginx-selfsigned.key;
    
    
        location / {
            proxy_pass http://wtfserver;
            proxy_http_version 1.1;
            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection $connection_upgrade;
        }
    }

### Connecting a client
However, it will whine about the self signed certificate. So need to import the certificate into the OS keyring

Visit https://wtf-dev-server.bitknife.se:888/packets from a browser, open/copy the certificate and add it to the
keyring. Change settings to trust it.

So, we should obtain a Lets Encrypt certificate instead.