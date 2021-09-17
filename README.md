# How to rebuild

# Setup Firebase RTDB

Follow the [firebase documentation](https://firebase.google.com/docs/database/admin/start#go) to set up your database.

> You will need to update key.json with the service account file that you get from the Firebase console.
> 
> Update the DatabaseURL inside firebase.go accordingly.

Why Firebase?

Firebase RTDB and the SDK allows you to listen to database changes out of the box. It's great for a quick project, an alternative to the pub/sub service.

# Alternative

Remove all Firebase related code and replace it with your own DaaS and the rest of the project should still function well.