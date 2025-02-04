---
id: user-management
---

# User Management

Creating users in Pocket ID can be handled in two ways:

1. Manually create users via the admin interface.
2. Sync users from a [LDAP](/configuration/ldap) source.

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

> **This method requires a valid SMTP server setup in Pocket ID**

> **Allowing users to sign in with a link sent to their email reduces the security significantly as anyone with access to the user's email can gain entry.**

1. Navigate to the `Application Configuration` section in the Pocket ID admin dashboard.
2. Expand the `Email` section and enable the `Email One Time Access` option.
3. Have the User navigate to Pocket ID for example: `https://id.example.com`
4. Have the user click on the `Don't have access to your passkey?` link at the bottom of the page.
5. Have the user enter their email associated with their Pocket ID account and click submit.
6. This will then send a email to the user with a One Time Access link to setup their passkey.