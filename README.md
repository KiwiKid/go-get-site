```bash
# (gvm use go1.21.1)
# Assuming you are executing the script from the root directory of the project
templ generate
mkdir -p app
go build -o app/go-get-site && ./app/go-get-site



templ generate --watch
go build -o app/go-get-site && ./app/go-get-site


# Migrate (only create) new DB changes
./app/go-get-site migrate

```


NB: un-comment the Migrate() call to apply db migraitons

TODO:

- Add air for development hot-reloading



