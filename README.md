# autofilterbot
[![telegram badge](https://img.shields.io/badge/Telegram-Channel-30302f?style=flat&logo=telegram)](https://telegram.dog/FractalProjects)
[![Go Report Card](https://goreportcard.com/badge/github.com/Jisin0/autofilterbot)](https://goreportcard.com/report/github.com/Jisin0/autofilterbot)
[![Go Tests](https://github.com/Jisin0/autofilterbot/workflows/Tests/badge.svg)](https://github.com/Jisin0/autofilterbot/actions?query=workflow%3ABuild+event%3Apush+branch%3Amain)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)

**autofilterbot** is a heavily customisable, automatic, fast telegram bot written in [GO](https://go.dev) that can automatically save, filter, sort and index files.

## Commands
```
start       - Check if the bot is alive.
about       - Basic Information About the bot.
help        - Short Guide on How to Use the Bot.
privacy     - Read the user  privacy policy.
settings    - Customise the bot.                              [Admin Only]
broadcast   - Broadcast a message to all users of the bot.    [Admin Only]
batch       - Bunch up messages.                              [Admin Only]
genlink     - Generate link to single file.                   [Admin Only]
logs        - Send app logs.                                  [Admin Only]
index       - Import existing files from a channel.           [Admin Only]
delete      - Assassinate a single file.                      [Admin Only]
deleteall   - Massacre all matching files.                    [Admin Only]
```

## Features
The bot is jam-packed with useful features and tools to maximise productivity.
- [x] Automatically Save & Filter Files
- [x] Import Large Volumes of Files.
- [x] Batch Files into a Single Link.
- [x] 25+ Configuration Options
- [x] Multiple Databases
- [ ] URL Shortener
- [x] Force Subscribe Channels
- [x] Broadcast Messages
- [x] Configuration Panel
- [x] Select Multiple Files
- [x] Request Fsub
- [x] Auto Delete

## Variables
The variables below can be configured by setting them as environment variables, or adding them to a .env file at the root of the project.
[Sample .env file](https://github.com/Jisin0/autofilterbot/tree/main/.env.sample) can be found at the root of the repository. Remember to name the file .env

### Required
- `BOT_TOKEN`     : Bot token obtained from [@botfather](https://t.me/botfather) by running the /newbot command.
- `ADMINS`        : List of telegram ids of bot admins separated by whitespaces. IDs can be obtained using [@myidbot](https://t.me/myidbot).
- `MONGODB_URI`   : Mongodb cluster uri form mongodb atlas. Watch [this video](https://www.youtube.com/watch?v=SMXbGrKe5gM) to learn how to create one. Multiple urls can be added by setting MONGODB_URI1, MONGODB_URI2 etc. Note: The main database will be MONGODB_URI, this is where all configuration and user data will be saved. Secondary databases are only used to save files. The database to save files to can be changed from settings.
- `FILE_CHANNELS` : List of telegram ids of channels where new files will be posted separated by whitespaces. Should be in the format -100xxxxxxxxx. IDs can be obtained using [@myidbot](https://t.me/myidbot).

### Optional
- `LOG_LEVEL` : Level of logs to be output. Possible values are debug, info, warn, error. Please set to debug if reporting an issue or it is recommended to be left at info.
- `APP_ID`    : Telegram MTProto App ID from my.telegram.org.
- `APP_HASH`    : Telegram MTProto App Hash from my.telegram.org.

## Deploy
Deploy your bot to any server or VPS of choice. The project comes with a plethera of pre-built platform-specific configurations.

<details><summary>Deploy To Heroku</summary>
<p>
<br>
<a href="https://heroku.com/deploy?template=https://github.com/Jisin0/autofilterbot/tree/main">
  <img src="https://www.herokucdn.com/deploy/button.svg" alt="Deploy">
</a>
</p>
</details>

<details><summary>Deploy To Scalingo</summary>
<p>
<br>
<a href="https://dashboard.scalingo.com/create/app?source=https://github.com/Jisin0/autofilterbot#main">
   <img src="https://cdn.scalingo.com/deploy/button.svg" alt="Deploy on Scalingo" data-canonical-src="https://cdn.scalingo.com/deploy/button.svg" style="max-width:100%;">
</a>
</p>
</details>


<details><summary>Deploy To Render</summary>
<p>
<br>
<a href="https://dashboard.render.com/select-repo?type=web">
  <img src="https://render.com/images/deploy-to-render-button.svg" alt="deploy-to-render">
</a>
</p>
<p>
Make sure to have the following options set :

<b>Environment</b>
<pre>Go</pre>

<b>Build Command</b>
<pre>go build .</pre>

<b>Start Command</b>
<pre>./autofilterbot</pre>

<b>Advanced >> Health Check Path</b>
<pre>/</pre>
</p>
</details>


<details><summary>Deploy To Koyeb</summary>
<p>
<br>
<a href="https://app.koyeb.com/deploy?type=git&repository=github.com/Jisin0/autofilterbot&branch=main">
  <img src="https://www.koyeb.com/static/images/deploy/button.svg" alt="deploy-to-koyeb">
</a>
</p>
<p>
You must set the Run command to :
<pre>./bin/autofilterbot</pre>
</p>
</details>

<details><summary>Deploy To Okteto</summary>
<p>
<br>
<a href="https://cloud.okteto.com/deploy?repository=https://github.com/Jisin0/autofilterbot">
  <img src="https://okteto.com/develop-okteto.svg" alt="deploy-to-okteto">
</a>
</p>
</details>

<details><summary>Deploy To Railway</summary>
<p>
<br>
<a href="https://railway.app/new/template?template=https%3A%2F%2Fgithub.com%2FJisin0%2Fautofilterbot">
  <img src="https://railway.app/button.svg" alt="deploy-to-railway">
</a>
</p>
</details>

<details><summary>Run Locally/VPS</summary>
<p>
You must have the latest version of <a href="https://go.dev/dl">GO</a> installed first
<pre>
git clone https://github.com/Jisin0/autofilterbot
cd autofilterbot
go build .
./autofilterbot
</pre>
</p>
</details>

## Contributing
All contributions and enhancements are welcome and we'll need all the help we can get. 

Please remember to follow [conventional commits](https://conventionalcommits.org) when contributing. Feel free to contact me [@FractalProjects](https://telegram.dog/FractalDiscussions) for suggestions, reporting issues or for a quick review before contributing.

If reporting an issue, please provide the app logs (can be obtained using /logs) after setting the LOG_LEVEL to debug for easy debugging.

## Thanks
 - Thanks to Paul for his awesome [Library](https://github.com/PaulSonOfLars/gotgbot)
 - Thanks To [AlbertEinsteinTG](https://github.com/AlbertEinsteinTG) for his awesome [project](https://github.com/AlbertEinsteinTG/Adv-Auto-Filter-Bot-V2) that inspired this.

## Disclaimer
Any content distributed using this project will be the liability of the user alone. Any code used from the project must be clearly cited.

[![GNU General Public License 3.0](https://www.gnu.org/graphics/gplv3-127x51.png)](https://www.gnu.org/licenses/gpl-3.0.en.html#header)    
Licensed under [GNU GPL 3.0.](https://github.com/Jisin0/autofilterbot/blob/main/LICENSE).