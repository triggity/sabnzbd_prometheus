# Sabnzbd Prometheus

A prometheus exporter for Sabnzbd. This will gather metrics from the configured Sabnzbd AP and make them available for prometheus, 

## Usage

This can be run via docker or golang binary

### Docker-compose

Sample `docker-compose.yaml` section

```
sabnzbd_prometheus:
  image: triggity/sabnzbd_prometheus
  container_name: sabnzbd_prometheus
  restart: always
  links:
    - sabnzbd
  ports:
    - 8081:8081
  environment:
    - SABNZBD_URI=http://sabnzbd.mydomain.com:8080
    - SABNZBD_APIKEY=*************
```

### Docker Cli

```

docker run --name sabnzbd_prometheus -e SABNZBD_URI=http://sabnzbd.mydomain.com:8080 -e SABNZBD_APIKEY=************* --restart unless-stopped -p 8081:8081 -d triggity/sabnzbd_prometheus:latest 

```


## Configuration

| Environment Variable | Cli Flag | Default | Description |
| -------------------- | -------- | ------- | ----------- |
| SABNZBD_URI          | -sabnzbd_uri | | Uri to connect to sabnzbd. in the format of `http://ip:port` 
| SABNZBD_APIKEY       | -sabnzbd_apiKey | | apiKey to connect to Sabnzbd
| LISTEN_ADDRESS       | -listen-address | http://0.0.0.0:8081 | Address for server to listen on

