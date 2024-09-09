export type AllAppConfig = {
	appName: string;
	sessionDuration: string;
	emailEnabled: string;
	smtpHost: string;
	smtpPort: string;
	smtpFrom: string;
	smtpUser: string;
	smtpPassword: string;
};

export type AppConfig = AllAppConfig;

export type AppConfigRawResponse = {
	key: string;
	type: string;
	value: string;
}[];
