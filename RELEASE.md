# Building a new weldr-client release

Make sure you have setup the following in your git environment:
* user.signingkey
* user.email
* user.name

The signingkey is the KEYID of the gpg key you will use to sign the release
with.  This is used for both the git tag and for the detached signature for the
source archive file.

Create a branch for the release, this is needed because commits to main are
restricted to approved pull requests.  eg. `git checkout -b main-release-X.Y`

Run `make release` on this branch.  This will run the tests, bump the version,
tag the release, and create an archive with a detached gpg signature.  The key
used to sign the archive will also be output as `gpg-KEYID.key` which will be
needed for the Fedora release.

Push your branch to GitHub and create a pull request, ask someone for an ack
if you are not on the list of people who can push to main w/o review.  Also
wait for the GitHub CI to finish and pass.

**DO NOT** use the GitHub UI to merge the commit.  It may rebase/merge the
commit which would make the gpg signed tag not match.

**DO** use git on the cmdline to do a fast-forward merge of your release branch
into main.

This is so that the signed tag matches the commit you made on your release
branch.

    git co main
    git merge main-release-X.Y
    git push --follow-tags origin main

Check github to make sure the new tag has the 'Verified' icon next to it.


## Create a GitHub release

In the GitHub UI select the 'Tags' page and click on the new tag.  Select 'Create
release from tag'.

Set the title to 'weldr-client version X.Y' with the new version.  Paste the
changelog since the last release into the description and trim out irrelevant
commits -- I try to collapse multiple build(deps) entries for the same
dependency into one entry to make things more readable, and move them all to
the end of the list.

Add the .tar.gz and .tar.gz.asc archive file and signature to the release and
click on 'Publish release'


## Create a Fedora release

Run `make weldr-client.spec` to generate a new .spec file. The changelog is in
`./clog` and should be used as the commit message for the spec change.  Copy
the new `weldr-client.spec` and `./clog` to your weldr-client dist-git repo.
Also copy the `gpg-KEYID.key` file if it isn't already there.

In the dist-git repo update the vendor archive. This requires that the
`go-vendor-tools` package be installed. Run this to create the vendor archive:

    make update-vendor

Add the source archive, vendor archive, signature, and public key:

    fedpkg new-sources weldr-client-36.1.tar.gz* weldr-client-36.1-vendor.tar.bz2 gpg-*key

Commit the changes with `git add -u && git commit -F clog`.

Check that the changes look ok with `git show`, do a mock build with `fedpkg mockbuild`
and if that all looks ok, push and build:

    fedpkg push && fedpkg build
