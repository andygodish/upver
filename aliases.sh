# aliases.sh
# Source this:  source ./aliases.sh

# Change this if you tag the image differently
export UPVER_IMAGE="${UPVER_IMAGE:-upver:latest}"

# Run upver in a container with your current repo mounted at /work
upverc() {
  docker run --rm -it \
    -v "$(pwd):/work" \
    -w /work \
    "${UPVER_IMAGE}" "$@"
}

# Common shortcuts
upver_plan() {
  upverc --plan "$@"
}

upver_apply() {
  upverc --apply "$@"
}

upver_changelog() {
  upverc --changelog "$@"
}

upver_apply_changelog() {
  upverc --apply --changelog "$@"
}