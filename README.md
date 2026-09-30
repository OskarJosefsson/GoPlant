# GoPlant
CLI-App that should remind you when to water your plants

## Configuration

Copy `.env.example` to `.env` and add your Trefle API key. The app loads `.env` on startup, and `.env` is ignored by Git so the key stays local across restarts.

```sh
cp .env.example .env
```

For a fresh clone, provide your key in the new local `.env` file. Do not commit API keys.
