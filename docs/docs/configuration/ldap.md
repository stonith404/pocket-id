---
id: ldap
---

# LDAP Synchronization

Pocket ID can sync users and groups from an LDAP Source (lldap, OpenLDAP, Active Directory, etc.).

### LDAP Sync

- The LDAP Service will sync on Pocket ID startup and every hour once enabled from the Web UI.
- Users or groups synced from LDAP can **NOT** be edited from the Pocket ID Web UI.

### Generic LDAP Setup

1. Follow the installation guide [here](/pocket-id/setup/installation).
2. Once you have signed in with the initial admin account, navigate to the Application Configuration section at `https://pocket.id/settings/admin/application-configuration`.
3. Client Configuration Setup

| LDAP Variable      | Example Value                      | Description                                                   |
| ------------------ | ---------------------------------- | ------------------------------------------------------------- |
| LDAP URL           | ldaps://ldap.mydomain.com:636      | The URL with port to connect to LDAP                          |
| LDAP Bind DN       | cn=admin,ou=users,dc=domain,dc=com | The full DN value for the user with search privileges in LDAP |
| LDAP Bind Password | securepassword                     | The password for the Bind DN account                          |
| LDAP Search Base   | dc=domain,dc=com                   | The top-level path to search for users and groups             |

<br />

4. LDAP Attribute Configuration Setup

| LDAP Variable                     | Example Value      | Description                                                                      |
| --------------------------------- | ------------------ | -------------------------------------------------------------------------------- |
| User Unique Identifier Attribute  | uuid               | The LDAP attribute to uniquely identify the user, **this should never change**   |
| Username Attribute                | uid                | The LDAP attribute to use as the username of users                               |
| User Mail Attribute               | mail               | The LDAP attribute to use for the email of users                                 |
| User First Name Attribute         | givenName          | The LDAP attribute to use for the first name of users                            |
| User Last Name Attribute          | sn                 | The LDAP attribute to use for the last name of users                             |
| Group Unique Identifier Attribute | uuid               | The LDAP attribute to uniquely identify the groups, **this should never change** |
| Group Name Attribute              | uid                | The LDAP attribute to use as the name of synced groups                           |
| Admin Group Name                  | \_pocket_id_admins | The group name to use for admin permissions for LDAP users                       |
