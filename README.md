**Disclaimer 1:** This tool is currently in an alpha version. It is not intended for production use. Use it at your own risk.  
**Disclaimer 2:** This README was created with the assistance of ChatGPT. The source code **is not**

# locker 

locker is a lightweight command-line tool for encrypting and decrypting data and files.
It can securely encrypt sensitive text (like passwords) as well as files.

Unlike traditional tools, locker does not require a password for encryption.
Instead, it uses an asymmetric RSA key pair:

When you create a login, locker generates a public/private RSA key pair.

The public key is used to encrypt a randomly generated, cryptographically secure password.

To decrypt, you authenticate with your login password, which unlocks your private key.
This private key decrypts the random password, which is then used to unlock your data.

You can create multiple logins, each encryption the same private key

## Features  

- 🔑 Asymmetric encryption (encryption without a password)
- 🔒 Store encrypted text data (e.g. passwords, notes)
- 📂 Encrypt and decrypt files
- 👥 Support for multiple logins
- 🧹 Automatically deletes the original file after encryption
- 💻 Cross-platform support: works on Windows, Linux, and macOS
(currently tested on Windows and Linux)

### Getting Started
1. Create your first login

Before encrypting anything, you need to create a login with a password:

`.\locker.exe add login --newlogin jakob`


This generates a new RSA key pair and associates it with your login.

2. Encrypt data

You can encrypt text (e.g. a password):

`.\locker.exe encrypt data --name "TestPass" --data "mysecretPass"`


To view all stored encrypted data:

.\locker.exe show data

3. Decrypt data

Each entry has an ID. Use it with your login to decrypt:

`.\locker.exe decrypt data --login jakob --id c1e54db5-6197-4893-b132-3f36cce40498`

4. Encrypt files

Encrypt a file by specifying the source path:

`.\locker.exe encrypt file --source "C:\Users\Jakob\Downloads\memtest86-usb.zip" --destination C:\Users\Jakob\Desktop`


The encrypted file will have the .lock extension.

If no destination is provided, the encrypted file is created in the same directory as the source.

After encryption, the original source file is deleted.

5. Decrypt files

To restore an encrypted file:

`.\locker.exe decrypt file --source "C:\Users\Jakob\Desktop\memtest86-usb.zip.lock" --destination C:\Users\Jakob\Downloads\ --login jakob`


Roadmap 

Road to 100% Unit/Integration tests 



License

MIT License – feel free to use and contribute.
