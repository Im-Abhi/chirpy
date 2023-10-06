# CHIRPY
A REST api made in golang.

### Setting Up in your local Environment
1. Make sure you have go installed. If not [go](https://go.dev/dl/)
2. After installation check your go version on command line using this command
```
go version
```
This command should return something like 
```
go version go1.21.1 windows/amd64
```

3. Copy the [.evn.example](./.env.example) file, rename the file to <b>.env</b> and add the required environment varibles. After you are done the file content should look like this
```env
PORT=8000
JWT_SECRET="chirpy_secret"
POLKA_KEY="f271c81ff7084ee5b99a5091b42d126e"
```
4. Once steps 1-3 are done you are good to go.
Run this command to spin up the chirpy server.
```
go build
.\chirpy.exe
```
5. If the build fails the executable won't be formed and you should be able to see the error in the terminal.
6. If the build succeeds you should see something like this.
```
2023/10/06 20:47:40 Serving files from . on port: 8000
```

### Want to Contribute to the project
You can do so by creating a fork of this repo and creating a pull request for the changes.\
Looking for someone who can help in the API documentation (someone new to github and wants to contribute)