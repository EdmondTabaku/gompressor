
# 🤐 Gompressor

Introducing Gompressor, a simple and slow text-based file compressor written in go. 


## ⁉️ How does it work?

The Gompressor's idea is simple: find occurrences in the content and replace them with references of the occurrences content(only if the variables size is smaller than the occurrence itself)

Here is a sample of a compressed content:
```
\e.txt\\d\world \d\hello \d\\o0\\o1\\o0\\o1\\o0\\o1\\o0\\o1\
```
The first variable is the **extension variable**, which will be used when we unzip the file.

 The format is like this: **``` \e[the file extension]\ ```**

 Then there are the *content variables*, which are sortet from the last to the beginning.

 The format is like this: **``` \d\[second variable content]\d\[first variable content]\d\```**

 When all the content variables finish, the content of the file starts which also contains the **content variables references** in their corresponding place.

 The format of the references is like this: **``` \o[index of the variable]\ ```**


 This is not the best way to compress a text file, but it works(expecially when the text file has too much occurrences eg. log files).

 The current implementation that I have done has some problems.

 First - it is very slow 

 Second - when the file has too much content sometimes the decompressor breaks the structure

 Third - works only with text based files




## 📦 Installation

### 1. Pre-built Binaries:

#### Linux:

```bash
# Download the binary
wget https://github.com/EdmondTabaku/gompressor/releases/download/v1.0.0/gompressor-linux -O gompressor

# Make the binary executable
chmod +x gompressor

# Move the binary to a directory in your PATH
sudo mv gompressor /usr/local/bin/
```


#### macOS:

```bash
# Download the binary
curl -Lo gompressor https://github.com/EdmondTabaku/gompressor/releases/download/v1.0.0/gompressor-darwin

# Make the binary executable
chmod +x gompressor

# Move the binary to a directory in your PATH
sudo mv gompressor /usr/local/bin/
```

#### Windows:
Download the .exe file from gompressor.exe.  

Move the .exe file to a directory of your choice.  

Add the directory to your system's PATH.  

Open Command Prompt and verify the installation with ```gompressor```

### 2. Build from Source:
All Operating Systems:
```bash
# Clone the repository
git clone https://github.com/EdmondTabaku/gompressor.git
cd gompressor/cmd

# Build the project
GOOS={your environment} GOARCH=amd64 go build -o gompressor

# Move the binary to a directory in your PATH (you might need sudo for Linux/macOS)
mv gompressor /path/to/directory/in/your/PATH/
```


## 🤝 Contributing
Yes you can contribute




