# Changelog

## v1.5 (2021-04-03)

- Handle the case of multiple registries in AWS (contributed by [Veraticus](https://github.com/Veraticus) in [pull request 18](https://github.com/nabsul/k8s-ecr-login-renew/pull/18))

## v1.4 (2021-02-13)

- Update Docker secrets instead of delete+create (suggested by [xavidop](https://github.com/xavidop) in [issue 15](https://github.com/nabsul/k8s-ecr-login-renew/issues/15))
- Fall back to old delete+create if update fails to avoid breaking old users

## v1.3 (2020-06-07)

- Added support for ARM (tested by [kuskoman](https://github.com/kuskoman) on Raspberry Pi)

## v1.2 (2020-06-07)

- `TARGET_NAMESPACE` now supports multiple namespaces and wildcards (Suggested by [Q-Nimbus](https://github.com/Q-Nimbus))
- Added automated tests
- Various refactoring and code restructuring

## v1.1 (2020-04-25)

- Added environment variable to specify a namespace (Contribution from [YoSmudge](https://github.com/YoSmudge))
- Add a changelog
- Add a contributor list

## v1.0 (2020-03-22)

Initial version release.
More info here: https://nabeel.dev/2020/03/22/k8s-ecr-login-renew/
