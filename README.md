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

