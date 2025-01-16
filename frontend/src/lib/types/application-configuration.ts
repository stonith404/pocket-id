export type AppConfig = {
	appName: string;
	allowOwnAccountEdit: boolean;
};

export type AllAppConfig = AppConfig & {
	// General
	sessionDuration: number;
	emailsVerified: boolean;
	// Email
	emailEnabled: boolean;
	smtpHost: string;
	smtpPort: number;
	smtpFrom: string;
	smtpUser: string;
	smtpPassword: string;
	smtpTls: boolean;
	smtpSkipCertVerify: boolean;
	// LDAP
	ldapEnabled: boolean;
	ldapUrl: string;
	ldapBindDn: string;
	ldapBindPassword: string;
	ldapBase: string;
	ldapSkipCertVerify: boolean;
	ldapAttributeUserUniqueIdentifier: string;
	ldapAttributeUserUsername: string;
	ldapAttributeUserEmail: string;
	ldapAttributeUserFirstName: string;
	ldapAttributeUserLastName: string;
	ldapAttributeGroupUniqueIdentifier: string;
	ldapAttributeGroupName: string;
};

export type AppConfigRawResponse = {
	key: string;
	type: string;
	value: string;
}[];

export type AppVersionInformation = {
	isUpToDate: boolean;
	newestVersion: string;
	currentVersion: string;
};
