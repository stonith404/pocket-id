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

