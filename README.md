# docIt

Hello! This is a tool to  help you draft documentation of your code.

## Installation

To get started you would need an OPENAI_API_KEY. Check the OpenAI websited for details on getting one and pricing. Once that is set fetch the module with : 

```sh
go get github.com/RashikShahjahan/docIt v0.0.0-20230316030938-3af0abae5213
```

## Usage

To get the documentations of file1.py and file2.go run:
```sh
docIt file1.py file2.go
```

This will create corresponding Markdown files containing the draft documentation in the same directory as the files.