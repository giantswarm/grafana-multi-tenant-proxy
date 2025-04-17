# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

- Stop pushing the proxy to capi collections

### Fixed

- Fix golangci by using filepath.Clean

## [0.8.0] - 2024-10-07

### Changed

- Add support to disable the whole chart if logging is disabled at the installation level.

## [0.7.0] - 2024-09-17

### Added

- Add basic http metrics.

## [0.6.0] - 2024-09-03

### Added

- Add helm chart to be able to deploy the proxy as a standalone component

### Changed

- Move config package from `internal` to public packages so it can be imported in our operators.
- Use giantswarm CI tooling to build
- Improve error handling.
- ⚠️ [BREAKING] Refactor configuration so that `target server` or `keep-org-id` are now configured via config instead of via a cli flag.

## [0.5.0] - 2024-02-23

### Changed

- Renamed project to `grafana-multi-tenant-proxy`
- Renamed parameter `loki-server` to `target-server`

## [0.4.0] - 2024-02-19

### Added

- Allow reload of config via the `/-/reload` endpoint.

## [0.3.0] - 2024-02-08

### Added

- Add OAuth token support for 'read' user.

## [0.2.0] - 2023-12-11

### Changed:

- Bump github.com/urfave/cli from 1.22.10 to v2.26.0
- Use zap as logger and log status code.

## [0.1.0] - 2022-10-11

## [0.0.0] - 2022-10-11

- changed: add '--keep-orgid' option
- changed: image renamed to loki-multi-tenant-proxy-gs
- changed: Setup CI
- Bump github.com/urfave/cli from 1.21.0 to 1.22.10
- Bump gopkg.in/yaml.v2 from 2.2.2 to 2.4.0

[Unreleased]: https://github.com/giantswarm/grafana-multi-tenant-proxy/compare/v0.8.0...HEAD
[0.8.0]: https://github.com/giantswarm/grafana-multi-tenant-proxy/compare/v0.7.0...v0.8.0
[0.7.0]: https://github.com/giantswarm/grafana-multi-tenant-proxy/compare/v0.6.0...v0.7.0
[0.6.0]: https://github.com/giantswarm/grafana-multi-tenant-proxy/compare/v0.5.0...v0.6.0
[0.5.0]: https://github.com/giantswarm/grafana-multi-tenant-proxy/compare/v0.4.0...v0.5.0
[0.4.0]: https://github.com/giantswarm/loki-multi-tenant-proxy/compare/v0.3.0...v0.4.0
[0.3.0]: https://github.com/giantswarm/loki-multi-tenant-proxy/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/giantswarm/loki-multi-tenant-proxy/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/giantswarm/loki-multi-tenant-proxy/compare/v0.0.0...v0.1.0
[0.0.0]: https://github.com/giantswarm/loki-multi-tenant-proxy/releases/tag/v0.0.0
