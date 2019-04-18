# Windows Docker Registry Service
Registry service will install / un-install, start / stop, and run as a service (daemon).

Build the registry:

```pwsh
PS C:\windows-docker-registry-service> go build -o registry.exe main.go 
```

Install and start the service:

```pwsh
PS C:\windows-docker-registry-service> ./registry -service install
PS C:\windows-docker-registry-service> ./registry -service start 
```

Stop and uninstall the service:

```pwsh
PS C:\windows-docker-registry-service> ./registry -service stop
PS C:\windows-docker-registry-service> ./registry -service uninstall 
```
