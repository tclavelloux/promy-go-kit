# Changelog

## [0.2.0](https://github.com/tclavelloux/promy-template-go/compare/v0.1.1...v0.2.0) (2026-11-19)


### ⚠ BREAKING CHANGES

* **deploy:** ship .railway/railway.ts and a migrate-database command, migrate at pre-deploy ([#82](https://github.com/tclavelloux/promy-template-go/issues/82))

### Features

* **deploy:** ship .railway/railway.ts and a migrate-database command, migrate at pre-deploy ([#82](https://github.com/tclavelloux/promy-template-go/issues/82)) ([68bf70f](https://github.com/tclavelloux/promy-template-go/commit/68bf70f8b1afbf417cb0e7d7392187c3a5ec62a5))

## [0.1.1](https://github.com/tclavelloux/promy-template-go/compare/v0.1.0...v0.1.1) (2026-11-11)


### Bug Fixes

* add the missing main.go entrypoint ([ecb0611](https://github.com/tclavelloux/promy-template-go/commit/ecb0611149d3d6f0e656b8d592bc414efc33c036))
* add the missing main.go entrypoint ([#17](https://github.com/tclavelloux/promy-template-go/issues/17)) ([96f4d40](https://github.com/tclavelloux/promy-template-go/commit/96f4d40afd9325c052528538acbd2fee7b1420af))
* **deps:** bump otel and grpc to clear GO-2026-6505 and GO-2026-6348 ([#34](https://github.com/tclavelloux/promy-template-go/issues/34)) ([16efcc8](https://github.com/tclavelloux/promy-template-go/commit/16efcc8e8e9df5a7a2c9988a628d4856d4ddd182))
* **deps:** patch 7 Go vulnerabilities ([#56](https://github.com/tclavelloux/promy-template-go/issues/56)) ([79fae59](https://github.com/tclavelloux/promy-template-go/commit/79fae599338a1dc0c08275f6df7ea12cbca7226a))
* **deps:** patch GO-2026-6809 ([#61](https://github.com/tclavelloux/promy-template-go/issues/61)) ([82c377d](https://github.com/tclavelloux/promy-template-go/commit/82c377d3c00f3202073e60ca78d29496f9c0c73d))
* **deps:** patch GO-2026-6859 ([#78](https://github.com/tclavelloux/promy-template-go/issues/78)) ([d250e5b](https://github.com/tclavelloux/promy-template-go/commit/d250e5b0f8a91c199b91838f81682616cc2a0922))
* **deps:** remediate called vulnerabilities in the scaffold ([644b851](https://github.com/tclavelloux/promy-template-go/commit/644b8515492442993984516534c9325e7e61ab57))
* **deps:** remediate called vulnerabilities in the scaffold ([#14](https://github.com/tclavelloux/promy-template-go/issues/14)) ([d9cca05](https://github.com/tclavelloux/promy-template-go/commit/d9cca051d3c37f77ba5de534ae714d9a708c7a79))
* **scaffold:** repair the rename path and drop the stale go.mod.template ([#19](https://github.com/tclavelloux/promy-template-go/issues/19)) ([801dc18](https://github.com/tclavelloux/promy-template-go/commit/801dc1820c355ed12ece0ea6b932bcd0cc43b0ad))
* **scaffold:** rewrite the module path on rename ([830b067](https://github.com/tclavelloux/promy-template-go/commit/830b067ab90ff4f4f6f1159889804b3c437c636a))


### Build System

* **deps:** bump actions/checkout from 7.0.0 to 7.0.1 in the github-actions group ([#49](https://github.com/tclavelloux/promy-template-go/issues/49)) ([83024d2](https://github.com/tclavelloux/promy-template-go/commit/83024d2ad9df0684c0156c52c3ee1ac0c5742f56))
* **deps:** bump actions/upload-artifact from 7.0.1 to 7.0.2 in the github-actions group across 1 directory ([#52](https://github.com/tclavelloux/promy-template-go/issues/52)) ([6288ab8](https://github.com/tclavelloux/promy-template-go/commit/6288ab807f20774b9de58cdd2907d3361b75fc7a))
* **deps:** bump github.com/go-playground/validator/v10 from 10.30.4 to 10.30.5 in the gomod group across 1 directory ([#50](https://github.com/tclavelloux/promy-template-go/issues/50)) ([2a077bd](https://github.com/tclavelloux/promy-template-go/commit/2a077bd6259f50ad7340ca845b571c52692eeb5b))
* **deps:** bump github.com/labstack/echo/v4 from 4.16.0 to 4.16.1 in the gomod-security group across 1 directory ([#64](https://github.com/tclavelloux/promy-template-go/issues/64)) ([0d70e71](https://github.com/tclavelloux/promy-template-go/commit/0d70e7149ae49c7bf346796eadba8448a6c3c538))
* **deps:** bump github.com/magiconair/properties from 1.8.7 to 1.18.12 ([#45](https://github.com/tclavelloux/promy-template-go/issues/45)) ([26cf46c](https://github.com/tclavelloux/promy-template-go/commit/26cf46c5661304aa8e79edbc58d8c1026674cc17))
* **deps:** bump gitleaks/gitleaks-action from 2.3.9 to 3.0.0 ([#25](https://github.com/tclavelloux/promy-template-go/issues/25)) ([55613d8](https://github.com/tclavelloux/promy-template-go/commit/55613d86f5e44168810521a009ef907da48e5217))
* **deps:** bump googleapis/release-please-action from 4.4.1 to 5.0.0 ([#33](https://github.com/tclavelloux/promy-template-go/issues/33)) ([3c861ec](https://github.com/tclavelloux/promy-template-go/commit/3c861ec7aab40bb2ca985ae897904ed5fd0c2237))
* **deps:** bump the gomod group across 1 directory with 25 updates ([#58](https://github.com/tclavelloux/promy-template-go/issues/58)) ([6d9d388](https://github.com/tclavelloux/promy-template-go/commit/6d9d388e0a6b2c10d422931f1e969b61473c454a))
* **deps:** bump the gomod group across 1 directory with 6 updates ([#63](https://github.com/tclavelloux/promy-template-go/issues/63)) ([a986b3e](https://github.com/tclavelloux/promy-template-go/commit/a986b3ef25de4ed0f836a540b63409c5b27b9c7c))
* **deps:** bump the gomod group with 4 updates ([#60](https://github.com/tclavelloux/promy-template-go/issues/60)) ([ae17a74](https://github.com/tclavelloux/promy-template-go/commit/ae17a74f8072b2aa197c93808ffc7e13c10cd55c))
* **deps:** bump the gomod group with 8 updates ([#46](https://github.com/tclavelloux/promy-template-go/issues/46)) ([1b821a0](https://github.com/tclavelloux/promy-template-go/commit/1b821a0fbae7f6e706ea5df64df3b9214beb2c38))
* **deps:** bump the gomod group with 9 updates ([#71](https://github.com/tclavelloux/promy-template-go/issues/71)) ([a0e74fa](https://github.com/tclavelloux/promy-template-go/commit/a0e74fa0fa69735686b4001e9b9182d8931a174b))
* **deps:** downgrade checkout and validator by one patch to exercise the automerge gates ([#48](https://github.com/tclavelloux/promy-template-go/issues/48)) ([716106e](https://github.com/tclavelloux/promy-template-go/commit/716106e4687769712c76252f288ef639d6736ca6))
* **deps:** follow governance hooks v1.9.0 (go-test-coverage 2.21.0) ([#70](https://github.com/tclavelloux/promy-template-go/issues/70)) ([240f19e](https://github.com/tclavelloux/promy-template-go/commit/240f19e0fc4bf00d14e0371293613b9b505212b1))
* **hooks:** ship a working pre-commit stack with a pinned linter ([#9](https://github.com/tclavelloux/promy-template-go/issues/9)) ([821cf5c](https://github.com/tclavelloux/promy-template-go/commit/821cf5cb2f6f4f0e8983e2acdda9bcf7018b098e))
* pin go-test-coverage to the version CI runs ([#39](https://github.com/tclavelloux/promy-template-go/issues/39)) ([efa05c6](https://github.com/tclavelloux/promy-template-go/commit/efa05c6bae3a7c2c56e657934fd2bd6ec9ba8628))


### Continuous Integration

* add the workflow-sync fleet check ([#41](https://github.com/tclavelloux/promy-template-go/issues/41)) ([4134478](https://github.com/tclavelloux/promy-template-go/commit/4134478bcb17b31ba610c12d4f00fce3eeade172))
* automate dependency bumps and vuln fixes ([#43](https://github.com/tclavelloux/promy-template-go/issues/43)) ([d59c4d5](https://github.com/tclavelloux/promy-template-go/commit/d59c4d52a0ad26e921c82b34d56f905f3f6c9568))
* **coverage:** adopt go-coverage.yml v2 ([6941cb4](https://github.com/tclavelloux/promy-template-go/commit/6941cb4580b11aef7ef866ad4f3043c4c16dfb61))
* **coverage:** adopt go-coverage.yml v2 ([#11](https://github.com/tclavelloux/promy-template-go/issues/11)) ([4ac7c0c](https://github.com/tclavelloux/promy-template-go/commit/4ac7c0c9898e5f868c544f584467880fb3505cd8))
* **deps:** bump golangci-lint to 2.15.0, hooks to v1.12.0 and migrate test job to go-test ([#74](https://github.com/tclavelloux/promy-template-go/issues/74)) ([aa213b3](https://github.com/tclavelloux/promy-template-go/commit/aa213b3ab0133cd0dc45814b9738966227c01bd8))
* **deps:** bump governance hooks to hooks/v1.14.0 ([#77](https://github.com/tclavelloux/promy-template-go/issues/77)) ([54e83b1](https://github.com/tclavelloux/promy-template-go/commit/54e83b1be259d6a77e3d088af0a35237931579ea))
* **governance:** adopt hooks/v1.6.0 for release-please and fleet-bot automerge ([#62](https://github.com/tclavelloux/promy-template-go/issues/62)) ([ebf1f77](https://github.com/tclavelloux/promy-template-go/commit/ebf1f773721d356af134bf9fdb567d3b72e48db1))
* **lint:** adopt the reusable go-lint workflow ([#10](https://github.com/tclavelloux/promy-template-go/issues/10)) ([90f016d](https://github.com/tclavelloux/promy-template-go/commit/90f016d0c976867604c2671af3c63df2cfcd2119))
* pin actions by SHA, add Dependabot and the PR title check ([#23](https://github.com/tclavelloux/promy-template-go/issues/23)) ([8862679](https://github.com/tclavelloux/promy-template-go/commit/8862679f8c9c053ca184124f84af47bc6f7207fd))
* pin release-please-action by SHA ([#31](https://github.com/tclavelloux/promy-template-go/issues/31)) ([3208f51](https://github.com/tclavelloux/promy-template-go/commit/3208f51c6d6a1db21dcca8f25513ae339d086651))
* replace the build job with the docker gate and tighten CI hygiene ([#21](https://github.com/tclavelloux/promy-template-go/issues/21)) ([93e1948](https://github.com/tclavelloux/promy-template-go/commit/93e1948d5a10a5f493f18eb41456234ab069a54f))
* run tests through the reusable go-test workflow ([#73](https://github.com/tclavelloux/promy-template-go/issues/73)) ([12490a4](https://github.com/tclavelloux/promy-template-go/commit/12490a4f9005a9015e0a0c3cefbf01aff66ab934))
* **vuln-fix:** pass the GitHub App secrets to go-vuln-fix ([#57](https://github.com/tclavelloux/promy-template-go/issues/57)) ([875634f](https://github.com/tclavelloux/promy-template-go/commit/875634f379608f89ff609c6b76953452a8ded673))


### Documentation

* align scaffold documentation and setup script with pre-commit ([#37](https://github.com/tclavelloux/promy-template-go/issues/37)) ([da72c72](https://github.com/tclavelloux/promy-template-go/commit/da72c723e5df8b266e116104c88856e175964e59))
* refresh HOWTO commit conventions and the README watermark ([#40](https://github.com/tclavelloux/promy-template-go/issues/40)) ([5e3fe7d](https://github.com/tclavelloux/promy-template-go/commit/5e3fe7dc80801a08e88aab375390bcb4ac2d77e2))
* remove the release guide, drop a stale githooks comment, use read-only Actions permissions ([#38](https://github.com/tclavelloux/promy-template-go/issues/38)) ([e2565e9](https://github.com/tclavelloux/promy-template-go/commit/e2565e95bdc6775add714f7365edd2930fed1fee))
* state that a PR is one business unit ([#42](https://github.com/tclavelloux/promy-template-go/issues/42)) ([0979b5a](https://github.com/tclavelloux/promy-template-go/commit/0979b5aaff4b66695c916a66318bad2bfdd117d2))


### Miscellaneous

* **claude:** enable promy plugin ([#69](https://github.com/tclavelloux/promy-template-go/issues/69)) ([3ffad94](https://github.com/tclavelloux/promy-template-go/commit/3ffad94bde75396b4c2afd3d5cfae2f04df24c24))
* drop the stale go.mod.template ([2f42351](https://github.com/tclavelloux/promy-template-go/commit/2f4235132f07a84ae0d6023b11bf199db7ccffdb))
* **gitignore:** ignore coverage artifacts ([7ee7e86](https://github.com/tclavelloux/promy-template-go/commit/7ee7e86012ca0dcbe9d0509449abe24204164dba))
* **hooks:** bump governance hooks to v1.7.0 and go-test-coverage to v2.20.1 ([#67](https://github.com/tclavelloux/promy-template-go/issues/67)) ([5696184](https://github.com/tclavelloux/promy-template-go/commit/5696184bee7c0a06ce992be485bdc55ebb5979ef))
* **hooks:** run each pre-commit hook once per commit ([#66](https://github.com/tclavelloux/promy-template-go/issues/66)) ([891fe55](https://github.com/tclavelloux/promy-template-go/commit/891fe558e862befd35fcaf7c5e2f8adadba168c3))
* migrate from .githooks to pre-commit ([#32](https://github.com/tclavelloux/promy-template-go/issues/32)) ([74a5fec](https://github.com/tclavelloux/promy-template-go/commit/74a5fec7d126577b0444a90fe056bd4da1178779))
* **pre-commit:** bump promy-github-workflows hooks to v1.3.0 ([#59](https://github.com/tclavelloux/promy-template-go/issues/59)) ([f3e5cc5](https://github.com/tclavelloux/promy-template-go/commit/f3e5cc50438eb20ad738af07c512a3a5707cfb0d))
* remove the unused release_please directory ([#51](https://github.com/tclavelloux/promy-template-go/issues/51)) ([13c34b4](https://github.com/tclavelloux/promy-template-go/commit/13c34b4b3d21c5990e0df7f1837eb7acd22f12f7))
