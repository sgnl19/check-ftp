# check-ftp
Icinga Checks based on ftp commands

## Default Flags



| Flag | Short | Description |
| --- | --- | --- |
| COMMAND |  | The command to execute |
| --help | -h | Help for check-ftp |
| --verbose | -v | Enable verbose output |

## Commands

### user-restriction

Check restrictions for a given user.

- User is able to log in`
- User is located in his home directory
- User cannot navigate outside his home directory

If any restriction is violated a critical is thrown

```
check-ftp user-restriction [ftp-host]:[port] -u [user] -p [password]
```

### Flags

| Flag | Short | Description |
|--- |--- |--- |
| --user | -u | the ftp user |
| --password | -p | the password |
| --help | -h | Help for user-restriction |
