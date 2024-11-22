# Changelog

## [2.0.0] - 2024-11-22

This release changes the database schema significantly. Make a backup.

### Main changes

- Use more readable environment variable names. Old names are temporarily supported.
- Store sessions in the database, this improves security. Secret key is no longer needed.
- Add support for multiple users. Link groups now have an owner. The registration page can be disabled.

### Additional changes

- Make margins in tile mode smaller.
- Use familiar names for web modules.
- Replace logrus with log/slog.
- Use username as a primary key.
