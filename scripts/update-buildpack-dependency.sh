uri() {
  sed 's|gs://|https://storage.googleapis.com/|' "${ROOT}"/dependency/url
}

sha256() {
  shasum -a 256 "${ROOT}"/dependency/streaming-http-adapter-linux-amd64-*.tgz | cut -f 1 -d ' '
}
