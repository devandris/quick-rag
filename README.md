# shortcut function to whisper. change lang as you wish

- `$1` is the audio file.
- use Audacity to record the sink that you have created with `./scripts/start_meeting_mix`. Export it in e.g. `wav`
  ```bash
  tc ()
  {
  whisper -f srt --model small.en --language en --device cuda --verbose False "$1"
  }
  ```

# ollama

- I use qwen3:4b
- prompt as you want

```bash
ollama run llama3.2 "Your prompt. e.g.: Summarize me the meeting by listing the main topics, the agreements and next actions" < input.txt
```

# Next

[build your own model.](https://www.hostinger.com/tutorials/ollama-cli-tutorial#Creating_custom_models)
