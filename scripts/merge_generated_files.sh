#!/bin/sh

# discard changes in generated files
git checkout --theirs -- editor/*
git checkout --theirs -- gen/*

# regenerate files
go generate ./...

# add and commit regenerated files
git add editor/* gen/*
git commit -m "Regenerate files after merge conflict"
#!/bin/sh

# discard changes in generated files
git checkout --theirs -- editor/*
git checkout --theirs -- gen/*

# regenerate files
go generate ./...

# add and commit regenerated files
git add editor/* gen/*
