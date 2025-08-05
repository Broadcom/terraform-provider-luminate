schema_version = 1

project {
  license        = "MPL-2.0"
  copyright_year = 2025
  copyright_holder = "Broadcom Inc."


  header_ignore = [
    "**/*_test.go",
    ".circleci/**/*",
    "utils/**/*",
    "build-scripts/**/*",
    ".goreleaser.yml",
    ".release/*.hcl",
    "service/utils/**/*",
    "service/roundtripper/**/*",
    "service/dto/**/*",
  ]
}
