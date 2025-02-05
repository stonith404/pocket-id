---
id: user-management
---

# User Management

Creating users in Pocket ID can be handled in two ways:

1. Manually create users via the admin interface.
2. Sync users from an [LDAP](/configuration/ldap) source.

Once users have been created using one of the methods above, follow the steps below to help configure passkeys for them.

## Setting Up User Passkeys

> As the admin, you cannot add passkeys for users; end users must configure them on their own.

> Passkeys can be stored in services like Bitwarden, LastPass, iCloud, or even locally on certain devices using platform authenticators.

### One-Time Link

1. Navigate to the **Users** page in the Pocket ID admin dashboard.
2. Locate the user you want to set up a passkey for.
3. Click the **three dots** on the right side of the user row.
4. Click **One-Time Link**.
5. Select an **Expiration Time** for the link.
6. Click **Generate Link** and send it to the user to allow them to set up their new passkey.

### One-Time Access Email

> **This method requires a valid SMTP server set up in Pocket ID.**

> **Allowing users to sign in with a link sent to their email significantly reduces security, as anyone with access to the user's email can gain entry.**

1. Navigate to the **Application Configuration** section in the Pocket ID admin dashboard.
2. Expand the **Email** section and enable the **Email One-Time Access** option.
3. Instruct the user to navigate to Pocket ID, e.g., `https://id.example.com`.
4. Have the user click on the **Don't have access to your passkey?** link at the bottom of the page.
5. Have the user enter their email associated with their Pocket ID account and click **Submit**.
6. The user will receive an email with a **One-Time Access** link to set up their passkey.
