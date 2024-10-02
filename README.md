# NotesAPI

This project is a simple API for creating, reading, updating, and deleting notes. The API is written in Go and uses the [Gorilla Mux](http://www.gorillatoolkit.org/pkg/mux) framework for routing.

## Description

The API allows you to perform the following operations with notes:

- Get a list of all notes
- Get a note by ID
- Create a new note
- Update an existing note by ID
- Delete a note by ID

The API will be available at http://localhost:8080.

## Usage
Here are some examples of using the API with curl:

Get all notes:
curl -X GET http://localhost:8080/notes

Get note by ID:
curl -X GET http://localhost:8080/notes/test

Create a new note:
curl -X POST http://localhost:8080/notes -H "Content-Type: application/json" -d '{"id":"test","content":"This is a new test note"}'

Update note:
curl -X PUT http://localhost:8080/notes/test -H "Content-Type: application/json" -d '{"id":"test","content":"This is an updated test note"}'

Delete note:
curl -X DELETE http://localhost:8080/notes/test
