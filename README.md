# GO INIT

Go init is a simple boilerplate for starting a new Go project. It includes a basic directory structure, and boilerplate code for the http server and openapi documentation with swagger ui.

## Changing the git remote URL
To change the git remote URL, run the following command:

```bash
git remote set-url origin <your_git_remote_url>
```

## Changing the project module name
To change the project module name, run the following command:

```bash
sh scripts/change-module-name.sh <your_module_name>
```

(optional)
```bash
sh scripts/change-module-name.sh <your_module_name> <old_module_name_to_replace>
```


## Build Project
To build the project, run the following command:

```bash
sh scripts/build.sh
```

## Run Project
To run the project, run the following command:

```bash
./bin/web
```

## Documentation
Visit Documentation at [http://localhost:8070/docs](http://localhost:8080/docs)

You can edit your OpenAPI documentation in the `docs/openapi.yaml` file.

### health check
You can check the health of the server by visiting [http://localhost:8080/api/v1/health](http://localhost:8080/api/v1/health)
