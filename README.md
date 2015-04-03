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

### Fail unlocked

If a check fails (e.g. remote server license check etc.) then the failure
should result in the plugin being unlocked.

### License pre-loaded

If possible, the software should come with the license pre-loaded (i.e. bundled
with the code) such that upon first run the user doesn't have to do anything, it
"just works".
