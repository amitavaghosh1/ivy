#### Requirements
- go
- nodejs
- yarn

#### Building
- make build (build executable for all environment)
    - make darwin (for mac)
    - make linux (for linux)
    - make windows (for windows)

The respective builds can be found in
- out/windows
- out/linux
- out/darwin

#### Running

##### Running the cli app
- ./out/${GOOS}/ivy --help

##### Running the web app
- make run.server

For development of the UI:
- cd web/frontend
- yarn start
- in web/frontend/src/index.js, uncomment `window.BaseURL = "http://localhost:4000"`


### Web App In Action:



https://user-images.githubusercontent.com/71380768/149381958-084f0d8a-eebd-42ff-92d4-920731bef536.mov


