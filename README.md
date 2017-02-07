# library: an interview coding challenge
This challenge was given to me by an unnamed company for an interview question.
I am leaving the name out of this to allow them to reuse the question and for me
to showcase it for other companies if this one didn't work out :)

I have written CLI's in Click and web servers using Flask and Django.
I knew I could do the challenge faster using those tools, however I like to use Go any chance I get.
I absolutely love the language and think its very well suited for web servers and CLI's.

Library was not developed using a framework like Gin. I Prefer to use the standard libs and not go for
external dependencies because often these lead to unnecessary abstractions. The standard libs are so rich
anyways, why add complexity?

The API is your standard JSON RESTful API. It's unauthenticated, but that's OK for the sake of this short demo.
Since I have added routes for the user, it would be a very easy to add a tokenized auth scheme like JSON Web Tokens.

## Improvements
Things I would have wanted to work on for this coding challenge:
- Writing a client side api wrapper that would provide a clean interface for the CLI application.
- Comments! Lots of comments, package level explanations, db structure, etc.
- Adding full integration tests
- Dockerize the build
- Dockerize the tests
- Dockerize the deployment
- Add a makefile to: build, test, bring up and down
- Validation on the models for input and output
- Standardize the response so the api always returns valid json
- Add authentication
- Add a get on the undo handler to allow they user to check what will be undone

## server
The server package contains all handler functions and the database api layer

###  /api/${user}/book?by=string&read=bool
POST: add a book into your library
GET: gets a list of all of your books

The ${user} in the url structure defines which user to store and retrieve user information.
Query params are used to provide context to the GET request. The `by` param is used to add a
search query to select only the book by the given author name. The `read` param is used to filter
books by whether they have been read or not.

The post method takes a json body format as follows:
```
{
    "read": bool, # optional. Defaults to false
    "author": string,
    "title": string
}
```

### /api/${user}/book/${name}
POST: sets a value in the database indicating the book has been read
GET: gets all information about the book given identified by the name

The post method takes a json body format as follows:
```
{
    "read": bool
}
```

### /api/${user}/undo
POST: server caches the last operation for the user and this endpoint will undo the last operation


## Persistent Storage: boltdb
Instead of going with a more traditional route like postgres or sqlite3 I decided to use
a key value store called boltdb. Bolt is a fantastic datastore written in pure Go using
a B+ tree under the hood. It's a great storage choice for this application because the
question is essentially asking to make a key value store anyways. Generally boltdb is great because
its simple, written in Go with a clean API, embedded into the built binary, and its fast.

The organization of the key/val store is `/user/title: Book`. Meaning the user name is
the high level bucket which contains keys of the normalized title and a value of the
Book struct marshaled into JSON. This makes writing CRUD style endpoints very easy.

In order to view the data I used a tool called boltbrowser which is a terminal based
boltdb browser. It's great for introspecting into bucket structure and checking to see
if your data is getting updated.


# Library: command line interface
OK, when I mentioned above that I am using standard libs only, I lied.
The cli of library is built using the awesome cobra library. Docker,
Kubernetes, etdc, CockroachDB, and many more successful cli tools use Cobra as its cli
framework. It takes all of the headache around handling input and running commands
based on user input.

There's a root command that collects global flags and adds subcommands to it which are
defined in package level structs. Very similar to python's Click cli framework.
These structures define functions to run at the various execution points in Cobra.

I intentionally copy & pasted a lot of code that could have been abstracted into an
api wrapper in the cli package. This was to save time on the coding so I could focus
more on the writeup.
