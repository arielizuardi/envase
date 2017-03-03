# envase
Envase means container, is a docker-invoker testing-helper for Go language

## Why ?

Sometimes it is really hard to mock and test your code against the implementation of several engine like RDBMS (mysql, oracle, sql-server), fluentd, etc.
So with the help of docker engine, we can test the real implementation of our code against the engine. 
For example: In MySQL if you want to make sure your query is working like expected.

Hopefully, this library will make testing easier.

## Usage 

```go get github.com/arielizuardi/envase```

If you are using glide

```glide get github.com/arielizuardi/envase```

## Example

## Start `fluentd` Container

```
ctx := context.Background()
dockerClient, err := client.NewEnvClient()
if err != nil {
    panic(err)
}

fluentdImage := `fluent/fluentd:v0.12.32`
fluentHost := `127.0.0.1`
containerPort := `24224`
exposedPort := `24224`
containerName := `fluent_test`

envConfig := []string{}

envConfig := []string{}
container := envase.NewDockerContainer(ctx, dockerClient, fluentdImage , fluentHost, containerPort, exposedPort, containerName, envConfig)

err := container.Start()

if err != nil {
    panic(err)
}
  
```

## Start `mysql` Container

```
dockerClient, err := client.NewEnvClient()
if err != nil {
    panic(err)
}
   
envConfig := []string{
	"MYSQL_USER="your_user",
	"MYSQL_ROOT_PASSWORD="your_password",
	"MYSQL_DATABASE="your_database",
}

mysqlImage := `mysql:5.7`
mysqlHost := `127.0.0.1`
containerPort := `3306`
exposedPort := `33060` 
containerName := `mysql_test`
  
container := envase.NewDockerContainer(ctx, dockerClient, mysqlImage, mysqlHost, containerPort, exposedPort, containerName, envConfig)
  
if err != nil {
   panic(err)
}
```

### Stop Your Container
```
err := container.Stop()
if err != nil {
    panic(err)
}
```
