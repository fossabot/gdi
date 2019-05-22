package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/urfave/cli"
	"golang.org/x/crypto/ssh/terminal"
)

type GDIClient struct {
	url        string
	username   string
	password   string
	HttpClient *http.Client
}

func NewClient(server_url string, server_user string, server_user_pwd string) *GDIClient {
	client := &GDIClient{
		url:        server_url,
		username:   server_user,
		password:   server_user_pwd,
		HttpClient: http.DefaultClient,
	}

	return client
}

func (c GDIClient) Text(text string) (string, error) {
	form := url.Values{}
	form.Add("text", text)

	request, err := http.NewRequest("POST", c.url+"/apps/dropit/text", strings.NewReader(form.Encode()))
	request.SetBasicAuth(c.username, c.password)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	response, err := c.HttpClient.Do(request)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return parseResponse(response)
}

func (c GDIClient) Drop(file *os.File) (string, error) {
	rbody := &bytes.Buffer{}
	writer := multipart.NewWriter(rbody)
	part, _ := writer.CreateFormFile("file", file.Name())
	io.Copy(part, file)
	writer.Close()

	request, err := http.NewRequest("POST", c.url+"/apps/dropit/drop", rbody)
	request.SetBasicAuth(c.username, c.password)
	request.Header.Add("Content-Type", writer.FormDataContentType())
	response, err := c.HttpClient.Do(request)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return parseResponse(response)
}

func parseResponse(response *http.Response) (string, error) {
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode == http.StatusOK {
		var result map[string]string

		if err := json.Unmarshal(body, &result); err != nil {
			log.Fatal(err)
			return "", err
		} else {
			return result["link"], nil
		}
	} else {
		return "", errors.New("Failed, got response code " + strconv.Itoa(response.StatusCode))
	}
}

func main() {
	var server string
	var username string
	var password string
	var filename string

	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "server, s",
			Usage:       "Nextcloud server url, e.g. https://localhost",
			Destination: &server,
		},
		cli.StringFlag{
			Name:        "username, u",
			Usage:       "Nextcloud username",
			Destination: &username,
		},
		cli.StringFlag{
			Name:        "password, p",
			Usage:       "Nextcloud user password",
			Destination: &password,
		},
		cli.StringFlag{
			Name:        "file, f",
			Usage:       "File to upload (reads from STDIN otherwise)",
			Destination: &filename,
		},
	}

	app.Action = func(c *cli.Context) error {
		for _, option := range [2]string{"server", "username"} {
			v := c.String(option)
			if len(v) == 0 {
				err := fmt.Errorf("The mandatory '%s' argument is missing.", option)
				fmt.Println(err.Error())
				os.Exit(1)
			}
		}

		if len(password) == 0 {
			fmt.Printf("Enter %s's password: ", username)
			pwd, err := terminal.ReadPassword(int(os.Stdin.Fd()))
			fmt.Println()
			if err != nil {
				log.Fatal(err)
			}

			if len(pwd) > 0 {
				password = string(pwd)
			} else {
				fmt.Println("The password is mandatory!")
				os.Exit(1)
			}
		}

		client := NewClient(server, username, password)

		fmt.Print("Generating link... ")
		var link string
		if len(filename) > 0 {
			filepath, err := filepath.Abs(filename)
			if err != nil {
				log.Fatal(err)
			}

			file, err := os.Open(filepath)
			if err != nil {
				log.Fatal(err)
			}
			link, err = client.Drop(file)
			defer file.Close()
			if err != nil {
				log.Fatal(err)
			}
		} else {
			input, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				log.Fatal(err)
			}

			link, err = client.Text(string(input))
			if err != nil {
				log.Fatal(err)
			}
		}
		fmt.Println(link)

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
