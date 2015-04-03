# Licensing Server

This repository has the code that runs the Volcanic Pixels licensing system.
The system is responsible for periodically regenerating and signing the
license revocation lists that are stored in Google Storage.

It has a simple frontend that supports:

 - Creating a license
 - Revoking a license
 - Decoding a license


## License Architecture

A license is a JSON Web Token that is signed using RSA256, the private key is
kept secure on systems that need to issue licenses (this system and the main
Volcanic Pixels site where licenses are purchased).

The public key is distributed as part of the software that the license is for.

In addition to the claims defined in the JWT spec the license uses the `_attr`
claim to store an object of attributes, namely:

 - **email** - the email address of the customer
 - **name** - the name of the customer
 - **chargeId** - the charge ID relating to the license


## Revoking licenses

The licensing google storage bucket has a revocations.txt file that contains
one revocation per line in the following format (key id, and optional comment)

```
daS7y8sioiecYy # 2015-03-31 12:31 (key posted online)
```

Initially you might recoil at the thought of keeping this as just a simple text
file but it works very well and has some very nice properties, if any system
needs to be able to revoke licenses then all that needs to be done is to
download the google storage client library for that language, generate some
credentials that are valid for the revocations.txt file and write a line to a
file.

Equally, if a system wants to know when a license is revoked then using google
storage it can watch the file and get notified via webhook when it changes.
These things would be non-trivial and time consuming to setup another way.

The revocations.txt file is not intended for public consumption, instead the
server periodically (or when a change occurs) parses the file and generates a
revocations.json file (which is public). This file is signed with an expiry
date.

By default the revocations.json file is valid for 72 hours and is regenerated
every hour (to ensure that it always has ~48 hours of life left).  


## Circumventing the licensing

Since the products that this licensing system supports are written in
interpreted languages it will always be possible for a determined actor to
disable the licensing. The goal of the licensing protection is simply to make it
hard enough that someone would have to not value their time at all for it to
be economically beneficial.


## Transparent to the end user

I am a firm believer that any protection should not inhibit or even hinder the
ability of a paying customer to use the software in any way that is permitted.
Here are some guidelines that should be followed in new products:

### Err unlocked

If a check errors (e.g. remote server license check etc.) then the error
should result in the plugin remaining in its current state.

### License pre-loaded

If possible, the software should come with the license pre-loaded (i.e. bundled
with the code) such that upon first run the user doesn't have to do anything, it
"just works".
