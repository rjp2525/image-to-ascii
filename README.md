## Image to ASCII

Converts an input image into ASCII art and can output the result in multiple formats (plain text, HTML, and an image). It preserves the aspect ratio of the input image to avoid distortion in the ASCII output, dynamically adjusting based on the input image's dimensions.

## Features

- Converts any PNG image into ASCII art.
- Outputs resized and grayscale images.
- Saves the ASCII art in a `.txt` file, an HTML file, and an ASCII-rendered image.
- Automatically adjusts the output to maintain the image's aspect ratio.

## Usage

You can specify the input image, output paths for the resized image, grayscale image, ASCII text file, HTML file, and ASCII-rendered image via command-line arguments.

```bash
go run main.go -input=input.png -resize_output=resize.png -gray_output=gray.png -ascii_output=result.txt -html_output=result.html -image_output=output.png
```

### Arguments

| Flag             | Description                                     | Default Value   |
|------------------|-------------------------------------------------|-----------------|
| `-input`         | Path to the input PNG image file                | `input.png`     |
| `-resize_output` | Path to save the resized image                  | `resize.png`    |
| `-gray_output`   | Path to save the grayscale image                | `gray.png`      |
| `-ascii_output`  | Path to save the ASCII art in plain text format | `result.txt`    |
| `-html_output`   | Path to save the ASCII art as an HTML file      | `result.html`   |
| `-image_output`  | Path to save the ASCII art rendered as an image | `output.png`    |

## Running Tests

```bash
go test
```

## Requirements

- Go 1.23+
- Go image package

## License

This project is licensed under the [MIT license](LICENSE.md).
