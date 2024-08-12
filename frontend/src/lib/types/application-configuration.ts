
export type AllApplicationConfiguration = {
	appName: string;
};

export type ApplicationConfiguration = AllApplicationConfiguration;

export type ApplicationConfigurationRawResponse = {
	key: string;
    type: string;
    value: string;
}[];