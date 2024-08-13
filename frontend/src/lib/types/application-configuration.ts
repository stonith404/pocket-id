
export type AllApplicationConfiguration = {
	appName: string;
    sessionDuration: string;
};

export type ApplicationConfiguration = AllApplicationConfiguration;

export type ApplicationConfigurationRawResponse = {
	key: string;
    type: string;
    value: string;
}[];