# Movie-Api
In this application, the member with admin authority will be able to create Movie, adding movies and delete movies. Customers will be able to see existing movies. Using these we aim to develop a movie application.

## Requests of application

1. You should get the required information from the user and create a user in the database and return JWT token in response.(Sign-up)
2. Users registered in the database must log in to the system with their email and password, if both information is correct, you must create a JWT token and return to the user.(Login)
3. All active and not deleted movies in the database should be listed.(List Movie)
4. Users in the admin role can delete products.(Delete Movie)
5. Users in the admin role can update the product.(Update Movie)

## TO-DO

- add more test
- with using cobra in configuration management separating services based on domain 
- communication with client between services
- add docker
- configure docker compose

## Used Files

* main.go: Main application file.
* Config: Configuration files.
* Repositories: CRUD operations handling.
* Models: It is our domain data.
* Handlers: This files , they receive the request from the user, they ask the services to perform an action for them on the database.
* Services: contains some business logic for each model, and for authorization.
* Middlewares: it contains middlewares(golang functions) that are triggered before the controller action, for example, a middleware which reads the request looking for the Jwt token and trying to authenticate the user before forwarding the request to the corresponding controller action.
* PKG: it contains the packages that are used by the application.
## Requirements

* Go Language
* Git
* Go Module
* GORM
* Gin
* Postgres
* JWT
* Viper
* Swagger
* Zerolog

## Acknowledgments

* Hat tip to anyone whose code was used
* Inspiration
* etc

## License
[MIT](https://choosealicense.com/licenses/mit/)
