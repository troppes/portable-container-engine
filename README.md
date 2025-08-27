# Portable Container Engine

This project want to create a small executable that makes it possible to download and explore the file system of docker images. For this it has two main modes of operation: run and download. Run downloads an image, extracts the fs and runs it. The download command just downloads the full image and saves it as a tar file.

## How to Use

```
go run cmd/pce/main.go <run|download> <image> <command> [-x]
```

Example:

Download and Extract Alpine
```
go run cmd/pce/main.go download alpine:latest -x
```

## Not yet supported:

- not all namespaces
- cGroups
- does not run the docker intened command

## Special Thanks

Liz Rice for the idea and base for this project: [https://github.com/lizrice/containers-from-scratch](https://github.com/lizrice/containers-from-scratch)
