## go-fiber-auth-2024

### Deploy

```shell script
git clone https://github.com/peterdee/go-fiber-auth-2024
cd ./go-fiber-auth-2024
gvm use go1.22
go mod download
```

### Environment variables

`ENV_SOURCE` environment variable determines the origin of environment variables: `file` or `env`

If `ENV_SOURCE` variable is set to `file`, then `.env` file is required and the app will not launch without it

Required environment variables are located in [.env.example](./.env.example) file

### Launch

**Without Docker**

```shell script
ENV_SOURCE=file go run ./

# With AIR
ENV_SOURCE=file air
```

**With Docker**

```shell script
docker run -p 2024:2024 --env-file .env -it $(docker build -q .)
```

### License

[MIT](./LICENSE.md)
