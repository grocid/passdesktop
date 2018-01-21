# Pass for macOS

Pass is a GUI for [pass](https://github.com/grocid/pass), but completely independent of it. It communicates with [Vault](https://www.vaultproject.io), where all accounts with associated usernames and passwords are stored.

Pass is password protected on the local computer, by storing an `encrypted token` on disk, along with a `nonce` and `salt`. The `token` is encrypted with AES-GCM-256. The encryption key is derived as 

```key := PBKDF2(password, salt)``` 

and 

```token := AES-GCM-Decrypt(encrypted token, key, nonce).```

`token` is kept in memory only.

The program is fairly short and easily auditable, if anyone feels encouraged to do it.

## Not implemented...

 - Generation of master password (locally)
 - Generation of new accounts
 - Changing account details remotely (username and password)
 - Remote deletion of accounts
 - Mutual authentication as an option for communicating with Vault

These, quite essential, features are of course already supported in [pass](https://github.com/grocid/pass). 

## Seurity considerations

 - To use mutual authentication. This requires each accessor/user to have a valid private key. Private key can be encrypted with master password.
 - Using [fw](https://github.com/grocid/fw) to only allow white-listed users. Requires the user to authenticate with Google Authenticator to white list its IP address. Makes it harder for attackers, but does not yield any real security.
 - To use root token or regular tokens: when sharing a server with multiple users and associated (disjoint) storage areas, different tokens is needed. In even in single-user mode, use of root token is not recommended.
 - Trusting a third-party server. The holder of the root token (or a group/individual holding of the unseal keys) will be able to read all data stored in Vault.

### How to sync between devices

If you created an App with macpack, then you can simply copy the App, because `ca.crt` and `config.json` will be included. Otherwise these two files need to be copied as well. Even though the token in `config.json` is encrypted, I suggest not storing the config on insecure media like Dropbox or unencrypted mail.

### Possible leaks

When an Apple computer goes into hibernation (not regular sleep), in-memory contents are transferred to disk. I am not entirely sure if the disk data is encrypted if you are not running FileVault. The token could theoretically be scraped here.

## Building

It is as simple as
```
go build
```
which creates a standalone executable. To build a real .App, I suggest using [macpack](https://github.com/murlokswarm/macpack).

The application will try to load your CA certificate `ca.crt`, located in the same folder as the executeable, along with `config.json`.

The file `ca.crt` will be used to authenticate the server you are running Vault on. When you setup your server, you generated a CA. This is the file you need.

The file `config.json` is a file of the format

```
{
	"encrypted": {
		"token": "..."
		"nonce": "..."
		"salt": "..."
	},
	"host": "myserver.com",
	"port": "8001"
}
```

To get the encrypted part, you need to invoke the function `LockToken (plaintext string, password string)` in `crypto.go`:

```
LockToken("your token", "your master password")
```

and put these into the JSON.

## Screenshot

![Pass](pass.gif)

Note: the passwords shown are of course fully working. Do not try to use them!