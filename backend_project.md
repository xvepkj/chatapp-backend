### Backend Project

**<u>Requirements : </u>**

- My backend has a database - **OK**

- My backend has an API the frontend uses **OK**

- I can also use postman, curl or wget to call my API **OK**

- My API uses POST/PUT calls to update the database and GET calls to return data **OK**

- I know SQL well enough to create my database, create the tables, insert the 
  
  records, update or delete the records, add indexes and manage the data **OK**

- I know what an ORM is and how to use at least one, with the pros and cons compared to using pain SQL **OK**

- I understand authentication and security to some degree, including white hat hacking and things like sql injection, cross-site man in the middle attacks, phishing etc. **OK**

- My source code is stored in version control (git) **OK**

- you understand logging and alerts **OK**

- I can output different file types from my database (excel, word, html emails, text) **OK**

- I understand how to write and execute unit tests for my code

- Authorization & roles

- Can integrate with external identity providers (Google Oauth)
  
  **You get extra bonus points if:**

- you are good at technical documentation (can use swagger) **OK**

- you understand kubernetes or similar tools

- I have a basic understanding of containerization ( most likely docker)

- I have a basic understanding of CICD pipelines

- you have some cloud deployment knowledge

### Project Idea : Multilanguage chat app (ChatVoyage)

- People can register using 
  
  - Username
  
  - Password 
  
  - Confirm Password 

- People can login using 
  
  - Username 
  
  - Password

- Can view their existing chats and send messages 

- Can send messages to new users by searching using usernames

- They can have a preferred language and any chats from the other side will be translated in their preffered language 

#### Technical Details

- Install go

- `go mod init github.com/yourusername/chatapp-backend`

- `go get -u github.com/gin-gonic/gin`

- Install postgres

- `psql -U postgres`

- `CREATE DATABASE chatapp;`

- `CREATE USER chatuser WITH PASSWORD 'yoyoyoyoyo';`

- `GRANT ALL PRIVILEGES ON DATABASE chatapp to chatuser;`

- `go get -u gorm.io/gorm`

- `go get -u gorm.io/driver/postgres`

- Write db functions for `user` model

- Write wrappers for db functions in handlers 

- Write routes (currently in `main` go file only)

- `swag init` to generate swagger docs 



### Design


