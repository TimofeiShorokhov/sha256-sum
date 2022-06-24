# sha256-sum
Implementation of hasher on golang, which computes checksum of files by different algorithms.
Improved algorithms: sha256, sha512, md5.

## How to build 

### Installation of app

`````
local usage: go run cmd/main.go
docker: docker build -t <image_name> .
kubernetes:
kubectl apply -f manifests/configMap.yaml
kubectl apply -f manifests/service.yaml
kubectl apply -f manifests/deploy.yaml
`````
To set path to files/directory for kubernetes change MOUNT_PATH in configMap file
````
MOUNT_PATH: "/root"
````
### Run database
```
kubectl apply -f manifests/postgres-db-secret.yaml
kubectl apply -f manifests/postgres-db-service.yaml
kubectl apply -f manifests/postgres-db-deployment.yaml
```

## How to use app in a local

1. To check checksum of file or files on directory with different algorithm use -d flag , example:
````
go run cmd/main.go -d="path"
or
go run cmd/main.go -d="path/file_name"
````

2. If you want to check if the checksum was changed, use -c flag , example:
````
go run cmd/main.go -c="path"
````

3. To log data from database use -g flag, example:
``````
go run cmd/main.go -g
``````

4. To get help with commands use -h flag, example:
```
go run cmd/main.go -h
```

5. To check if file deleted and update deleted status, use -u flag with path, example:
```
go run cmd/main.go -u="path"
```

6. If you want to change algorithm, change ALG variable in .env file