# shortcut function to whisper. change lang as you wish
# $1 is the audio file. 
# use Audacity to record the sink that you have created with ./start_meeting_mix. Export it in e.g. wav
# use srt or txt
tc ()
{
whisper -f txt --model small.en --language en --device cuda --verbose False "$1"
}

# ollama
# I use qwen3:4b
# prompt as you want
ollama run llama3.2 "Summarize the content of this conversation. Highlight the main topics, decisions and action items" < input.txt > summary.note
