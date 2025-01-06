#!/bin/sh

# --batch to prevent interactive command
# --yes to assume "yes" for questions
gpg --quiet --batch --yes --decrypt --passphrase="$PROJECT_YML_CRYPT_KEY" --output project.yml project.yml.gpg
