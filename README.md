# Ki-Hocche

Ki-Hocche is a utility to help you find out what is happening in your SCM repositories. It creates a icalendar file from your SCM repositories. It is a simple tool that can be used to keep track of what is happening in your repositories.

## Installation

```bash
go install github.com/ki-hocche/ki-hocche@latest
```

## Usage

```bash
ki-hocche [flags]
```

## Flags

```sh
./ki-hocche --help
```

### Export to icalendar file

```bash
export GITHUB_TOKEN=<your_github_token>
./ki-hocche journey --repos=repo_1 --repos=repo_2 -t ics -o /tmp/mycal.ics;
```

### Export to JSON

```bash
export GITHUB_TOKEN=<your_github_token>
./ki-hocche journey --repos=repo_1 --repos=repo_2 -t json -o /tmp/events.json;
```

## Supported SCMs

- [X] Github
- [ ] Gitlab
- [ ] Bitbucket

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.
