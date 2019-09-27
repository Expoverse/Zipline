package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"time"
)

type Config []struct {
	Server struct {
		RemoteSource     string `yaml:"remoteSource"`
		Host             string `yaml:"host"`
		PrivateKey       string `yaml:"privateKey"`
		Username         string `yaml:"username"`
		LocalDestination string `yaml:"localDestination"`
	}
}

func main() {
	configs := Config{}
	source, err := ioutil.ReadFile("config.yml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(source, &configs)
	if err != nil {
		panic(err)
	}

	for _, config := range configs {
		source := config.Server.RemoteSource
		host := config.Server.Host
		privateKey := config.Server.PrivateKey
		username := config.Server.Username
		localDestination := config.Server.LocalDestination

		download(source, host, privateKey, username, localDestination)
	}
}

func clientConfigSetup(keyName string, username string) *ssh.ClientConfig {
	file, err := ioutil.ReadFile("privateKeys/" + keyName + ".pem")
	if err != nil {
		panic(err.Error())
	}

	signer, _ := ssh.ParsePrivateKey(file)
	clientConfig := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	return clientConfig
}

func isOlderThanSixyDays(t time.Time) bool {
	return time.Now().Sub(t) > 1440*time.Hour
}

func download(cmd, hostname string, pem string, username string, destination string) {
	config := clientConfigSetup(pem, username)
	fmt.Println("Backup started... [" + destination + "]")

	conn, err := ssh.Dial("tcp", hostname+":22", config)
	if err != nil {
		panic(err.Error())
	}

	session, err := conn.NewSession()
	if err != nil {
		panic(err.Error())
	}
	defer session.Close()

	r, err := session.StdoutPipe()
	if err != nil {
		panic(err.Error())
	}

	// Make parent directory
	if _, err := os.Stat("backups"); os.IsNotExist(err) {
		_ = os.Mkdir("backups", 0655)
	}

	// Make destination defined directory
	if _, err := os.Stat("backups/" + destination); os.IsNotExist(err) {
		_ = os.Mkdir("backups/"+destination, 0655)
	}

	//Delete backups older than 60 days
	tmpfiles, err := ioutil.ReadDir("backups/" + destination)
	if err != nil {
		return
	}

	for _, file := range tmpfiles {
		if file.Mode().IsRegular() {
			if isOlderThanSixyDays(file.ModTime()) {
				err = os.Remove("backups/" + destination + "/" + file.Name())
				if err != nil {
					panic(err.Error())
				}
			}
		}
	}

	t := time.Now()
	name := fmt.Sprintf("backups/%s/%v.tar.gz", destination, t.Format("2006.01.02.15.04.05"))
	file, err := os.OpenFile(name, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err.Error())
	}
	defer file.Close()

	if err := session.Start("tar -zcf - " + cmd); err != nil {
		panic(err.Error())
	}

	_, err = io.Copy(file, r)
	if err != nil {
		panic(err.Error())
	}

	if err := session.Wait(); err != nil {
		panic(err.Error())
	}

	fmt.Println("Backup finished...")
}
