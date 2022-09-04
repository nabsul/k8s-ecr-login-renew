# Changelog

## V1.7.2 (2022-09-10) - Helm Charts!

- Added support for deploying via Helm chart. Thanks to:
  - [devec0](https://github.com/devec0): v1.7.1
  - [PawelLipski](https://github.com/PawelLipski): v1.7.1
  - [xavidop](https://github.com/xavidop): v1.7.1
  - [armenr](https://github.com/armenr): v1.7.1
 

## V1.7.1 (2022-06-05) - V1.7.0 for Real

I forgot to merge the change that stops using root user in the container. Thanks to [PawelLipski](https://github.com/PawelLipski) for spotting this.

## V1.7.0 (2022-06-05) - Security and Updates

- The job now runs as a user in the container instead of root (#30)
- Updated to latest version of Go
- Updated dependencies to latest versions

## v1.6 (2021-10-31) - Spooky Separators!

- Support multi-line and whitespace in Namespace list
  - Contributed by [Jeremy Ruffell](https://github.com/jeremyruffell) in [pull request 24](https://github.com/nabsul/k8s-ecr-login-renew/pull/24)
- CI and Docker build improvements:
  - Fix broken GitHub CI actions
  - Combine amd64 and ARM builds into one using `buildx`
  - Reduce size of container image
  - Contributed by [Jeremy Ruffell](https://github.com/jeremyruffell) in [pull request 25](https://github.com/nabsul/k8s-ecr-login-renew/pull/24)
- Minor bug fix in an error message
- Fix some broken integration tests
- Implemented in [pull request 26](https://github.com/nabsul/k8s-ecr-login-renew/pull/26) 

## v1.5 (2021-04-03)

- Allow custom/multiple registry URLs with a new environment variable: `DOCKER_REGISTRIES`
  - Contributed by [Veraticus](https://github.com/Veraticus) with feedback from [PawelLipski](https://github.com/PawelLipski) in [pull request 18](https://github.com/nabsul/k8s-ecr-login-renew/pull/18)
- Implemented in [pull request 19](https://github.com/nabsul/k8s-ecr-login-renew/pull/19)

## v1.4 (2021-02-13)

- Update Docker secrets instead of delete+create
  - Suggested by [xavidop](https://github.com/xavidop) in [issue 15](https://github.com/nabsul/k8s-ecr-login-renew/issues/15)
- Fall back to old delete+create if update fails to avoid breaking old users

## v1.3 (2020-06-07)

- Added support for ARM
  - Tested by [kuskoman](https://github.com/kuskoman) on Raspberry Pi

## v1.2 (2020-06-07)

- `TARGET_NAMESPACE` now supports multiple namespaces and wildcards
  - Suggested by [Q-Nimbus](https://github.com/Q-Nimbus) in [issue 5](https://github.com/nabsul/k8s-ecr-login-renew/issues/5)
- Added automated tests
- Various refactoring and code restructuring

## v1.1 (2020-04-25)

- Added environment variable to specify a namespace
  - Contribution from [YoSmudge](https://github.com/YoSmudge) in [pull request 1](https://github.com/nabsul/k8s-ecr-login-renew/pull/1)
- Add a changelog
- Add a contributor list

## v1.0 (2020-03-22)

Initial version release.
More info here: https://nabeel.dev/2020/03/22/k8s-ecr-login-renew/
