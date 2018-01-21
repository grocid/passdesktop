# Pass for macOS

Pass is a GUI for [pass](https://github.com/grocid/pass), but completely idenpendent of it. It communicates with Vault, where all accounts with associated usernames and passwords are stored.

Pass is password protected locally, by storing an `encrypted token` on disk, along with a `nonce` and `salt`. The `token` is encrypted with AES-GCM-256. The encryption key is derived as 
```key := PBKDF2(password, salt)``` 
and 
```token := AES-GCM-Decrypt(encrypted token, key, nonce)```. 
`token` is kept in memory only.

The program is fairly short and easily auditable, if anyone want to do it.

## Screenshot

![Pass](pass.gif)