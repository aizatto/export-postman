# Generate documentation

```sh
go run main.go file.postman_collection.json
```

## Automatically Run

If using [gomon](http://github.com/aizatto/gomon)

```sh
~/go/bin/gomon "**.go" -- go run main.go file.json
```

# Convert to other file format (using Pandoc)

https://pandoc.org/getting-started.html

```
pandoc file.md -f markdown -t docx -s -o api.docx
pandoc file.md -f markdown -t html -s -o api.html
```