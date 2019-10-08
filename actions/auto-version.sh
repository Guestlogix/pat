echo "Checking Autoversion"
CURRENT_VERSION=$(pat release rv .)
NEW_VERSION=$(pat release nv .)
echo "Current Version: $CURRENT_VERSION"
echo "New Version: $NEW_VERSION"
if [ "$CURRENT_VERSION" != "$NEW_VERSION" ]; then
    echo "Updating Commit with new semver tag"
fi