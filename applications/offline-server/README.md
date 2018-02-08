# HTML page

Installation instruction of the offline-server.

# Prerequisities:

Administration rights are required to setup the web server.

1. Download and install GIT https://git-scm.com/downloads
2. Define the GOPATH environment variable.

![environment variable](img/env_variables.png)

![new environment variable](img/new_env_variable.png)

3. Open new cmd window.

# Setup:

1. mkdir -p %GOPATH%\src\github.com\WPTechInnovation\
2. cd %GOPATH%\src\github.com\WPTechInnovation\
3. git clone https://github.com/WPTechInnovation/wpw-sdk-go.git
4. cd wpw-sdk-go
5. git checkout develop
6. cd applications\offline-server
7. go get
8. go build

Additional Info:
- To run this server type: 
```
cd %GOPATH%\src\github.com\WPTechInnovation\wpw-sdk-go\applications\offline-server
offline-server.exe
```