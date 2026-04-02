# dnsmixer

`dnsmixer` is a fast, lightweight Go-based DNS permutation tool. It generates potential subdomains based on existing subdomains and optional wordlists, helping security researchers and bug hunters discover hidden assets.

## Features

- **Subdomain Permutation:** Automatically mixes and matches existing subdomains.
- **Wordlist Support:** Injects custom words into existing subdomain structures.
- **Level Filtering:** Option to target only first-level subdomains (`-l1`).
- **Scoped Execution:** Target specific subdomains and their children (`-sc`).
- **Live Progress:** Real-time counter showing remaining permutations to be generated.
- **Clean Output:** Designed for easy piping into other tools like `httpx` or `dnsx`.

## Installation
```
git clone https://github.com/codejavu-inc/dnsmixer.git
cd dnsmixer
bash install.sh
sudo cp dnsmixer /usr/local/bin # if you want to execute dnsmixer from any directory
```

# Usage
```
~ dnsmixer -h
Usage: dnsmixer [options]

Options:
  -w <file>   Wordlist file (if omitted, only permutes existing subdomains)
  -t <file>   Target list of subdomains (Required)
  -l1         Process only first-level subdomains
  -sc <sub>   Scope execution to only this subdomain and its children
  -o <file>   Output file path to save results (default is stdout)
```
