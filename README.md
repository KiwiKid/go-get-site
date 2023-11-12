```bash
# (gvm use go1.21.1)
# Assuming you are executing the script from the root directory of the project
templ generate
mkdir -p app
go build -o app/go-get-site && ./app/go-get-site


# Build and log to file
templ generate && go build -o app/go-get-site && ./app/go-get-site &> logs4.txt

```


NB: un-comment the Migrate() call to apply db migraitons

TODO:

- Add air for development hot-reloading



