there have been multiple vulnerabilities where exif data including location was public

although they have been fixed, videos sent before the fix still include this metadata

- videos sent through the share button feature on android (fixed ~sep 2022)
- photos sent from iphone (fixed ~2020)
- videos sent from iphone (fixed unknown)

this tool scans your [discord data package](https://support.discord.com/hc/en-us/articles/360004027692-Requesting-a-Copy-of-your-Data) for attachments that include location metadata

hopefully you dont get doxxed like me!!!

# usage

```
go run . [unzipped package path]
```
