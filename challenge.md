We'd like you to write a small system for managing a personal library. The system should be accessible from the command line. A user would interact with it like so:

  $ ./library

	Welcome to your library!

	> add "The Grapes of Wrath" "John Steinbeck"

	Added "The Grapes of Wrath" by John Steinbeck

	> add "Of Mice and Men" "John Steinbeck"

	Added "Of Mice and Men" by John Steinbeck

	> add "Moby Dick" "Herman Melville"

	Added "Moby Dick" by Herman Melville

	> show all

	"The Grapes of Wrath" by John Steinbeck (unread)
	"Of Mice and Men" by John Steinbeck (unread)
	"Moby Dick" by Herman Melville (unread)

	> read "Moby Dick"

	"Moby Dick" by Herman Melville marked as read

	> undo

	"Moby Dick" by Herman Melville marked as unread

	> read "Of Mice and Men"

	"Of Mice and Men" by John Steinbeck marked as read

	> add "Infinite Jest" "David Foster Wallace"

	Added "Infinite Jest" by David Foster Wallace

	> undo

	Removed "Infinite Jest" by David Foster Wallace

	> show all

	"The Grapes of Wrath" by John Steinbeck (unread)
	"Of Mice and Men" by John Steinbeck (read)
	"Moby Dick" by Herman Melville (unread)

	> show unread

	"The Grapes of Wrath" by John Steinbeck (unread)
	"Moby Dick" by Herman Melville (unread)

	> show all by "John Steinbeck"

	"The Grapes of Wrath" by John Steinbeck (unread)
	"Of Mice and Men" by John Steinbeck (read)

	> show unread by "John Steinbeck"

	"The Grapes of Wrath" by John Steinbeck (unread)

	> quit

	Bye!

	$

--------------------------

As shown above, the program should accept the following commands:

- **add "$title" "$author"**: adds a book to the library with the given title and author. All books are unread by default.
- **read "$title"**: marks a given book as read.
- **show all**: displays all of the books in the library
- **show unread**: display all of the books that are unread
- **show all by "$author"**: shows all of the books in the library by the given author.
- **show unread by "$author"**: shows the unread books in the library by the given author
- **undo**: undoes the last mutational command (if a book was marked as read it marks it as unread; if a book was added, it gets removed)
- **quit**: quits the program.

Some other stipulations:

- You can use whatever language you want.
- There's no need to use a persistance mechanism (ie, a SQL database) for the books. You can just store them in memory. That is, every time you run the program, the list of books should be empty.
0
