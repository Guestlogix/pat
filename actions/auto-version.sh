echo "Checking Autoversion"
CURRENT_VERSION=$(pat release rv .)
NEW_VERSION=$(pat release nv .)
echo "Current Version: $CURRENT_VERSION"
echo "New Version: $NEW_VERSION"
if [ "$CURRENT_VERSION" != "$NEW_VERSION" ]; then
    echo "Updating commit with new semver tag."
    # get repo name from git
    remote=$(git config --get remote.origin.url)
    repo=$(basename $remote .git)

    # POST a new ref to repo via Github API
    curl -X POST https://api.github.com/repos/$GITHUB_REPOSITORY/git/refs -H 'Authorization: token $GITHUB_ACCESS_TOKEN' -d '{"ref": "refs/tags/$NEW_VERSION","sha": "$GITHUB_SHA"}'
else
    echo "No new tag required."
fi