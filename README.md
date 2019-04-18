# Windows Docker Registry Service
Registry service will install / un-install, start / stop, and run a program as a service (daemon).

Step 1) build the registry:

```pwsh
PS C:\windows-docker-registry-service> go build -o registry.exe main.go 
```

Step 2) install and start the service:

```pwsh
PS C:\windows-docker-registry-service> ./registry -service install
PS C:\windows-docker-registry-service> ./registry -service start 
```
