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

This project uses settings.json file with such format:

```json
{
    "app_id": "app_id",
    "api_version": "5.41"
}
```
where "app_id" is vk application id and "api_version" is vk API version. The program will fail without this file.

This project is a console application. To run it type:

```bash
./vk_public_saver #on *NIX
./vk_public_saver.exe #on Windows
go run *.go #or run as a script if you have Golang environment
```

There are 3 modes:

- user login/password input, where program asks about username, password and public or user id to save it's wall data;
- user token input, where program asks about token and public or user id to save it's wall data;
- json input, where program uses json file userdata.json;

The first mode is the default one. For token input run

```bash
./vk_public_saver --token

```

For using json input run vk_public_saver with --userdata option. For example:

```bash
./vk_public_saver --userdata

```

In this mode there should be json file userdata.json with such format:

```json
{
    "email": "my@mail.ru",
    "pass": "MyStrongSecretPassword",
    "source": "MyPublicID"
} 
```
where "email" is user email, "pass" is user password, "source" is public or user id which wall we are saving.

While running vk_public_saver will create data_"YourPublicID".db file with sqlite database storing posts and some attachments info:

```sql
CREATE TABLE IF NOT EXISTS posts (
    id INTEGER,
    from_id INTEGER,
    owner_id INTEGER,
    signer_id INTEGER,
    date INTEGER,
    marked_as_ads INTEGER,
    post_type TEXT,
    text TEXT,
    is_pinned INTEGER,
    comments_count INTEGER,
    likes_count INTEGER,
    reposts_count INTEGER,
    views_count INTEGER,
    attachments_count INTEGER,
    PRIMARY KEY (id, from_id)
);

CREATE TABLE IF NOT EXISTS attachments (
    type TEXT,
    id INTEGER,
    owner_id INTEGER,
    post_id INTEGER,
    url TEXT,
    additional_info text,
    additional_info2 text,
    PRIMARY KEY (id, type, post_id)
);
```

## License

This program uses GNU GPL3. For more information see the LICENSE file.

<!-- ## Known usage

TODO: write it -->

## Thanks

This project uses some additional packages:

- EasyVK (https://github.com/Vorkytaka/easyvk-go) for login in vk.com;
- golang.org/x/crypto/ssh/terminal for password input.
