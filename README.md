# interview coding challenge: library
This challege was given to me by an unnamed company for an interview question.
I am leaving the name out of this to allow them to reuse the question and for me
to showcase it for other companies if this one didn't work out :)

# API & CLI
The api is a json rest server that has the folowing endpoints:

##  /api/${user}/book?by=string&read=bool
POST: add a book into your library
GET: gets a list of all of your books

The ${user} in the url structure defines which user to store and retrive user information.
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

## /api/${user}/book/${name}
POST: sets a value in the database indicating the book has been read
GET: gets all information about the book given identified by the name

The post method takes a json body format as follows:
```
{
    "read": bool
}
```

## /api/${user}/undo
POST: server caches the last operation for the user and this endpoint will undo the last operation


# Persistent Storage
Instead of going with a more traditional route like postgres or sqlite3 I decided to use
a key value store called boltdb. Bolt is a fantastic datastore written in pure Go using
a B+ tree under the hood. Its a great storage choice for this application becuase the
its simple, written in Go with a clean API, embedded and bundled into the built binary,
and its fast.

The structure of the key/val store is `/user/title: Book`. Meaning the user name is
the high level bucket which contains keys of the normalized title and a value of the
Book struct marshaled into JSON. This makes writing CRUD style endpoints very easy.

In order to view the data I used a tool called boltbrowser which is a terminal based
boltdb browser. It's grat for introspecting into bucket structure and checking to see
if your data is getting updated.
