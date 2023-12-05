# Jervis: OpenAI chat completion CLI

A simple CLI written using Cobra and Viper to interact with OpenAI's chat completion API.

## Usage

You will need to setup your OpenAI API key in your environment variable.
See OpenAI API documentation for more details.

### Chat command

By default, input taken from stdin, `-n` flag to start a new session
```bash
$ echo hello how are you | jv chat -n
```
Alternatively, you can type or paste in lines of text using `-r` flag.
```bash
$ jv chat -r
```
Without the `-n` flag, we will pick up from the last session and continue the conversation.

By default, the session are savely in the current directory in a file called `.jchat.json`. You can change this with
the `-s {session_name}` flag.

Each session can be converted into a markdown file and viewed using a markdown viewer like `mdcat` or `glow`.
```bash
$ jv chat -f | glow
```

### Edit command

The `edit` command only provide a convenient way to edit text. The options are the same as the `chat` command.

Edit command serves as a quick text editor. You can just provide the text without additional prompt.
```bash
$ cat document.md | jv edit
```

### Code command

The `code` command only provide a convenient way to diagnose code. The options are the same as the `chat` command.

Code command serves as a quick way to diagnose code. You can just provide the code without additional prompt.
```bash
$ cat code.py | jv code
```

### Configuration

You can configure the system prompts and other settings using the `config` command.
```bash
$ jv config
```
This would open the config file using the default editor in your terminal.
The config file is store in `~/.jervis.json` by default.

## Installation

You can build and install locally using
```bash
$ go build -o $GOPATH/bin/jv jervis
```
