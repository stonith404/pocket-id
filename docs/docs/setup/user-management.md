---
id: user-management
---

# User Management

Creating users in Pocket ID can be handled in two ways:

1. Manually create users via the admin interface.
2. Sync users from a LDAP source.

Once you have created your users using one of the methods above, follow the steps below to help configure passkeys for the users.

## Setting Up User Passkeys

You as the admin can not add passkeys for users, the end users will have to configure those on their own. Below are the eays the user can do this.

Passkeys can be stored in services like `Bitwarden`, `LastPass`, `iCloud` or Even just locally on certain devices using platform authenticators.

### One Time Link

1. Navigate to the `Users` page in the Pocket ID admin dashboard.
2. Locate the user you want to setup a passkey for.
3. Click the `3 dots` on the right side of the user row.
4. Click `One Time Link`
5. Select a Expiration Time for the link.
6. Click `Generate Link` and you can now send this link to the user to have them setup their new `Passkey`.

### One Time Access Email


- Insert steps here