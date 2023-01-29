# Contributing

<details>
<summary>Table of Contents</summary>

- [Test](#test)
- [Lint](#lint)
- [Release](#release)

</details>

All pull requests are welcome! By participating in this project, you
agree to abide by our **[code of conduct]**.

[code of conduct]: https://github.com/demdxx/gocast/blob/master/.github/CODE_OF_CONDUCT.md

[Fork], then clone the repository:

[fork]: https://github.com/demdxx/gocast/fork

```sh
# replace `<user>` with your username
git clone git@github.com:<user>/gocast.git && cd gocast
```

Write a commit message that follows the [Conventional Commits][commit] specification.

The commit message will be linted during the pre-commit Git hook.

Push to your fork and [submit a pull request][pr].

[pr]: https://github.com/demdxx/gocast/compare/

At this point you're waiting on us. We like to comment on pull requests
within three business days (and, typically, one business day). We may suggest
changes, improvements, or alternatives.

Things that will improve the chance that your pull request will be accepted:

- [ ] Write tests that pass [CI].
- [ ] Write good documentation.
- [ ] Write a [good commit message][commit].

[ci]: https://github.com/demdxx/gocast/actions/workflows/tests.yml
[commit]: https://www.conventionalcommits.org/

## Test

Run tests with coverage:

```sh
make test
```

## Lint

Lint codebase:

```sh
make lint
```

## Release

Release and publish.
