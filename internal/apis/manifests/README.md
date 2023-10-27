The `conversion_generated.go` file is generated with this command:

```
./.cache/deps/bin/conversion-gen --input-dirs ./apis,./internal/apis/manifests -h /dev/null
```

Please note that the generated files need to be adjusted by hand at the moment to work around some codegen limitations.
