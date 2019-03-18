# vk public saver

Project for saving vk posts

## Requirements

This project is written in Go language so you need modern version of Go to build and run this application. You can find instructions and some more info here - https://golang.org/doc/install

You should also install some additional packages:

```bash
golang.org/x/crypto/ssh/terminal
github.com/vorkytaka/easyvk-go/easyvk
```

## Building

For build binary executable run

```bash
go build
```

You can also run vk_public_saver as script:

```bash
go run *.go
```

## Usage

This project is a console application. To run it type:

```bash
./vk_public_saver #on *NIX
./vk_public_saver.exe #on Windows
go run *.go #or run as a script
```

TODO: write about settings.json file

There are 2 modes:

- user input, where program asks about username, password and public or user id;
- json input, where program uses json file userdata.json in such format:

```json
{
    "email": "my@mail.ru",
    "pass": "MyStrongSecretPassword",
    "source": "MyPublicID"
} 
```

The first mode is the default one. For using json input run vk_public_saver with --userdata option. For example:

```bash
./vk_public_saver --userdata

```

While running vk_public_saver will create data_"YourPublicID".db file with sqlite database storing posts and some attachments info. See [db](db.go) file for additional info about database structure and tables' list. TODO: write this info here.

## License

This program uses GNU GPL3. For more information see the LICENSE file.

## Known usage

TODO: write it

## Thanks

TODO: This project uses EasyVK package for login  
https://github.com/Vorkytaka/easyvk-go

golang.org/x/crypto/ssh/terminal
