## Zipline

The community is welcomed to ask questions, contribute, and make issues, and pull requests.

#### Description
Go program to securely archive and backup data from multiple servers.

This software tarballs (backup.tar.gz) the remote folders and then downloads to local destination.

#### Installing Zipline (Linux)
```bash
$ curl -sSL https://git.io/zipline | bash
```
  

#### Usage
To create backups edit the config.yml file. An example of the config.yml file below:
```
- server:
    remoteSource: ~/example/directory/products
    host: 10.0.2.1
    privateKey: id_rsa_testing
    username: ubuntu
    localDestination: products
- server:
    remoteSource: ~/docker/apache/app/var/www/members
    host: 52.26.27.120
    privateKey: id_rsa_production
    username: root
    localDestination: members
```
The yml contains an array of servers with properties. Each server should contain:
- **remoteSource** This is the remote directory src
- **host** This is the IP address of the remote machine assuming the port is ***22***
- **privateKey** This is the name of the .pem file in the directory created by Zipline. The directory
is called ***privateKeys***
- **username** Username to log into the machine via ssh
- **localDestination** The child directory name of the location to store the backups in ***backups*** 
folder created by Zipline.

#### Creating more servers
To create more servers add another server definition to the list in the config.yml file. Fill in the 
properties according to your setup.
```
- server:
    remoteSource:
    host:
    privateKey:
    username:
    localDestination:
```

