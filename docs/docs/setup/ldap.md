---
id: ldap
---

# LDAP Synchronization

Pocket ID can sync Users and Group from a LDAP Source (lldap, OpenLDAP, Active Directory, etc.)

### LDAP Sync

- The LDAP Service Will sync on Pocket ID Startup and Every Hour once Enabled from the Web UI
- Users or Groups Synced from LDAP can **NOT** be Edited from the Pocket ID Web UI.

### Generic LDAP Setup

1. Follow the Installation guide [Here](/pocket-id/setup/installation)
2. Once you have Signed in with the initial admin account navigate to the Application Configuration section at `https://pocket.id/settings/admin/application-configuration`
3. Client Configuration Setup

| LDAP Variable                     | Example Value                         | Description                                                            |
| ---------------------------- | ------------------------------------- | --------------------------------------------------------------------------- |
| LDAP URL                     | ldaps://ldap.mydomain.com:636         | The URL with Port to Connect to LDAP                                        |
| LDAP Bind DN                 | cn=admin,ou=users,dc=domain,dc=com    | The full DN Value for the User to Search Privileges in LDAP                 |
| LDAP Bind Password           | securepassword                        | The Password for the Bind DN Account                                        |
| LDAP Search Base             | dc=domain,dc=com                      | The Top Level Path to search for Users and Groups                           |

<br />

4. LDAP Attribute Configuration Setup

| LDAP Variable                     | Example Value                         | Description                                                            |
| ----------------------------      | ------------------------------------- | ------------------------------------------------------------------------------ |
| User Unique Identifier Attribute  | uuid                                  | The LDAP Attribute to Unique Identify the User, **This Should never Change**   |
| Username Attribute                | uid                                   | The LDAP Attribute to use as the Username of Users                             |
| User Mail Attribute               | mail                                  | The LDAP Attribute to Use for the Email for Users                              |
| User First Name Attribute         | givenName                             | The LDAP Attribute to Use for the First Name for Users                         |
| User Last Name Attribute          | sn                                    | The LDAP Attribute to Use for the Last Name for Users                          |
| Group Unique Identifier Attribute | uuid                                  | The LDAP Attribute to Unique Identify the Groups, **This Should never Change** |
| Group Name Attribute              | uid                                   | The LDAP Attribute to use as the Name of Synced Groups                         |
| Admin Group Name                  | _pocket_id_admins                     | The Group Name to Use for Admin Permissions for LDAP Users                     |
