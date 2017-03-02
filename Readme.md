This program listens on port 3000. 

example usage: 
`curl --data "hey @bonnie, good work!" http://localhost:3000`

Known Issues:
One of the tests doesn't pass - The example text did not interpolate quotes, `&quot;` and for some reason they are interpolating as `"`

The json is not marshalling correctly - the `omitempty` param on the Msg struct (see `msg.go`) should be well, omitting empty struct values. see https://golang.org/pkg/encoding/json/
>The "omitempty" option specifies that the field should be omitted from the encoding if the field has an empty value, defined as false, 0, a nil pointer, a nil interface value, and any empty array, slice, map, or string.
 
 It's not omitting empty struct values. 
 
lines 79-83 on main_test.go has some commented out lines that dumps the response if you want to see for yourself (make sure to uncomment the relevant import `net/http/httputil` as well)  - the response looks like
```javascript
 {
 	"mentions": ["chris"],
 	"emoticons": null,
 	"links": []
 }
```
rather than the given example:

```javascript 
{
   "mentions": [
     "chris"
   ]
 }
```

I have some ideas as to why this might be happening, but time prevented me from exploring them further.

If I had more time, I would include:
- unit tests. This was quick and dirty so I just did  black box style integration tests
- make the integration tests a bit more explicit, ie, not just checking the entire response object, but checking each value in the response object. 
- Logging. There are some errors that should be logged properly. right now its just sending errors to stout
- handle  client facing  (http) errors more explicitly than sending out an internal server error
- client input validation (to help prevent internal errors)


If this were an actual service, I would: 
- add security. it just leaves a port wide open. How I would secure it depends entirely on  it's planned usage
- Dockerize it 
- ask a lot more questions about the purpose of this service and road map out more features

outside libraries used: 
- `title.go` is adapted from [this gist](https://siongui.github.io/2016/05/10/go-get-html-title-via-net-html/)
- [github.com/mvdan/xurls](https://github.com/mvdan/xurls) for capturing URLs from messages
- [github.com/gorilla/mux](https://github.com/gorilla/mux) as my http router

