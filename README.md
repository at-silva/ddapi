# DDAPI (Database Direct API)

DDAPI is a set of Go http handlers designed to allow SQL (or pretty much any other query language) to be safely executed from the frontend.

The idea is to be able to safely embed database queries ([DQL](https://en.wikipedia.org/wiki/Data_query_language)) and statements ([DML](https://en.wikipedia.org/wiki/Data_definition_language)) in frontend code.

## But, why?

The goal is to help you build truly flexible backends in a fast and frictionless way. 

The project was designed with the [backend for frontend](https://docs.microsoft.com/en-us/azure/architecture/patterns/backends-for-frontends) pattern in mind and it's not intended to replace your backend but to rather supplement it.

## How's that safe? Wouldn't this allow an ill-intentioned individual to just `cURL` SQL into my API?

No! only queries signed with your secret key are allowed to run, and the key is never deployed with your frontend code. We are using the same technology people use secure [JWTs](https://jwt.io).

## But what about the query parameters? Wouldn't it be possible to manipulate the request to read or write data outside of a user's authorization scope?

No! Well... Yes, but that's a problem we already have to deal with when we're working with any other equivalent solutions like GraphQL for instance. 

With that said, we do have two mechanisms in place to help us mitigate the risk though:

 - **JWT claims injection**: parameters like user id, or tenant id can be read straight from a signed JWT, that makes impersonating another user pretty much impossible.
 - **Server-side parameters validation**: we're leveraging a JSON schema validation engine to declaratively to restrict the possible values any single input parameter can contain.

## And what about caching, rate limiting, instrumentation and all that stuff? 

You can keep using MemCached, redis, NewRelic and any other tool you already use today, in the end DDAPI is just a set of http handlers and validators that can be wrapped and extended with pretty much anything you want.

## Is it cumbersome to work with?

Not really, no. Let me show you:

![Screencast 1]( https://github.com/at-silva/ddapi/raw/main/docs/screencast1.gif "Screencast")

Pretty cool huh?

## It is pretty cool... Is it production ready yet?

No, not yet, I'm just trying to put the idea out there and gather some feedback from the community.

## Does this work with (insert your preferred tech stack here)?

I believe DDAPI is more concept than tool at this point, the basic idea - to have signed DML and DQL embedded into the frontend - can be implemented in any stach, this implementation is specifically designed to be plugged in Go backends and Javascript (React, Vue, Angular, etc...) frontends though.

## This is crazy man, you're crazy.

Awww, really?! You're so sweet, thank you!

## Contributing
All pull requests and discussions are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](https://choosealicense.com/licenses/mit/)