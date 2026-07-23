# Drop

A simple tool for uploading files to your VPS and serving them from your domain.

## Server Setup

On your VPS, start a `tmux` session and run:

```bash
git clone https://github.com/wazzydotdev/drop
cd drop/server

export DROP_TOKEN="your-token"

go run .
```

Detach from the `tmux` session and your server will continue running.

## Client Setup

On your local machine:

```bash
git clone https://github.com/wazzydotdev/drop

cd drop/cli

go install .

mkdir -p ~/.config/drop

echo "your-token" > ~/.config/drop/key.txt
echo "https://yourdomain.com" > ~/.config/drop/server.txt
```

## Usage

Upload a file:

```bash
drop README.md
```

Output:

```json
{
  "id": "BPOTV9DAR",
  "file": "README.md",
  "url": "/d/BPOTV9DAR"
}
```

Your file will then be available at:

```
https://yourdomain.com/d/BPOTV9DAR
```
