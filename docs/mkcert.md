# MKCERT

## Install

```bash
brew install mkcert
```

## Setup

```bash
mkcert -install
```

## Generate Certs

```bash
mkcert {{TLD}} {{SANs/WILDCARD}} localhost 127.0.0.1 ::1
```

### Example: **(myawesomesite.com)**

```bash
mkcert myawesomesite.com "*.myawesomesite.com" localhost 127.0.0.1 ::1
```
