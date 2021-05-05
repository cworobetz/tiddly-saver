# tiddly-saver

tiddly-saver is a helper program for Tiddlywiki. It moves your saved Tiddlywiki from a watched path to your desired save folder.

## Movitation

When using Tiddlywiki on Windows 10 Chrome, the default save location is `%USERPROFILE%/Downloads/<wiki-name>.html`. This makes sense, because you're downloading an HTML page and the default save location is in a user's download folder. 

However, I save my Tiddlywiki in `%USERPROFILE%/Tiddlywiki`. It is then synced from there to my other devices so I can access it anywhere. With my syncing program, I can't select individual files from a folder to sync - I have to sync the whole folder.

This is where this program comes in: It runs in the background, waiting for the watched file path to be present, e.g. `%USERPROFILE%/Downloads/<wiki-name>.html`. When it see writes to this file, it waits a user-configurable amount of time for the writes to finish (e.g. 2 seconds), then it moves the file to the desired location, e.g. `%USERPROFILE%/Tiddlywiki/<wiki-name>.html`

## Installation

Coming soon

## Usage

Coming soon

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License
[GNU AGPLv3](https://choosealicense.com/licenses/agpl-3.0/)