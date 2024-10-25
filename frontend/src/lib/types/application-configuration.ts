export type AppConfig = {
	appName: string;
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
