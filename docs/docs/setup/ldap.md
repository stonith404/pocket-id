---
id: ldap
---

# LDAP Synchronization

Pocket ID can sync Users and Group from a LDAP Source (lldap, OpenLDAP, Active Directory, etc.)

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
