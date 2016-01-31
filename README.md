# go-workshop

This is an educational purpose application.

Application is splitted in three phases.

You are encouraged to use `go test` and resolve all test errors.
To mange development in phases use `go test -phase X` where `X` is phase number.

## Phases

1. Develop CLI application that reads list of CSS files from JSON file and merge them in single CSS file. Use example:

```bash
app.exe -list my_list.js -out merged.css
```

2. Upgrade application with `watch mode` that will watch for changes in list of provided CSS files and rebuild merged CSS file continuously. Use example:

```bash
app.exe -watch -list my_list.js -out merged.css
```

3. Upgrade application so it can serve merged CSS file via HTTP protocol. Use example:

```bash
app.exe -watch -serve 8080 -list my_list.js -out merged.css
```
