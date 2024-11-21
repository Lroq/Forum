# Forum Project

This project consists in creating a web forum that allows:

- Communication between users.
- Associating categories to posts.
- Liking and disliking posts and comments.
- Filtering posts.

## SQLite

To store the data in your forum (like users, posts, comments, etc.), you will use the database library SQLite.

SQLite is a popular choice as an embedded database software for local/client storage in application software such as web browsers. It enables you to create a database as well as control it by using queries.

To structure your database and achieve better performance, we highly advise you to take a look at the entity relationship diagram and build one based on your own database.

You must use at least one SELECT, one CREATE, and one INSERT query. To know more about SQLite, you can check the [SQLite page](https://www.sqlite.org/).

## Authentication

In this segment, the client must be able to register as a new user on the forum by inputting their credentials. You also have to create a login session to access the forum and be able to add posts and comments.

You should use cookies to allow each user to have only one open session. Each of these sessions must contain an expiration date. It is up to you to decide how long the cookie stays "alive." The use of UUID is a Bonus task.

### Instructions for User Registration

- Must ask for email
  - When the email is already taken, return an error response.
- Must ask for username
- Must ask for password
  - The password must be encrypted when stored (this is a Bonus task)

The forum must be able to check if the email provided is present in the database and if all credentials are correct. It will check if the password is the same as the one provided, and if the password is not the same, it will return an error response.

## Communication

To enable communication between users, they will be able to create posts and comments.

- Only registered users will be able to create posts and comments.
- When registered users create a post, they can associate one or more categories with it.
  - The implementation and choice of the categories is up to you.
- The posts and comments should be visible to all users (registered or not).
  - Non-registered users will only be able to see posts and comments.

## Likes and Dislikes

Only registered users will be able to like or dislike posts and comments.

- The number of likes and dislikes should be visible to all users (registered or not).

## Filter

You need to implement a filter mechanism that will allow users to filter the displayed posts by:

- Categories
- Created posts
- Liked posts

You can look at filtering by categories as subforums. A subforum is a section of an online forum dedicated to a specific topic.

Note that the last two are only available for registered users and must refer to the logged-in user.

## Docker

For the forum project, you must use Docker. You can read about docker basics in the ascii-art-web-dockerize subject.

## Authentication with External Providers

The goal of this project is to implement new ways of authentication into your forum. You must be able to register and log in using at least Google and GitHub authentication tools.

### Examples of Authentication Providers

- Facebook
- GitHub
- Google

### Instructions

- Your project must have implemented at least the two authentication examples given.
- Your project must be written in Go.
- The code must respect good practices.

You must follow the same principles as the first subject.

For this project, you must take into account the security of your forum.

### HTTPS

You should implement a Hypertext Transfer Protocol Secure (HTTPS) protocol:

- Encrypted connection: For this, you will have to generate an SSL certificate. You can think of this as an identity card for your website. You can create your certificates or use "Certificate Authorities" (CA's).

We recommend you take a look into cipher suites.

### Rate Limiting

The implementation of Rate Limiting must be present in this project.

### Encryption

You should encrypt at least the clients' passwords. As a Bonus, you can also encrypt the database. For this, you will have to create a password for your database.

### Sessions and Cookies

Clients' session cookies should be unique. For instance, the session state is stored on the server, and the session should present a unique identifier. This way, the client has no direct access to it. Therefore, there is no way for attackers to read or tamper with the session state.

You must follow the same principles as the first subject.

## Image Upload

In the forum, registered users have the possibility to create a post containing an image as well as text.

- When viewing the post, users and guests should see the image associated with it.
- There are several extensions for images like: JPEG, SVG, PNG, GIF, etc. In this project, you have to handle at least JPEG, PNG, and GIF types.
- The max size of the images to load should be 20 MB. If there is an attempt to load an image greater than 20 MB, an error message should inform the user that the image is too big.

You must follow the same principles as the first subject.

## Forum Moderation

The forum moderation will be based on a moderation system. It must present a moderator who, depending on the access level of a user or the forum setup, approves posted messages before they become publicly visible.

The filtering can be done depending on the categories of the post being sorted by irrelevant, obscene, illegal, or insulting.

### User Types

You should take into account all types of users that can exist in a forum and their levels. You should implement at least four types of users:

- **Guests**: Unregistered users who can neither post, comment, like, or dislike a post. They only have permission to see those posts, comments, likes, or dislikes.
- **Users**: Users who can create, comment, like, or dislike posts.
- **Moderators**: Users who have granted access to special functions:
  - Monitor content in the forum by deleting or reporting posts to the admin.
  - To create a moderator, the user should request an admin for that role.
- **Administrators**: Users who manage the technical details required for running the forum. This user must be able to:
  - Promote or demote a normal user to or from a moderator user.
  - Receive reports from moderators. If the admin receives a report from a moderator, he can respond to that report.
  - Delete posts and comments.
  - Manage the categories by being able to create and delete them.

You must follow the same principles as the first subject.

## Advanced Features

You will have to implement the following advanced features:

- Notify users when their posts are:
  - Liked/disliked
  - Commented
- Create an activity page that tracks the user's activity. This page should:
  - Show the user's created posts.
  - Show where the user left a like or a dislike.
  - Show where and what the user has been commenting. The comment must be shown, as well as the post commented on.
- Create a section where users will be able to edit/remove posts and comments.

We encourage you to add any other additional features that you find relevant.

## Commands

To set up and install the necessary dependencies for this project, run the following commands:

```sh
go mod init Forum
go get github.com/mattn/go-sqlite3
go get golang.org/x/crypto/bcrypt
go get golang.org/x/oauth2/google@v0.21.0
go get golang.org/x/oauth2
go get github.com/gorilla/sessions
go get github.com/gorilla/websocket
go get -u github.com/gorilla/mux
go get github.com/google/uuid
go mod tidy
