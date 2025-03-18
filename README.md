# jserve

jserve is a handy cli tool for front-end developers to mock up a back-end in a snap, making it easier for devs to test and prototype without any hassle.

## routes

| Route                    | Description                           |
| ------------------------ | ------------------------------------- |
| **GET** /                | List database tables.                 |
| **GET** /{table}         | List a certain table records.         |
| **POST** /{table}        | Add a record to a given table.        |
| **PUT** /{table}/{id}    | Edit a record by id.                  |
| **DELETE** /{table}/{id} | Delete a record from a table.         |
| **GET** /{table}/{id}    | Retrieve a record from a table by id. |

## Filtering and sorting

| Operation | Route                              | Description                                                                                  | Example                            |
| --------- | ---------------------------------- | -------------------------------------------------------------------------------------------- | ---------------------------------- |
| Filtering | **GET** /{table}?col=*{}*&val=*{}* | Search through a given table, based on a column value                                        | GET /posts?col=author&val=Jane+Doe |
| Sorting   | **GET** /{table}?sort=*{}*         | sort a table based on a column, for descending sorting append the column name with a "**-**" | GET /posts?sort=-id                |

## CLI arguments

| Argument | Default   | Description                               |
| -------- | --------- | ----------------------------------------- |
| **db**   | ./db.json | the path to the json databse.             |
| **port** | :4000     | the port the server will be listening to. |

## Local usage

jserve can be downloaded as a standalone binary directly from the [releases page](https://github.com/aymenhta/jserve/releases).

```bash
./jserve --db=./db.json
```

**Alternative**: If you do not prefer that method you can simply install it using go:

```bash
go install github.com/aymenhta/jserve@latest
```
