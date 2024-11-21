export type AppConfig = {
	appName: string;
	allowOwnAccountEdit: boolean;
};

export type AllAppConfig = AppConfig & {
	sessionDuration: number;
	emailsVerified: boolean;
	emailEnabled: boolean;
	smtpHost: string;
	smtpPort: number;
	smtpFrom: string;
	smtpUser: string;
	smtpPassword: string;
	smtpSkipCertVerify: boolean;
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
