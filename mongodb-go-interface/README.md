# MongoDB example

This example meant to test basic functionalities of mongoDB. The example allows to create,query and delete users and belonging user settings.
The documents look as follows:

- Documents in users collection
  ```
  users {
      _id -> objectID,
      name -> string,
      email -> string,
      password -> bcrypt blob
      _settings_id -> objectID
  }
- Documents in user_settings collection
  ```
  user_settings {
      _id -> objectID,
      2steps_on -> bool,
  }
# Usage
# Build

Run ```docker-compose up --build --force-recreate -d main-server``` to generate and start all containers.

In order to access the db run: ```docker exec -it user-data-db bash -c "mongo mongodb://root:123secure@user-db:27017"```

## Execution examples

Use the browser or curl command to execute the following:
- add new user that automatically generates a belonging settings document: ```http://localhost:8080/insert?name=testName&email=testEmail&password=testPass```
- get user with specified email: ```http://localhost:8080/get?email=testEmail```
- delete user specified by the email: ```http://localhost:8080/delete?email=testEmail```
- check user password and email: ```http://localhost:8080/check?email=testEmail&password=testPass```
- get user settings belonging to the user with specified email: ```http://localhost:8080/get-settings?email=testEmail```
- delete user settings belonging to the user with specified email: ```http://localhost:8080/delete-settings?email=testEmail```
