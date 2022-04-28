# sha256-sum
Implementation of hasher on golang, which computes checksum of files by different algorithms.
Improved algorithms: sha256, sha512, md5.

## How to build 

### Installation of app

`````
go mod download
go build cmd/main.go
`````
To set path to files/directory for docker change .env file
````
TAG=your_local_path
````
### Run database
```
docker-compose build
docker-compose up
```

## How to use app

1. To check checksum of file or files on directory with different algorithm use -d flag with -a flag , example:
````
go run cmd/main.go -d="path" -a="algorithm"
or
go run cmd/main.go -d="path/file_name" -a="algorithm"
via docker
docker-compose run timdb -d=/local/"your_path" -a="algorithm"
````

2. If you want to check if the checksum was changed, use -c flag , example:
````
go run cmd/main.go -c="path" -a="algorithm"
via docker
docker-compose run timdb -c=/local/"your_path" -a="algorithm"
````

3. To log data from database use -g flag, example:
``````
go run cmd/main.go -g -a="algorithm"
via docker
docker-compose run timdb -g=/local/"your_path" -a="algorithm"
``````

4. To get help with commands use -h flag, example:
```
go run cmd/main.go -h
via docker
docker-compose run timdb -h
```

## Documentation

To see documentation use:
````
godoc -http=:port(use port like 6060 or other)
````
And then visit "localhost:port"