# NGINX Load Balancer Updater API

Allows updating a nginx L4 load balancer setup through a JSON API

## How it works
1 / The API receives the updates
Example:
````json
{
    "backendName": "default_myservice",
    "lbPort": 8080,
    "lbProtocol": "tcp",
    "upstreamServers":[
        {
            "host": "192.168.64.5", 
            "port": 30291
        },
        {
            "host": "192.168.64.6", 
            "port": 30291
        }
    ],
    "proxyTimeoutSeconds": 5,
    "proxyConnectTimeoutSeconds": 2,    
}
````

2/ The updated config file is written to disk
Example:
````
# /etc/nginx/conf.f/default_myservice.conf
upstream default_myservice {
    server 192.168.64.5:30291;
    server 192.168.64.6:30291;
}
    
server {
    listen        8080;
    proxy_pass    default_myservice;
    proxy_timeout 5s;
    proxy_connect_timeout 2s;
}
````
You can see an example of a root nginx config that works with this setup under nginx/nginx.conf

3/ We issue an nginx hot reload command through 
````
sudo nginx -s reload
````

4/ Done.  

## Run locally
````bash
docker-compose up
````
This runs a container with the API exposed on the host on port 3000, as well as a container of the nginx instance where load balancers can be created on ports 8080-8084.
The 2 containers share a volume, and the nginx container watches for changes on that volume and reloads the server whenever config files are modified by the nginx-lb-updater API.

## API
### Create/Update LB
POST /api/v1/lb
````json
{
    "backendName": "default_myservice",
    "lbPort": 8080,
    "lbProtocol": "tcp",
    "upstreamServers":[
        {
            "host": "192.168.64.5", 
            "port": 30291
        },
        {
            "host": "192.168.64.6", 
            "port": 30291
        }
    ],
    "proxyTimeoutSeconds": 5,
    "proxyConnectTimeoutSeconds": 2,    
}
````

### Delete LB
DELETE /api/v1/lb
````json
{
    "backendName": "default_myservice"
}
````

## License 

MIT License

Copyright (c) 2023 Adil H

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.