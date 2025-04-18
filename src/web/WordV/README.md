---

# WordPress Version Finder 

That’s a tool written in Rust for finding WordPress versions from their RSS feed. It’s mainly for grabbing the `<generator>` tag from an RSS feed, which usually includes the WordPress version — and that’s super handy for vulnerability checks using CVE databases like [cve.mitre.org](https://cve.mitre.org) or [OpenCVE](https://www.opencve.io).

Still in progress, probably gonna add other methods for finding versions, and maybe even extend it to detect Apache versions (and their respective CVEs) too. 👀

## What it does

- Fetches the RSS feed from a WordPress site (currently hardcoded to `https://awesomemotive.com/feed`)
- Saves the feed locally as `feed.xml`
- Parses the feed to find the `<generator>` tag
- If it can't parse it properly, it tries a simpler fallback: line-by-line search
- Prints the WordPress version (if found)

## Example Output

```bash
File feed.xml saved successfully
Wordpress generator content:
 https://wordpress.org/?v=6.4.3
```

Or, if parsing fails:

```bash
Error while getting content, trying in another way...
Error: ...
Version found at line <generator>https://wordpress.org/?v=6.4.3</generator>
```

## How to run it

Make sure you’ve got Rust and Cargo set up, then just:

```bash
cargo run
```

Easy as that. 

##  What's next

- Add more ways to extract version info (e.g. from meta tags, headers, etc.)
- Apache version detection
- Maybe a CLI argument for custom URLs
- Possibly automatic CVE lookups 👁
---

# **NOTE**:
It does not work for every wordpress page
