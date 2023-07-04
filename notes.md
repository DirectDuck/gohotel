## Notes
1. bson (binary json) for mongo, json for fiber
2. `omitempty` useful for id, `json:"_"` useful for passwords
3. myVar.(type) -> cast (something like interface?) to type
4. bson.M is simpler, but doesn't preserve order. bson.D wher order matters
5. Use github link to repository in mod init like `go mod init github.com/DirectDuck/gohotel`
6. context.TODO as placeholder to update in the future, Background if you just need one

## Todo
1. Move urls to api file
2. Add api docs swagger
3. Make logging (with sentry maybe?) 
4. Put user in context in jwt success callback for easier access
5. where to put `Params`?
6. Write readme & .env.examples
