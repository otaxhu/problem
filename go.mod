module github.com/otaxhu/problem

go 1.22.0

retract (
    v1.0.0 // NewMap returned a pointer, changed to return not a pointer
)
