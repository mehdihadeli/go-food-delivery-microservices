# GitFlow

we are using [GitFlow](https://www.atlassian.com/git/tutorials/comparing-workflows/gitflow-workflow) for branch management

## Branches

- Production branch: master
- Develop branch: develop
- Feature prefix: feature/
- Release prefix: release/
- Hotfix prefix: hotfix/

## Setup

```bash
# install git-flow for mac
brew install git-flow-avh
```

```bash
# Start using git-flow by initializing it inside an existing git repository
git flow init [-d] # The -d flag will accept all defaults.
# Start a new feature
git flow feature start grpc
# Finish up a feature
git flow feature finish grpc
# Publish a feature
git flow feature publish grpc
# Get a feature published by another user.
git flow feature pull origin grpc
## Make a release
# Start a release
git flow release start '0.1.0'
# It's wise to publish the release branch after creating it to allow release commits by other developers
git flow release publish '0.1.0'
## Finish up a release,
# Merges the release branch back into 'master'
# Tags the release with its name
# Back-merges the release into 'develop'
# Removes the release branch
git flow release finish '0.1.0'
# Don't forget to push your tags with
git push origin --tags
```

### FAQ

- How to keep your feature branch in sync with your develop branch?

  While using GitFlow, it is a good practice to keep your feature branch in sync with the develop branch to make merging easy.

  > I do the following to keep my feature branch in sync with develop.

  ```bash
  git checkout develop    #if you don't have it already
  git checkout feature/x  #if you don't have it already
  git pull --all
  git merge develop
  ```

- Ever annoyed by the long list of local Git branches that are no longer relevant?

  Solution:

  `git remote prune origin`

  That removes all local branches that have been deleted from remote (typically GitHub)<br/>
  Add --dry-run to merely see a list first to confirm.

## Versioning

## Changelog

> generate changelog using [git-chglog](https://github.com/git-chglog/git-chglog)

```bash
git-chglog -c .github/chglog/config.yml -o CHANGELOG.md
git-chglog -c .github/chglog/config.yml -o CHANGELOG.md --next-tag 2.0.0
```

## Reference
- [xmlking/micro-starter-kit](https://github.com/xmlking/micro-starter-kit)
- [https://nvie.com/posts/a-successful-git-branching-model/](https://nvie.com/posts/a-successful-git-branching-model/)
- <https://www.atlassian.com/git/tutorials/comparing-workflows/gitflow-workflow>
- <https://danielkummer.github.io/git-flow-cheatsheet/>
