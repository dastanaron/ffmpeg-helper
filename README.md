# FFmpeg Helper

FFmpeg Helper is a command-line tool that simplifies the use of FFmpeg by providing a user-friendly interface for executing predefined FFmpeg commands. It allows users to select from a list of commands, modify them, and execute them with ease.

## Features

- **Predefined Commands**: Comes with default FFmpeg commands for common tasks.
- **Custom Commands**: Users can add their own FFmpeg commands via a configuration file.
- **Interactive UI**: A terminal-based UI for selecting and executing commands.
- **Progress Tracking**: Real-time progress bar for FFmpeg operations.

## Installation

1. **Install Go**: Ensure you have Go installed on your system. You can download it from [here](https://golang.org/dl/).

2. **Clone the Repository**:
   ```bash
   git clone https://gitlab.com/Dastanaron/ffmpeg-helper.git
   cd ffmpeg-helper
   ```

3. **Build the Program**:
   ```bash
   go build -o ffmpeg-helper
   ```

4. **Run the Program**:
   ```bash
   ./ffmpeg-helper input.mp4 output.mp4
   ```

   Replace `input.mp4` and `output.mp4` with your actual input and output file paths.

## Configuration

### Default Commands

On the first run, the program creates a configuration file with two default commands:

1. **Telegram**: Converts a video for Telegram.
   ```yaml
   - name: 'Telegram'
     cmd: 'ffmpeg -y -i {if} -pix_fmt yuv420p -codec:a aac -c:v libx264 {of}'
     description: 'Convert video for Telegram app'
   ```

2. **MP3 Audio**: Extracts audio from a video and saves it as an MP3 file.
   ```yaml
   - name: 'MP3 Audio'
     cmd: 'ffmpeg -y -i {if} -vn -ar 44100 -ac 2 -ab 192K -f mp3 {of}'
     description: 'Cut audio in mp3 format'
   ```

### Adding Custom Commands

You can add your own FFmpeg commands by editing the configuration file located at:
`~/.config/ffmpeg-helper/commands.yaml`

#### Example: Adding a Custom Command

1. Open the `commands.yaml` file in a text editor.

2. Add a new command in the following format:
   ```yaml
   - name: 'Custom Command'
     cmd: 'ffmpeg -y -i {if} -vf "scale=1280:720" {of}'
     description: 'Resize video to 1280x720'
   ```

   - `name`: The name of the command (displayed in the UI).
   - `cmd`: The FFmpeg command. Use `{if}` for the input file and `{of}` for the output file.
   - `description`: A brief description of the command (displayed in the UI).

3. Save the file and restart the program. Your custom command will now appear in the list.

## Usage

1. **Run the Program**:
   ```bash
   ./ffmpeg-helper input.mp4 output.mp4
   ```

2. **Select a Command**:
   - Use the arrow keys to navigate the list of commands.
   - Press `Enter` to execute the selected command.

3. **Monitor Progress**:
   - A progress bar will display the status of the FFmpeg operation.

4. **Quit the Program**:
   - Press `q` to exit the program.

## Example Commands

Here are some example commands you can add to your `commands.yaml` file:

### Convert to GIF
```yaml
- name: 'GIF'
  cmd: 'ffmpeg -y -i {if} -vf "fps=10,scale=320:-1" {of}'
  description: 'Convert video to GIF'
```

### Extract Frames
```yaml
- name: 'Extract Frames'
  cmd: 'ffmpeg -y -i {if} -vf "fps=1" frame_%04d.png'
  description: 'Extract frames from video'
```

### Add Watermark
```yaml
- name: 'Watermark'
  cmd: 'ffmpeg -y -i {if} -i watermark.png -filter_complex "overlay=10:10" {of}'
  description: 'Add watermark to video'
```

## Contributing

If you'd like to contribute to this project, please follow these steps:

1. Fork the repository.
2. Create a new branch for your feature or bugfix.
3. Submit a pull request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

This `README.md` provides a clear and concise guide for users to get started with your FFmpeg Helper tool. It explains how to install, configure, and use the program, as well as how to add custom commands.
