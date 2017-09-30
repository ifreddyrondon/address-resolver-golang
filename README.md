# Address resolver golang
Address API with geolocation. Converting addresses  (like "1600 Amphitheatre Parkway, Mountain View, CA") into geographic  coordinates, resolving then throw google maps geocoding and elevation api.

## Getting Started

It's necessary to set some environment variables to run the application.

### Database environment variables

The database environment variables use the following structure:

`postgres://username:password@localhost/db_name?sslmode=disable`

**DATABASE_URL** used for connect to a production database.

`export DATABASE_URL=postgres://localhost/addressresolver?sslmode=disable`

**DATABASE_URL_TEST** used for connect to the testing database. A running PostgreSQL server is required, with the ability to log in.

`export DATABASE_URL_TEST=postgres://localhost/addressresolver_test?sslmode=disable
`