# OS-level Password Storage

This document explains where and how passwords are stored locally on different operating systems in this project, the file format used, and important security considerations.

## Summary

- File name: `genp.yaml`
- Top-level YAML key: `password`
- Each stored secret is a key-value pair under `password`
- Directory permissions (Unix-like): `0700`
- File permissions (Unix-like): `0600`
- Storage locations follow OS conventions for per-user application data

## Storage Locations

The base application folder is `genp`. The password file is `genp.yaml` inside that folder.

- Windows
  - Preferred: `%LOCALAPPDATA%\genp\genp.yaml`
  - Fallback: `%APPDATA%\genp\genp.yaml`
  - Example: `C:\Users\<User>\AppData\Local\genp\genp.yaml`

- macOS
  - `~/Library/Application Support/genp/genp.yaml`

- Linux and other Unix-like systems
  - Preferred: `$XDG_CONFIG_HOME/genp/genp.yaml`
  - Fallback: `~/.config/genp/genp.yaml`

Note: Dotfiles are considered “hidden” on Unix-like systems. We use the standard XDG/macOS application support directories instead of dotfiles, which are not typically visible in standard file browsers. On Windows, there is no dotfile convention; the directory lives under the user’s `%LOCALAPPDATA%` or `%APPDATA%`.

## File Format

Passwords are stored in YAML as a map under the top-level key `password`. Each password entry is added by name.

Example:
```
password:
  email: "hunter2"
  github: "ghp_example_token"
  db_admin: "S3cur3P@ss!"
```

Characteristics:
- The keys under `password` are human-readable identifiers you provide when storing.
- The values are the corresponding secrets.
- Subsequent additions append or update entries under the `password` map.

## Permissions and Access Control

- On Unix-like systems (Linux/macOS), the directory is created with `0700` permissions and the file is written with `0600` permissions. This typically restricts access to the current user only.
- On Windows, file mode flags are best-effort. Ensure the directory resides in the per-user profile (`%LOCALAPPDATA%` or `%APPDATA%`). For stronger protection, use OS-native encryption (DPAPI) or the Credential Manager instead of plain files.

## Security Considerations

- Plain file storage is convenient but not inherently secure. Prefer OS-native secure stores whenever practical:
  - Windows: Credential Manager and DPAPI (`CryptProtectData`)
  - macOS: Keychain (`SecItemAdd`, `SecItemCopyMatching`)
  - Linux: Secret Service (`libsecret`/gnome-keyring) or KWallet
- If you must store secrets in a file:
  - Apply strict permissions (per-user only).
  - Consider encrypting contents at rest using a well-maintained library (AES-GCM or ChaCha20-Poly1305 with a unique random nonce and a strong key).
  - Do not hardcode encryption keys. Derive or retrieve keys from the OS keystore when possible.
  - Avoid storing long-lived secrets in environment variables; they may leak via logs or process inspection.
  - Exclude the storage path from backups that could be accessible to others, or ensure backups are encrypted.
  - Never commit `genp.yaml` to version control. Add it to ignore lists.
- Implement rotation and revoke procedures:
  - Provide a way to update or remove individual entries under `password`.
  - When a secret is rotated, write the new value and confirm dependent systems use the updated secret.
  - If the application is uninstalled, consider removing the local file.

## Operational Notes

- The code determines the correct per-OS base directory and ensures it exists with restrictive permissions.
- When adding a new password:
  - If `genp.yaml` does not exist, it is created with the proper header and the first entry.
  - If it exists, a new entry is appended or updated under the `password` key.
- The implementation avoids using a non-standard hidden path on macOS and Linux, in favor of standard application data locations, making lifecycle management consistent with OS expectations.

## Future Enhancements

- Use an actual YAML library to safely read/modify the document structure and handle edge cases (comments, formatting, duplicate keys).
- Add conflict detection for duplicate names and provide update vs. overwrite behavior.
- Offer optional encryption-at-rest with key management via OS keystores.
- Add export/import features that maintain structure and enforce security controls.