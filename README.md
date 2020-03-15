This is a personal project to learn Golang.

# Working features

- Sending files

# Not implementing

- Compression
- SSL
- File retrieval
- Retries
- Timeout
- Proxy


# About Tentacle

__Tentacle__ is a client/server file transfer protocol that aims to be:

- Secure by design.
- Easy to use.
- Versatile and cross-platform.

Tentacle was created to replace more complex tools like SCP and FTP for simple file transfer/retrieval, and switch from authentication mechanisms like .netrc, interactive logins and SSH keys to X.509 certificates. Simple password authentication over a SSL secured connection is supported too.

Tentacle runs on[IANA](www.iana.org/assignments/port-numbers) assigned port 41121/tcp.

The client and server are designed to be run from the command line or called from a shell script, and no configuration files are needed.

Tentacle is now the default file transfer method for [Pandora FMS](https://pandorafms.org/) and [Babel Enterprise](http://babel.sourceforge.net/en/index.php).

PERL (server and client) and ANSI C/POSIX (client) implementations are available.
