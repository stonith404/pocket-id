## [](https://github.com/stonith404/pocket-id/compare/v0.26.0...v) (2025-01-22)


### Features

* display private IP ranges correctly in audit log ([#139](https://github.com/stonith404/pocket-id/issues/139)) ([72923bb](https://github.com/stonith404/pocket-id/commit/72923bb86dc5d07d56aea98cf03320667944b553))


### Bug Fixes

* add save changes dialog before sending test email ([#165](https://github.com/stonith404/pocket-id/issues/165)) ([d02f475](https://github.com/stonith404/pocket-id/commit/d02f4753f3fbda75cd415ebbfe0702765c38c144))
* ensure the downloaded GeoLite2 DB is not corrupted & prevent RW race condition ([#138](https://github.com/stonith404/pocket-id/issues/138)) ([f7710f2](https://github.com/stonith404/pocket-id/commit/f7710f298898d322885c1c83680e26faaa0bb800))

## [](https://github.com/stonith404/pocket-id/compare/v0.25.1...v) (2025-01-20)


### Features

* support wildcard callback URLs ([8a1db0c](https://github.com/stonith404/pocket-id/commit/8a1db0cb4a5d4b32b4fdc19d41fff688a7c71a56))


### Bug Fixes

* non LDAP users get created with a empty LDAP ID string ([3f02d08](https://github.com/stonith404/pocket-id/commit/3f02d081098ad2caaa60a56eea4705639f80d01f))

## [](https://github.com/stonith404/pocket-id/compare/v0.25.0...v) (2025-01-19)


### Bug Fixes

* disable account details inputs if user is imported from LDAP ([a8b9d60](https://github.com/stonith404/pocket-id/commit/a8b9d60a86e80c10d6fba07072b1d32cec400ecb))

## [](https://github.com/stonith404/pocket-id/compare/v0.24.1...v) (2025-01-19)


### Features

* add LDAP sync ([#106](https://github.com/stonith404/pocket-id/issues/106)) ([5101b14](https://github.com/stonith404/pocket-id/commit/5101b14eec68a9507e1730994178d0ebe8185876))
* allow sign in with email ([#100](https://github.com/stonith404/pocket-id/issues/100)) ([06b90ed](https://github.com/stonith404/pocket-id/commit/06b90eddd645cce57813f2536e4a6a8836548f2b))
* automatically authorize client if signed in ([d5dd118](https://github.com/stonith404/pocket-id/commit/d5dd118a3f4ad6eed9ca496c458201bb10f148a0))


### Bug Fixes

* always set secure on cookie ([#130](https://github.com/stonith404/pocket-id/issues/130)) ([fda08ac](https://github.com/stonith404/pocket-id/commit/fda08ac1cd88842e25dc47395ed1288a5cfac4f8))
* don't panic if LDAP sync fails on startup ([e284e35](https://github.com/stonith404/pocket-id/commit/e284e352e2b95fac1d098de3d404e8531de4b869))
* improve spacing of checkboxes on application configuration page ([090eca2](https://github.com/stonith404/pocket-id/commit/090eca202d198852e6fbf4e6bebaf3b5ada13944))
* search input not displayed if response hasn't any items ([05a98eb](https://github.com/stonith404/pocket-id/commit/05a98ebe87d7a88e8b96b144c53250a40d724ec3))
* session duration ignored in cookie expiration ([bc8f454](https://github.com/stonith404/pocket-id/commit/bc8f454ea173ecc60e06450a1d22e24207f76714))

## [](https://github.com/stonith404/pocket-id/compare/v0.24.0...v) (2025-01-13)


### Bug Fixes

* audit log table overflow if row data is long ([4d337a2](https://github.com/stonith404/pocket-id/commit/4d337a20c5cb92ef80bb7402f9b99b08e3ad0b6b))
* optional arguments not working with `create-one-time-access-token.sh` ([8885571](https://github.com/stonith404/pocket-id/commit/888557171d61589211b10f70dce405126216ad61))
* remove restrictive validation for group names ([be6e25a](https://github.com/stonith404/pocket-id/commit/be6e25a167de8bf07075b46f09d9fc1fa6c74426))

## [](https://github.com/stonith404/pocket-id/compare/v0.23.0...v) (2025-01-11)


### Features

* add sorting for tables ([fd69830](https://github.com/stonith404/pocket-id/commit/fd69830c2681985e4fd3c5336a2b75c9fb7bc5d4))


### Bug Fixes

* pkce state not correctly reflected in oidc client info ([61d18a9](https://github.com/stonith404/pocket-id/commit/61d18a9d1b167ff59a59523ff00d00ca8f23258d))
* send test email to the user that has requested it ([a649c4b](https://github.com/stonith404/pocket-id/commit/a649c4b4a543286123f4d1f3c411fe1a7e2c6d71))

## [](https://github.com/stonith404/pocket-id/compare/v0.22.0...v) (2025-01-03)


### Features

* add PKCE for non public clients ([adcf3dd](https://github.com/stonith404/pocket-id/commit/adcf3ddc6682794e136a454ef9e69ddd130626a8))
* use same table component for OIDC client list as all other lists ([2d31fc2](https://github.com/stonith404/pocket-id/commit/2d31fc2cc9201bb93d296faae622f52c6dcdfebc))

## [](https://github.com/stonith404/pocket-id/compare/v0.21.0...v) (2025-01-01)


### Features

* add warning if passkeys missing ([2d0bd8d](https://github.com/stonith404/pocket-id/commit/2d0bd8dcbfb73650b7829cb66f40decb284bd73b))


### Bug Fixes

* allow first and last name of user to be between 1 and 50 characters ([1ff20ca](https://github.com/stonith404/pocket-id/commit/1ff20caa3ccd651f9fb30f958ffb807dfbbcbd8a))
* hash in callback url is incorrectly appended ([f6f2736](https://github.com/stonith404/pocket-id/commit/f6f2736bba65eee017f2d8cdaa70621574092869))
* make user validation consistent between pages ([333a1a1](https://github.com/stonith404/pocket-id/commit/333a1a18d59f675111f4ed106fa5614ef563c6f4))
* passkey can't be added if `PUBLIC_APP_URL` includes a port ([0729ce9](https://github.com/stonith404/pocket-id/commit/0729ce9e1a8dab9912900a01dcd0fbd892718a1a))

## [](https://github.com/stonith404/pocket-id/compare/v0.20.1...v) (2024-12-17)


### Features

* improve error state design for login page ([0716c38](https://github.com/stonith404/pocket-id/commit/0716c38fb8ce7fa719c7fe0df750bdb213786c21))


### Bug Fixes

* OIDC client logo gets removed if other properties get updated ([789d939](https://github.com/stonith404/pocket-id/commit/789d9394a533831e7e2fb8dc3f6b338787336ad8))

## [](https://github.com/stonith404/pocket-id/compare/v0.20.0...v) (2024-12-13)


### Bug Fixes

* `create-one-time-access-token.sh` script not compatible with postgres ([34e3519](https://github.com/stonith404/pocket-id/commit/34e35193f9f3813f6248e60f15080d753e8da7ae))
* wrong date time datatype used for read operations with Postgres ([bad901e](https://github.com/stonith404/pocket-id/commit/bad901ea2b661aadd286e5e4bed317e73bd8a70d))

## [](https://github.com/stonith404/pocket-id/compare/v0.19.0...v) (2024-12-12)


### Features

* add support for Postgres database provider ([#79](https://github.com/stonith404/pocket-id/issues/79)) ([9d20a98](https://github.com/stonith404/pocket-id/commit/9d20a98dbbc322fa6f0644e8b31e6b97769887ce))

## [](https://github.com/stonith404/pocket-id/compare/v0.18.0...v) (2024-11-29)


### Features

* **geolite:** add Tailscale IP detection with CGNAT range check  ([#77](https://github.com/stonith404/pocket-id/issues/77)) ([edce3d3](https://github.com/stonith404/pocket-id/commit/edce3d337129c9c6e8a60df2122745984ba0f3e0))

## [](https://github.com/stonith404/pocket-id/compare/v0.17.0...v) (2024-11-28)


### Features

* add option to disable TLS for email sending ([f9fa2c6](https://github.com/stonith404/pocket-id/commit/f9fa2c6706a8bf949fe5efd6664dec8c80e18659))
* allow empty user and password in SMTP configuration ([a9f4dad](https://github.com/stonith404/pocket-id/commit/a9f4dada321841d3611b15775307228b34e7793f))


### Bug Fixes

* email save toast shows two times ([f2bfc73](https://github.com/stonith404/pocket-id/commit/f2bfc731585ad7424eb8c4c41c18368fc0f75ffc))

## [](https://github.com/stonith404/pocket-id/compare/v0.16.0...v) (2024-11-26)


### ⚠ BREAKING CHANGES

* add option to specify the Max Mind license key for the Geolite2 db

### Features

* add option to specify the Max Mind license key for the Geolite2 db ([fcf08a4](https://github.com/stonith404/pocket-id/commit/fcf08a4d898160426442bd80830f4431988f4313))


### Bug Fixes

* don't try to create a new user if the Docker user is not root ([#71](https://github.com/stonith404/pocket-id/issues/71)) ([0e95e9c](https://github.com/stonith404/pocket-id/commit/0e95e9c56f4c3f84982f508fdb6894ba747952b4))

## [](https://github.com/stonith404/pocket-id/compare/v0.15.0...v) (2024-11-24)


### Features

* add health check ([058084e](https://github.com/stonith404/pocket-id/commit/058084ed64816b12108e25bf04af988fc97772ed))
* improve error message for invalid callback url ([f637a89](https://github.com/stonith404/pocket-id/commit/f637a89f579aefb8dc3c3c16a27ef0bc453dfe40))

## [](https://github.com/stonith404/pocket-id/compare/v0.14.0...v) (2024-11-21)


### Features

* add option to skip TLS certificate check and ability to send test email ([653d948](https://github.com/stonith404/pocket-id/commit/653d948f73b61e6d1fd3484398fef1a2a37e6d92))
* add PKCE support ([3613ac2](https://github.com/stonith404/pocket-id/commit/3613ac261cf65a2db0620ff16dc6df239f6e5ecd))


### Bug Fixes

* mobile layout overflow on application configuration page ([e784093](https://github.com/stonith404/pocket-id/commit/e784093342f9977ea08cac65ff0c3de4d2644872))

## [](https://github.com/stonith404/pocket-id/compare/v0.13.1...v) (2024-11-11)


### Features

* add audit log event for one time access token sign in ([aca2240](https://github.com/stonith404/pocket-id/commit/aca2240a50a12e849cfb6e1aa56390b000aebae0))


### Bug Fixes

* overflow of pagination control on mobile ([de45398](https://github.com/stonith404/pocket-id/commit/de4539890349153c467013c24c4d6b30feb8fed8))
* time displayed incorrectly in audit log ([3d3fb4d](https://github.com/stonith404/pocket-id/commit/3d3fb4d855ef510f2292e98fcaaaf83debb5d3e0))

## [](https://github.com/stonith404/pocket-id/compare/v0.13.0...v) (2024-11-01)


### Features

* add list empty indicator ([becfc00](https://github.com/stonith404/pocket-id/commit/becfc0004a87c01e18eb92ac85bf4e33f105b6a3))


### Bug Fixes

* errors in middleware do not abort the request ([376d747](https://github.com/stonith404/pocket-id/commit/376d747616b1e835f252d20832c5ae42b8b0b737))
* typo in Self-Account Editing description ([5b9f4d7](https://github.com/stonith404/pocket-id/commit/5b9f4d732615f428c13d3317da96a86c5daebd89))

## [](https://github.com/stonith404/pocket-id/compare/v0.12.0...v) (2024-10-31)


### Features

* add ability to define expiration of one time link ([2ccabf8](https://github.com/stonith404/pocket-id/commit/2ccabf835c2c923d6986d9cafb4e878f5110b91a))

## [](https://github.com/stonith404/pocket-id/compare/v0.11.0...v) (2024-10-28)


### Features

* add option to disable self-account editing ([8304065](https://github.com/stonith404/pocket-id/commit/83040656525cf7b6c8f2acf416c5f8f3288f3d48))
* add validation to custom claim input ([7bfc3f4](https://github.com/stonith404/pocket-id/commit/7bfc3f43a591287c038187ed5e782de6b9dd738b))
* custom claims ([#53](https://github.com/stonith404/pocket-id/issues/53)) ([c056089](https://github.com/stonith404/pocket-id/commit/c056089c6043a825aaaaecf0c57454892a108f1d))

## [](https://github.com/stonith404/pocket-id/compare/v0.10.0...v) (2024-10-25)


### Features

* add `email_verified` claim ([5565f60](https://github.com/stonith404/pocket-id/commit/5565f60d6d62ca24bedea337e21effc13e5853a5))


### Bug Fixes

* powered by link text color in light mode ([18c5103](https://github.com/stonith404/pocket-id/commit/18c5103c20ce79abdc0f724cdedd642c09269e78))

## [](https://github.com/stonith404/pocket-id/compare/v0.9.0...v) (2024-10-23)


### Features

* add script for creating one time access token ([a1985ce](https://github.com/stonith404/pocket-id/commit/a1985ce1b200550e91c5cb42a8d19899dcec831e))
* add version information to footer and update link if new update is available ([70ad0b4](https://github.com/stonith404/pocket-id/commit/70ad0b4f39699fd81ffdfd5c8d6839f49348be78))


### Bug Fixes

* cache version information for 3 hours ([29d632c](https://github.com/stonith404/pocket-id/commit/29d632c1514d6edacdfebe6deae4c95fc5a0f621))
* improve text for initial admin account setup ([0a07344](https://github.com/stonith404/pocket-id/commit/0a0734413943b1fff27d8f4ccf07587e207e2189))
* increase callback url count ([f3f0e1d](https://github.com/stonith404/pocket-id/commit/f3f0e1d56d7656bdabbd745a4eaf967f63193b6c))
* no DTO was returned from exchange one time access token endpoint ([824c5cb](https://github.com/stonith404/pocket-id/commit/824c5cb4f3d6be7f940c1758112fbe9322df5768))

## [](https://github.com/stonith404/pocket-id/compare/v0.8.1...v) (2024-10-18)


### Features

* add environment variable to change the caddy port in Docker ([ff06bf0](https://github.com/stonith404/pocket-id/commit/ff06bf0b34496ce472ba6d3ebd4ea249f21c0ec3))
* use improve table for users and audit logs ([11ed661](https://github.com/stonith404/pocket-id/commit/11ed661f86a512f78f66d604a10c1d47d39f2c39))


### Bug Fixes

* allow copy to clipboard for client secret ([29748cc](https://github.com/stonith404/pocket-id/commit/29748cc6c7b7e5a6b54bfe837e0b1a98fa1ad594))

## [](https://github.com/stonith404/pocket-id/compare/v0.8.0...v) (2024-10-11)


### Bug Fixes

* add key id to JWK ([282ff82](https://github.com/stonith404/pocket-id/commit/282ff82b0c7e2414b3528c8ca325758245b8ae61))

## [](https://github.com/stonith404/pocket-id/compare/v0.7.1...v) (2024-10-04)


### Features

* add location based on ip to the audit log ([025378d](https://github.com/stonith404/pocket-id/commit/025378d14edd2d72da76e90799a0ccdd42cf672c))

## [](https://github.com/stonith404/pocket-id/compare/v0.7.0...v) (2024-10-03)


### Bug Fixes

* initials don't get displayed if Gravatar avatar doesn't exist ([e095628](https://github.com/stonith404/pocket-id/commit/e09562824a794bc7d240e9d229709d4b389db7d5))

## [](https://github.com/stonith404/pocket-id/compare/v0.6.0...v) (2024-10-03)


### ⚠ BREAKING CHANGES

* add ability to set light and dark mode logo

### Features

* add ability to set light and dark mode logo ([be45eed](https://github.com/stonith404/pocket-id/commit/be45eed125e33e9930572660a034d5f12dc310ce))

## [](https://github.com/stonith404/pocket-id/compare/v0.5.3...v) (2024-10-02)


### Features

* add copy to clipboard option for OIDC client information ([f82020c](https://github.com/stonith404/pocket-id/commit/f82020ccfb0d4fbaa1dd98182188149d8085252a))
* add gravatar profile picture integration ([365734e](https://github.com/stonith404/pocket-id/commit/365734ec5d8966c2ab877c60cfb176b9cdc36880))
* add user groups ([24c948e](https://github.com/stonith404/pocket-id/commit/24c948e6a66f283866f6c8369c16fa6cbcfa626c))


### Bug Fixes

* only return user groups if it is explicitly requested ([a4a90a1](https://github.com/stonith404/pocket-id/commit/a4a90a16a9726569a22e42560184319b25fd7ca6))

## [](https://github.com/stonith404/pocket-id/compare/v0.5.2...v) (2024-09-26)


### Bug Fixes

* add space to "Firstname" and "Lastname" label ([#31](https://github.com/stonith404/pocket-id/issues/31)) ([d6a9bb4](https://github.com/stonith404/pocket-id/commit/d6a9bb4c09efb8102da172e49c36c070b341f0fc))
* port environment variables get ignored in caddyfile ([3c67765](https://github.com/stonith404/pocket-id/commit/3c67765992d7369a79812bc8cd216c9ba12fd96e))

## [](https://github.com/stonith404/pocket-id/compare/v0.5.1...v) (2024-09-19)


### Bug Fixes

* updated application name doesn't apply to webauthn credential ([924bb14](https://github.com/stonith404/pocket-id/commit/924bb1468bbd8e42fa6a530ef740be73ce3b3914))

## [](https://github.com/stonith404/pocket-id/compare/v0.5.0...v) (2024-09-16)


### Features

* **email:** improve email templating ([#27](https://github.com/stonith404/pocket-id/issues/27)) ([64cf562](https://github.com/stonith404/pocket-id/commit/64cf56276a07169bc601a11be905c1eea67c4750))


### Bug Fixes

* debounce oidc client and user search ([9c2848d](https://github.com/stonith404/pocket-id/commit/9c2848db1d93c230afc6c5f64e498e9f6df8c8a7))

## [](https://github.com/stonith404/pocket-id/compare/v0.4.1...v) (2024-09-09)


### Features

* add audit log with email notification ([#26](https://github.com/stonith404/pocket-id/issues/26)) ([9121239](https://github.com/stonith404/pocket-id/commit/9121239dd7c14a2107a984f9f94f54227489a63a))

## [](https://github.com/stonith404/pocket-id/compare/v0.4.0...v) (2024-09-06)


### Features

* add name claim to userinfo endpoint and id token ([4e7574a](https://github.com/stonith404/pocket-id/commit/4e7574a297307395603267c7a3285d538d4111d8))


### Bug Fixes

* limit width of content on large screens ([c6f83a5](https://github.com/stonith404/pocket-id/commit/c6f83a581ad385391d77fec7eeb385060742f097))
* show error message if error occurs while authorizing new client ([8038a11](https://github.com/stonith404/pocket-id/commit/8038a111dd7fa8f5d421b29c3bc0c11d865dc71b))

## [](https://github.com/stonith404/pocket-id/compare/v0.3.1...v) (2024-09-03)


### Features

* add setup details to oidc client details ([fd21ce5](https://github.com/stonith404/pocket-id/commit/fd21ce5aac1daeba04e4e7399a0720338ea710c2))
* add support for more username formats ([903b0b3](https://github.com/stonith404/pocket-id/commit/903b0b39181c208e9411ee61849d2671e7c56dc5))


### Bug Fixes

* non pointer passed to create user ([e7861df](https://github.com/stonith404/pocket-id/commit/e7861df95a6beecab359d1c56f4383373f74bb73))
* oidc client logo not displayed on authorize page ([28ed064](https://github.com/stonith404/pocket-id/commit/28ed064668afeec8f80adda59ba94f1fc2fbce17))
* typo in hasLogo property of oidc dto ([2b9413c](https://github.com/stonith404/pocket-id/commit/2b9413c7575e1322f8547490a9b02a1836bad549))

## [](https://github.com/stonith404/pocket-id/compare/v0.3.0...v) (2024-08-24)


### Bug Fixes

* empty lists don't get returned correctly from the api ([97f7fc4](https://github.com/stonith404/pocket-id/commit/97f7fc4e288c2bb49210072a7a151b58ef44f5b5))

## [](https://github.com/stonith404/pocket-id/compare/v0.2.1...v) (2024-08-23)


### Features

* add support for multiple callback urls ([8166e2e](https://github.com/stonith404/pocket-id/commit/8166e2ead7fc71a0b7a45950b05c5c65a60833b6))


### Bug Fixes

* db migration for multiple callback urls ([552d7cc](https://github.com/stonith404/pocket-id/commit/552d7ccfa58d7922ecb94bdfe6a86651b4cf2745))

## [](https://github.com/stonith404/pocket-id/compare/v0.2.0...v) (2024-08-19)


### Bug Fixes

* session duration can't be updated ([4780548](https://github.com/stonith404/pocket-id/commit/478054884389ed8a08d707fd82da7b31177a67e5))

## [](https://github.com/stonith404/pocket-id/compare/v0.1.3...v) (2024-08-19)


### Features

* add `INTERNAL_BACKEND_URL` env variable ([0595d73](https://github.com/stonith404/pocket-id/commit/0595d73ea5afbd7937b8f292ffe624139f818f41))
* add user info endpoint to support more oidc clients ([fdc1921](https://github.com/stonith404/pocket-id/commit/fdc1921f5dcb5ac6beef8d1c9b1b7c53f514cce5))
* change default logo ([9eec7a3](https://github.com/stonith404/pocket-id/commit/9eec7a3e9eb7f690099f38a5d4cf7c2516ea9ef9))

## [](https://github.com/stonith404/pocket-id/compare/v0.1.2...v) (2024-08-13)


### Bug Fixes

* add missing passkey flags to make icloud passkeys work ([cc407e1](https://github.com/stonith404/pocket-id/commit/cc407e17d409041ed88b959ce13bd581663d55c3))
* logo not white in dark mode ([5749d05](https://github.com/stonith404/pocket-id/commit/5749d0532fc38bf2fc66571878b7c71643895c9e))

## [](https://github.com/stonith404/pocket-id/compare/v0.1.1...v) (2024-08-13)


### Features

* add option to change session duration ([475b932](https://github.com/stonith404/pocket-id/commit/475b932f9d0ec029ada844072e9d89bebd4e902c))


### Bug Fixes

* a non admin user was able to make himself an admin ([df0cd38](https://github.com/stonith404/pocket-id/commit/df0cd38deeea516c47b26a080eed522f19f7290f))
* background image not loading ([7b44189](https://github.com/stonith404/pocket-id/commit/7b4418958ebfffffd216ef5ba7313cfaad9bc9fa))
* background image on mobile ([4a808c8](https://github.com/stonith404/pocket-id/commit/4a808c86ac204f9b58cfa02f5ceb064162a87076))
* disable search engine indexing ([8395492](https://github.com/stonith404/pocket-id/commit/83954926f5ee328ebf75a75bb47b380ec0680378))

## [](https://github.com/stonith404/pocket-id/compare/v0.1.0...v) (2024-08-12)


### Features

* add rounded corners to logo ([bec908f](https://github.com/stonith404/pocket-id/commit/bec908f9078aaa4eec03b730fc36b9fffb1ece74))


### Bug Fixes

* one time link not displayed correctly ([486771f](https://github.com/stonith404/pocket-id/commit/486771f433872d08164156d5d6fb0aeb5ae0d125))

##  (2024-08-12)

