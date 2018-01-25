![Travis](https://travis-ci.org/grocid/PassMacOS.svg?branch=master)

# Pass Desktop for macOS

Pass is a GUI for [pass](https://github.com/grocid/pass), but completely independent of it. It communicates with [Vault](https://www.vaultproject.io), where all accounts with associated usernames and passwords are stored.

Pass is password protected on the local computer, by storing an `encrypted token` on disk, along with a `nonce` and `salt`. The `token` is encrypted with AES-GCM-256. The encryption key is derived as 
```
key := PBKDF2(password, salt)
```
and 
```
token := AES-GCM-Decrypt(encrypted token, key, nonce).
```

The decrypted `token` is kept in memory only.

![Decrypting token](decryptingtoken.png)

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
 - To use root token or regular tokens: when sharing a server with multiple users and associated (disjoint) storage areas, different tokens are needed and, hence, root token cannot be used. In even in single-user mode, use of root token is not recommended.
 - Trusting a third-party server. The holder of the root token (or a group/individual holding of the unseal keys) will be able to read all data stored in Vault. However, Vault is quite light weight and can, for a limited amount of users, be run on a mere Raspberry Pi gen A. I would suggest that each user runs Vault on their own Raspberry Pi at home. Secrets can be shared over several VPS instances and providers using secret sharing. While at a higher cost, it would give higher security and accessibility (as a e.g. (3, 2) scheme would require only two out of three servers to be online).
 - PBKDF2 computes 4096 iterations of SHA-256. If the password has a lot lower entropy than 256 bits, then the iteration count need to be increased considerably.

### How to sync between devices

If you created an App with macpack, then you can simply copy the App, because `ca.crt` and `config.json` will be included. Otherwise these two files need to be copied as well. Even though the token in `config.json` is encrypted, I suggest not storing the config on insecure media like Dropbox or unencrypted mail.

### Possible leaks

When an Apple computer goes into hibernation (not regular sleep), in-memory contents are transferred to disk. The Filevault key itself can be recovered when device is in sleep (default for desktops), or deep/hybrid sleep, or whatever they call it (default for laptops), so presumably the token is also somewhere. The difference (in description) between hybrid and full hibernation is only that memory power is disconnected, so it might mean that in the hibernation file you also have all keys. And there is also a, notably non-standard, option to wipe the Filevault key each time the system goes to standby. By default Apple do not wipe the key -- this is usability consideration -- user gets faster response from system and no annoying passwords are needed.

There is the ```pmset somethingVaultKeysomethingsomething``` setting. If you are concerned about this, I suggest you do some own research.

## Performance

Pass Desktop keeps no information stored on disk. Search operations are done by performing a LIST (Hashicorp-specific operation), which fetches a JSON with all keys (account names) from the server, after which filtering operations are performed locally. This is done every search query, so with a slow server, the user experience may not be as intended. The same applies for very long lists of accounts. Moreever, Pass Desktop keeps an iconset, where each filename is associated with the account name (favicons are too small). Since there is a mapping betwen account names and the iconset, the recommended convention is to name accounts after the domain. The iconset can be extended by the user with minor effort. The memory usage is about 50 MBs of RAM.

## Building

It is as simple as
```
go build
```
which creates a standalone executable. To build a real .App, I suggest using [macpack](https://github.com/murlokswarm/macpack).

The application will try to load your CA certificate `ca.crt`, located in the same folder as the executeable, along with `config.json`.
The file `ca.crt` will be used to authenticate the server you are running Vault on. When you setup your server, you generated a CA. This is the file you need.
The configuration `config.json` is a file of the format

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

## Setting up the backend

To get Pass working, you need to install and configure Vault on the remote server. First, start the storage backend for Vault. This can be SQL, but I would recommend [Consul](https://www.consul.io). Start Consul as follows:
```
consul agent -server -config-dir=/etc/consul.d/bootstrap/ > /dev/null 2>&1 &
```
Let the contents of `/etc/consul.d/bootstrap/config.json` be
```
{
    "bootstrap": true,
    "server": true,
    "datacenter": "pass",
    "data_dir": "/var/consul",
    "encrypt": "<secret>",
    "ca_file": "/etc/consul.d/ssl/ca.crt",
    "cert_file": "/etc/consul.d/ssl/consul.crt",
    "key_file": "/etc/consul.d/ssl/consul.key",
    "verify_incoming": true,
    "verify_outgoing": true,
    "log_level": "INFO",
    "enable_syslog": false
}
```
Then, Vault can be started in the following way.
```
vault server -config=/etc/vault.d/config.json > /dev/null 2>&1 &
```
where `/etc/vault.d/config.json` contains
```
{
    "storage": {
        "consul": {
            "address": "127.0.0.1:8500",
            "advertise_addr": "https://127.0.0.1:8200",
            "path": "vault"
        }
    },
    "listener": {
        "tcp": {
            "address": "127.0.0.1:8200",
            "tls_cert_file": "/etc/vault.d/ssl/vault.crt",
            "tls_key_file": "/etc/vault.d/ssl/vault.key",
            "tls_disable": 0
        }
    }
}
```
Vault binds to `127.0.0.1`, so we need fw to access it (if you do not want to use fw, then bind the listener to `0.0.0.0:8001`). We can start it as
```
./fw 0.0.0.0:8001 127.0.0.1:8200 2>&1 &
```

The more proper way to start Consul, Vault and fw would be to create an init.d or systemd service, but this is fine for testing purposes.

## Screenshots

![Config](doc/config.png)
![Unlock](doc/unlock.png)
![Search](doc/search.png)
![Account](doc/account.png)
![Trash](doc/trash.png)
![Pass](doc/pass.gif)

Note: the passwords shown are of course fully working. Do not try to use them!