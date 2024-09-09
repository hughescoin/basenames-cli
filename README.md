# Basenames CLI

Basenames CLI is a command-line interface tool for interacting with the Basenames contract on the Base network. It provides easy-to-use commands for checking balances, availability, expiration, and more.

## Installation

### Prerequisites

- Go 1.21 or higher

### Steps

1. Clone the repository:

   ```
   git clone https://github.com/hughescoin/basenames-cli.git
   cd basenames-cli
   ```

2. Build the CLI:

   ```
   go build -o basenames
   ```

3. Move the executable to your system's PATH:

   For Unix-like systems (Linux, macOS):

   ```
   sudo mv basenames /usr/local/bin/
   ```

   For Windows:
   Add the directory containing the `basenames.exe` to your PATH environment variable.

## Usage

Here are some basic commands to get you started:

1. Check balance:

   ```
   basenames check balance --address 0x1234567890123456789012345678901234567890
   ```

2. Check availability:

   ```
   basenames check availability --tokenId example
   ```

3. Check expiration:

   ```
   basenames check expiration --tokenId example
   ```

4. Get help:
   ```
   basenames --help
   ```

For more commands and detailed usage, please refer to the full documentation.

## Configuration

The CLI uses a configuration file located at `~/.basenames/config.yaml`. You can edit this file to set default values for the RPC URL, private key, and other options.

## TO DO:

- Add versioning to basenamescli
- Add buy basenames if available
- check basename by basenam / not tokenId
- Add session keys
- Name resolution
- Welcome user by Basename if available
- Messaging?
- Add friends
