go blog aggregator to recieve RSS feeds

Depencencies: Go and Postgres

blog agggregator that you can have a user and follow rss feeds, then aggregate them to make post and add them to the database.

go install from root of the project to make the gator command work

create a ~/.gatorconfig.json file that has a struct like this
  {
  "db_url": "connection_string_goes_here",
  "current_user_name": "username_goes_here"
  }

start a posgressql server from your terminal to have a database to connect to
The url might look something like "postgres://wagslane:@localhost:5432/gator"

from cli run "gator" + (command) to do anything in the app

commands:
  login- given a user arguement, will log in a user that is in the database
  register- given a user arguemnt, will register a new user to the app
  users- show a list of all the registered users
  addfeed- given a name and url arguement, adds a feed to a users profile and follows it
  feeds- shows a list of feeds that are in the database
  following- shows list of feeds current user is following
  reset- resets the datatbase

