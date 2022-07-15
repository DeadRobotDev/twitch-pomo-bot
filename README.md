# twitch-pomo-bot
A Twitch chat bot for co-working streamers. Viewers can set their own task. Streamers can display them by adding a simple text source to OBS.

## Installation
- Download the [latest release](https://github.com/DeadRobotDev/twitch-pomo-bot/releases).

**NOTE:** Microsoft Defender Antivirus may flag the executable. It is a false positive, but you can read and build the source code yourself if you're unsure.

- Run `twitch-pomo-bot.exe`.

**WARNING:** Do **NOT** stream this part, or the `config.json` file. It includes an OAuth token that allows user to post to any Twitch chat from your bot account without logining in.

It will take you through the initial set up process, and create a `config.json` file. You can edit this manually if you wish to change the command prefix (default: `!`), or the bot responses.

### Display in OBS
- Add a `Text` source.
- Select `Read from File`.
- Select `Browse`.
- Select the `viewer_tasks.txt` file.

## Commands
- `!task` - Replies with the help message.
- `!task add <task name>`
- `!task edit <new task name>`
- `!task done` OR `!task complete`
- `!task cancel` OR `!task delete`

## Planned Features
- Mod commands.
- Timed tasks (aka pomodoros).
- Disable individual bot responses.

## Credits
- Thanks to [RumpleStudy](https://www.twitch.tv/RumpleStudy) for the inspiration and co-working streams.
- Thanks to [SeijiSoldier](https://www.twitch.tv/SeijiSoldier) for the word, logining, and co-working streams.

## License
[MIT](LICENSE)
