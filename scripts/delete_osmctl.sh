printf "INFO\tDeleting deployment......\n"
if ! ./$OS-$OSM_ARCH/osm uninstall; then
	printf "ERROR\tUnable to delete\n"
	exit 1
fi
printf "INFO\tDeleted successfully\n"