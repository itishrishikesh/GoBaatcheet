# GoBaatcheet
It's simple websocket based chat application, I decided to build this as a learning exercise. It makes use of couple of external components like authorizer (which is an open source authentication server build using go) and kafka (which I'm using to store messages for offline user).
The below image kind of show the flow of the application and interactions with different components.
![gobaatcheet.png](docs%2Fimages%2Fgobaatcheet.png)

### State of the project.

- The project is still work in progress, although the backend is functional.

### Testing/Installing/Running the project
_[This is expected to fail for a while]_
- Make use of the docker compose file.

```bash
docker compose up -d
```

- This will bring up kafka, authorizer, build backend and frontend.
- You should be able to access the frontend, firstly you'll need to register yourself with authorizer server.
- Once you register, you can now specify username of other registerd users.