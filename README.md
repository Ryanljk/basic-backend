# Deployment Writeup
The instructions given here shall relate to the deployment of the basic-backend codebase onto an AWS EC2 instance via Docker container. The steps are as follows:

1. SSH into the EC2 instance (either from your local machine or via AWS Cloudshell)
    
        ssh -i /<path>/<key.pem> ubuntu@<ip>

2. Install Docker Engine along with Compose plugin (ensure you run `sudo apt update` beforehand)

3. Copy the backend codebase files into the EC2 instance
    
        scp -i /<path>/<key.pem> /path/basic-backend ubuntu@<ip>:/home/ubuntu/

4. Run Docker Compose to start the application

## Enabling HTTPS
nginx will be used to act as a reverse proxy.

1. Install nginx and run the following commands to start it

        sudo systemctl start nginx
        sudo systemctl enable nginx

2. Install certbot and run with nginx plugin (assuming there is a domain. If there is no domain, then a self-signed certificate has to be used instead via openssl)

        sudo certbot --nginx -d <domain>

        OR

        openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
        -keyout /etc/ssl/private/selfsigned.key \
        -out /etc/ssl/certs/selfsigned.crt \
        -subj "/CN=<ip>"

3. Modify the config file for nginx

        server {
            listen 443 ssl;

            server_name <ip>;

            ssl_certificate /etc/ssl/certs/cert.crt;
            ssl_certificate_key /etc/ssl/private/cert.key;

            location / {
                proxy_pass http://127.0.0.1:8080;
                proxy_set_header Host $host;
                proxy_set_header X-Real-IP $remote_addr;
                proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
                proxy_set_header X-Forwarded-Proto $scheme;
            }
        }

        server {
            listen 80;
            server_name <ip>;
            return 301 https://$host$request_uri;
        }


