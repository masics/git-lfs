## Contributing to Git Large File Storage

Hi there! We're thrilled that you'd like to contribute to this project. Your
help is essential for keeping it great.

Contributions to this project are [released](https://help.github.com/articles/github-terms-of-service/#6-contributions-under-repository-license) to the public under the [project's open source license](LICENSE.md).

This project adheres to the [Open Code of Conduct](./CODE-OF-CONDUCT.md). By participating, you are expected to uphold this code.

## Feature Requests

Feature requests are welcome, but will have a much better chance of being
accepted if they meet the first principles for the project. Git LFS is intended
for end users, not Git experts. It should fit into the standard workflow as
much as possible, and require little client configuration.

* Large objects are pushed to Git LFS servers during git push.
* Large objects are downloaded during git checkout.
* Git LFS servers are linked to Git remotes by default. Git hosts can support
users without requiring them to set up anything extra. Users can access
different Git LFS servers like they can with different Git remotes.
* Upload and download requests should use the same form of authentication built
into Git: SSH through public keys, and HTTPS through Git credential helpers.
* Git LFS servers use a JSON API designed around progressive enhancement.
Servers can simply host off cloud storage, or implement more efficient methods
of transferring data.

## Project Management

The Git LFS project is managed completely through this open source project. The
[milestones][] show the high level items that are prioritized for future work.
Suggestions for major features should be submitted as a pull request that adds a
markdown file to `docs/proposals` discussing the feature. This gives the
community time to discuss it before a lot of code has been written.

[milestones]: https://github.com/git-lfs/git-lfs/milestones

The Git LFS teams mark issues and pull requests with the following labels:

* `bug` - An issue describing a bug.
* `enhancement` - An issue for a possible new feature.
* `review` - A pull request ready to be reviewed.
* `release` - A checklist issue showing items marked for an upcoming release.

## Branching strategy

In general, contributors should develop on branches based off of `master` and pull requests should be to `master`.

## Submitting a pull request

1. [Fork][] and clone the repository
1. Configure and install the dependencies: `make`
1. Make sure the tests pass on your machine: `make test`
1. Create a new branch based on `master`: `git checkout -b <my-branch-name> master`
1. Make your change, add tests, and make sure the tests still pass
1. Push to your fork and [submit a pull request][pr] from your branch to `master`
1. Pat yourself on the back and wait for your pull request to be reviewed

Here are a few things you can do that will increase the likelihood of your pull request being accepted:

* Follow the [style guide][style] where possible.
* Write tests.
* Update documentation as necessary.  Commands have [man pages](./docs/man).
* Keep your change as focused as possible. If there are multiple changes you
would like to make that are not dependent upon each other, consider submitting
them as separate pull requests.
* Write a [good commit message](http://tbaggery.com/2008/04/19/a-note-about-git-commit-messages.html).

## Building

### Prerequisites

Git LFS depends on having a working Go 1.11.0+ environment.

On RHEL etc. e.g. Red Hat Enterprise Linux Server release 7.2 (Maipo), you will neet the minimum packages installed to build Git LFS:

```ShellSession
$ sudo yum install gcc
$ sudo yum install perl-Digest-SHA
```

In order to run the RPM build `rpm/build_rpms.bsh` you will also need to:

```ShellSession
$ sudo yum install ruby-devel
```

(note on an AWS instance you may first need to `sudo yum-config-manager --enable rhui-REGION-rhel-server-optional`)

### Building Git LFS

The easiest way to download Git LFS for making changes is `git clone`:

```ShellSession
$ git clone git@github.com:git-lfs/git-lfs.git
$ cd git-lfs
```

From here, run `make` to build Git LFS in the `./bin` directory. Before
submitting changes, be sure to run the Go tests and the shell integration
tests:

```ShellSession
$ make test          # runs just the Go tests
$ cd t && make test  # runs the shell tests in ./test
$ script/cibuild     # runs everything, with verbose debug output
```

## Updating 3rd party packages

1. Update `go.mod`.
1. Run `make vendor` to update the code in the `vendor` directory.
1. Commit the change.  Git LFS vendors the full source code in the repository.
1. Submit a pull request.

## Releasing

If you are the current maintainer, see
[the release howto](./docs/howto/release-git-lfs.md) for how to perform a release.

## Resources

- [Contributing to Open Source on GitHub](https://guides.github.com/activities/contributing-to-open-source/)
- [Using Pull Requests](https://help.github.com/articles/using-pull-requests/)
- [GitHub Help](https://help.github.com)

[fork]: https://github.com/git-lfs/git-lfs/fork
[pr]: https://github.com/git-lfs/git-lfs/compare
[style]: https://github.com/golang/go/wiki/CodeReviewComments
