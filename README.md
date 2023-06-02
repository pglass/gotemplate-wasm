Overview
--------

Quick demo of compiling Go to WASM, to support rendering Go templates from JavaScript.

Build and run (requires Python3 for the HTTP server):

```
make build
make run
```

* This outputs `main.wasm`
* This copies `wasm_exec.js` from GOROOT into the current directory
* It starts a web server

Open http://localhost:8080. Open the dev console (cmd-i). Click 'Console', and enter the following
JavaScript: 

```js
> renderGoTemplate("name={{ .name }}", JSON.stringify({name: "wumbo"}))
```

You should see:

```
2023/06/02 17:03:35 renderGoTemplate: args=[name={{ .name }} {"name":"wumbo"}]
'name=wumbo'
```

Notes:

* I used `JSON.stringify` to pass the object to the Go function, because there doesn't seem to be a
  way to iterate the object keys in a `js.Value`. Instead, I parse a JSON string.
* `syscall/js` is experimental.
* `main.wasm` is 4.5 MB.
