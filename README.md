# JWT AUTH FORM

# How it works.

  - Client when login will be returned JWT token which it will send with every request to authenticate.
  - When client first signup we create a JWT token and store in DB with all other things like email, passwords, user ID(created from MD5 hash of email and password).
  - When user logs in with email id and password we query DB to see if user exists and if it does we use UserID to create a token and after making it digitally signed using our secret token_password we return the whole login payload with the newly created token.
  - Now when user will hit the API it will send this token in the authorization header and we will check if its valid. Such routes are autmatically checked from jwt.go servie in services/ folder. Routes like /login /signup are discared as these not require JWT auth. We add the token in user context.Ex. Image Upload has JWT auth as its an action done once user signup/logs in.
   
   