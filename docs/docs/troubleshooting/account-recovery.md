---
id: account-recovery
---

# Account recovery

There are two ways to create a one-time access link for a user:

1. **UI**: An admin can create a one-time access link for the user in the admin panel under the "Users" tab by clicking on the three dots next to the user's name and selecting "One-time link".
2. **Terminal**: You can create a one-time access link for a user by running the `scripts/create-one-time-access-token.sh` script. To execute this script with Docker you have to run the following command:
   ```bash
   docker compose exec pocket-id sh "sh scripts/create-one-time-access-token.sh <username or email>"
   ```
